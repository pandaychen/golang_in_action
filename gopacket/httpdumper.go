package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/reassembly"
)

var fname = flag.String("filename", "", "Pcap file to parse")
var filter = flag.String("filter", "", "BPF Filter to apply to PCAP")
var dumpsite = flag.String("output", "extracted", "Path where to extract files")

type httpReader struct {
	ident    string
	isClient bool
	bytes    chan []byte
	data     []byte
	parent   *tcpStream
}

func (hR *httpReader) Read(bytes []byte) (int, error) {
	ok := true
	for len(hR.data) == 0 && ok {
		hR.data, ok = <-hR.bytes
	}
	if !ok || len(hR.data) == 0 {
		return 0, io.EOF
	}
	l := copy(bytes, hR.data)
	hR.data = hR.data[l:]
	return l, nil
}

func (hR *httpReader) run(wg *sync.WaitGroup) {
	defer wg.Done()
	b := bufio.NewReader(hR)
	for {
		if hR.isClient {
			//Client Request
			req, err := http.ReadRequest(b)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			} else if err != nil {
				log.Println("error parsing request ", err)
				continue
			}
			body, err := ioutil.ReadAll(req.Body)
			req.Body.Close()
			if err != nil {
				log.Println("error reading body - ", err)
			}
			if len(body) > 0 {
				log.Println("Body length: ", len(body))
			}
		} else {
			// Server Response
			res, err := http.ReadResponse(b, nil)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			} else if err != nil {
				log.Println("error reading response - ", err)
				continue
			}
			// Read content from Body
			content, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Println("error reading body - ", err)
			}
			res.Body.Close()
			//contentType := res.Header.Get("Content-Type")
			contentType, ok := res.Header["Content-Type"]
			if !ok {
				contentType = []string{http.DetectContentType(content)}
			}
			encoding := res.Header["Content-Encoding"]
			if contentType == nil {
				log.Println("response Headers has no content in ['Content-Type']")
				continue
			} else if isFilteredContentType(contentType) {
				contentlength := len(content)
				if contentlength > 0 {
					// Store file content
					baseName := path.Join(*dumpsite, hR.ident)
					_, err = os.Stat(baseName)
					if err == nil {
						baseName = baseName + "_1"
					}
					f, err := os.Create(baseName)
					if err != nil {
						log.Println("could not create file ", baseName, " error: ", err)
					}
					defer f.Close()
					var r io.Reader
					r = bytes.NewBuffer(content)
					// If  it is encoded we might need to do something
					if len(encoding) > 0 && (encoding[0] == "deflate" || encoding[0] == "gzip") {
						r, err = gzip.NewReader(r)
						if err != nil {
							log.Println("error - could not decode with ", encoding, " error: ", err)
						}
					}
					//  If there is no error we can try to write the file
					if err == nil {
						w, err := io.Copy(f, r)
						if _, ok = r.(*gzip.Reader); ok {
							r.(*gzip.Reader).Close()
						}
						if err != nil {
							log.Println(err)
						} else {
							log.Printf("written %d bytes to file %s\n", w, baseName)
						}
					}
				}
			}
		}
	}
}

func isFilteredContentType(contentTypes []string) bool {
	filteredTypes := map[string]bool{
		"application/x-msdownload": true,
		"application/octet-stream": true,
	}
	for _, cType := range contentTypes {
		if _, ok := filteredTypes[cType]; ok {
			return true
		}
	}
	return false
}

type Context struct {
	CaptureInfo gopacket.CaptureInfo
}

func (ac *Context) GetCaptureInfo() gopacket.CaptureInfo {
	return ac.CaptureInfo
}

// Implements a reassembly.Stream
type tcpStream struct {
	net, transport gopacket.Flow
	tcpstate       *reassembly.TCPSimpleFSM
	optchecker     reassembly.TCPOptionCheck
	reversed       bool
	httpClient     httpReader
	httpServer     httpReader
	ident          string
}

func (tS *tcpStream) Accept(tcp *layers.TCP, ci gopacket.CaptureInfo, dir reassembly.TCPFlowDirection, nextSeq reassembly.Sequence, start *bool, ac reassembly.AssemblerContext) bool {
	if !tS.tcpstate.CheckState(tcp, dir) {
		return false
	}
	if err := tS.optchecker.Accept(tcp, ci, dir, nextSeq, start); err != nil {
		return false
	}
	return true
}

func (tS *tcpStream) ReassembledSG(sg reassembly.ScatterGather, ac reassembly.AssemblerContext) {
	dir, _, _, _ := sg.Info()
	length, _ := sg.Lengths()

	data := sg.Fetch(length)

	if length > 0 {
		if dir == reassembly.TCPDirClientToServer && !tS.reversed {
			tS.httpClient.bytes <- data
		} else {
			tS.httpServer.bytes <- data
		}
	}
}

func (tS *tcpStream) ReassemblyComplete(ac reassembly.AssemblerContext) bool {
	close(tS.httpClient.bytes)
	close(tS.httpServer.bytes)
	return false
}

// Implements Interface reassembly.StreamFactory
type tcpStreamFactory struct {
	wg sync.WaitGroup
}

func (tSF *tcpStreamFactory) New(netFlow gopacket.Flow, tcpFlow gopacket.Flow, tcp *layers.TCP, ac reassembly.AssemblerContext) reassembly.Stream {
	fsmOptions := reassembly.TCPSimpleFSMOptions{
		SupportMissingEstablishment: false,
	}

	stream := &tcpStream{
		net:        netFlow,
		transport:  tcpFlow,
		tcpstate:   reassembly.NewTCPSimpleFSM(fsmOptions),
		optchecker: reassembly.NewTCPOptionCheck(),
		reversed:   tcp.SrcPort == 80,
		ident:      fmt.Sprintf("%s - %s", netFlow, tcpFlow),
	}
	if tcp.SrcPort == 80 || tcp.DstPort == 80 {
		stream.httpClient = httpReader{
			ident:    fmt.Sprintf("%s - %s", netFlow, tcpFlow),
			bytes:    make(chan []byte),
			isClient: true,
			parent:   stream,
		}
		stream.httpServer = httpReader{
			ident:    fmt.Sprintf("%s - %s", netFlow, tcpFlow),
			bytes:    make(chan []byte),
			isClient: false,
			parent:   stream,
		}
		tSF.wg.Add(2)
		go stream.httpClient.run(&tSF.wg)
		go stream.httpServer.run(&tSF.wg)
	}
	return stream
}

func (tSF *tcpStreamFactory) WaitGoRoutines() {
	tSF.wg.Wait()
}

func main() {
	log.Println("start")
	defer log.Println("end")
	flag.Parse()

	var handle *pcap.Handle
	var err error

	if *fname == "" {
		log.Fatal("no --filename parameter passed")
	}

	handle, err = pcap.OpenOffline(*fname)
	if err != nil {
		log.Fatalf("could not oppen filename - %v - %s", *fname, err)
	}

	if *filter != "" {
		if err = handle.SetBPFFilter(*filter); err != nil {
			log.Fatalf("could not apply filter %v to capture - %s", *filter, err)
		}
	}

	source := gopacket.NewPacketSource(handle, handle.LinkType())
	source.Lazy = false
	source.NoCopy = true

	// Create StreamFactory
	streamFactory := &tcpStreamFactory{}
	// Create StreamPool
	streamPool := reassembly.NewStreamPool(streamFactory)
	// Create Assembler
	reassembler := reassembly.NewAssembler(streamPool)

	const closeTimeout time.Duration = time.Hour * 1
	const timeout time.Duration = time.Minute * 1

	count := 0

	for packet := range source.Packets() {

		count++
		//Parse Packet
		if packet == nil {
			return
		}
		if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
			continue
		}

		tcp := packet.Layer(layers.LayerTypeTCP)

		if tcp != nil {
			tcp := tcp.(*layers.TCP)
			context := Context{
				CaptureInfo: packet.Metadata().CaptureInfo,
			}
			if tcp.SrcPort == 80 || tcp.DstPort == 80 {
				reassembler.AssembleWithContext(packet.NetworkLayer().NetworkFlow(), tcp, &context)
			}
		}

		if count%1000 == 0 {
			timestamp := packet.Metadata().CaptureInfo.Timestamp
			reassembler.FlushWithOptions(reassembly.FlushOptions{T: timestamp.Add(-timeout), TC: timestamp.Add(-closeTimeout)})
		}

	}

	log.Println("iterated all packets")

	reassembler.FlushAll()
	log.Println("flushed all connections")
	streamFactory.WaitGoRoutines()
	log.Println("all go routines finished")

}
