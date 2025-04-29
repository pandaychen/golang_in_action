package main

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/vishvananda/netlink/nl"
	"golang.org/x/sys/unix"
)

// @TODO: /proc/net/packet
// https://github.com/osquery/osquery/blob/f9282c0f03d049e0f13670afa2cf8a87f8ddf0cc/osquery/filesystem/linux/proc.cpp
// osquery中用户态获取socket方式 https://github.com/osquery/osquery/blob/f9282c0f03d049e0f13670afa2cf8a87f8ddf0cc/osquery/tables/networking/linux/process_open_sockets.cpp
// 在 osquery issue 1094 中(https://github.com/osquery/osquery/issues/1094) 解释了为什么剔除了用 netlink 获取的方式
// 大致为 netlink 的方式在 CentOS/RHEL6 不稳定, 经常会 fallback
// 可以看到之前 readnetlink 他们也有出现 timeout 的情况 https://github.com/osquery/osquery/pull/543
// 其他相关 issue: https://github.com/osquery/osquery/issues/671

// In Elkeid, socket rebuild again for better performance. By the way, since there is no race condition
// of netlink function execution, no netlink socket singleton or lock is needed in such situation.
// The source code is from: https://github.com/vishvananda/netlink/blob/main/socket_linux.go
const (
	sizeofSocketID      = 0x30
	sizeofSocketRequest = sizeofSocketID + 0x8
	sizeofSocket        = sizeofSocketID + 0x18
	netlinkLimit        = 1500 // max socket we get from netlink
)

var (
	native       = nl.NativeEndian()
	networkOrder = binary.BigEndian
)

type readBuffer struct {
	Bytes []byte
	pos   int
}

func (b *readBuffer) Read() byte {
	c := b.Bytes[b.pos]
	b.pos++
	return c
}

func (b *readBuffer) Next(n int) []byte {
	s := b.Bytes[b.pos : b.pos+n]
	b.pos += n
	return s
}

// what we define
type Socket struct {
	DPort     string `json:"dport" mapstructure:"dport"`
	SPort     string `json:"sport" mapstructure:"sport"`
	DIP       string `json:"dip" mapstructure:"dip"`
	SIP       string `json:"sip" mapstructure:"sip"`
	Interface string `json:"interface" mapstructure:"interface"`
	Family    string `json:"family" mapstructure:"family"`
	State     string `json:"state" mapstructure:"state"`
	UID       string `json:"uid" mapstructure:"uid"`
	Username  string `json:"username" mapstructure:"username"`
	Inode     string `json:"inode" mapstructure:"inode"`
	PID       string `json:"pid" mapstructure:"pid"`
	Cmdline   string `json:"cmdline" mapstructure:"cmdline"`
	Comm      string `json:"comm" mapstructure:"comm"`
	Type      string `json:"type" mapstructure:"type"`
}

type socketRequest struct {
	Family   uint8
	Protocol uint8
	Ext      uint8
	pad      uint8
	States   uint32
	ID       _socketID
}

type writeBuffer struct {
	Bytes []byte
	pos   int
}

func (b *writeBuffer) Write(c byte) {
	b.Bytes[b.pos] = c
	b.pos++
}

func (b *writeBuffer) Next(n int) []byte {
	s := b.Bytes[b.pos : b.pos+n]
	b.pos += n
	return s
}

func (r *socketRequest) Serialize() []byte {
	b := writeBuffer{Bytes: make([]byte, sizeofSocketRequest)}
	b.Write(r.Family)
	b.Write(r.Protocol)
	b.Write(r.Ext)
	b.Write(r.pad)
	native.PutUint32(b.Next(4), r.States)
	networkOrder.PutUint16(b.Next(2), r.ID.SourcePort)
	networkOrder.PutUint16(b.Next(2), r.ID.DestinationPort)
	if r.Family == unix.AF_INET6 {
		copy(b.Next(16), r.ID.Source)
		copy(b.Next(16), r.ID.Destination)
	} else {
		copy(b.Next(4), r.ID.Source.To4())
		b.Next(12)
		copy(b.Next(4), r.ID.Destination.To4())
		b.Next(12)
	}
	native.PutUint32(b.Next(4), r.ID.Interface)
	native.PutUint32(b.Next(4), r.ID.Cookie[0])
	native.PutUint32(b.Next(4), r.ID.Cookie[1])
	return b.Bytes
}

func (r *socketRequest) Len() int { return sizeofSocketRequest }

// Add limitation of socket, in case too much of this.
func parseNetlink(family, protocol uint8, state uint32) (sockets []Socket, err error) {
	var (
		s   *nl.NetlinkSocket
		req *nl.NetlinkRequest
	)
	// precheck protocol
	if protocol != unix.IPPROTO_UDP && protocol != unix.IPPROTO_TCP {
		err = fmt.Errorf("unsupported protocol %d", protocol)
		return
	}
	// subscribe the netlink
	if s, err = nl.Subscribe(unix.NETLINK_INET_DIAG); err != nil {
		return
	}
	defer s.Close()
	// send the netlink request
	req = nl.NewNetlinkRequest(nl.SOCK_DIAG_BY_FAMILY, unix.NLM_F_DUMP)
	req.AddData(&socketRequest{
		Family:   family,
		Protocol: protocol,
		Ext:      (1 << (INET_DIAG_VEGASINFO - 1)) | (1 << (INET_DIAG_INFO - 1)),
		States:   uint32(1 << state),
	})
	if err = s.Send(req); err != nil {
		return
	}
loop:
	for i := 1; i < netlinkLimit; i++ {
		var msgs []syscall.NetlinkMessage
		var from *unix.SockaddrNetlink
		msgs, from, err = s.Receive()
		if err != nil {
			return
		}
		if from.Pid != nl.PidKernel {
			continue
		}
		if len(msgs) == 0 {
			break
		}
		for _, m := range msgs {
			switch m.Header.Type {
			case unix.NLMSG_DONE:
				break loop
			case unix.NLMSG_ERROR:
				break loop
			}
			sockInfo := &_socket{}
			if err := sockInfo.deserialize(m.Data); err != nil {
				continue
			}
			socket := Socket{
				SIP:       sockInfo.ID.Source.String(),
				DIP:       sockInfo.ID.Destination.String(),
				SPort:     strconv.Itoa(int(sockInfo.ID.SourcePort)),
				DPort:     strconv.Itoa(int(sockInfo.ID.DestinationPort)),
				UID:       strconv.FormatUint(uint64(sockInfo.UID), 10),
				Interface: strconv.FormatUint(uint64(sockInfo.ID.Interface), 10),
				Family:    strconv.FormatUint(uint64(sockInfo.Family), 10),
				State:     strconv.FormatUint(uint64(sockInfo.State), 10),
				Inode:     strconv.FormatUint(uint64(sockInfo.INode), 10),
				Type:      strconv.FormatUint(uint64(protocol), 10),
			}
			//socket.Username = user.Cache.GetUser(sockInfo.UID).Username
			sockets = append(sockets, socket)
		}
	}
	return
}

func parseIP(hexIP string) (net.IP, error) {
	var byteIP []byte
	byteIP, err := hex.DecodeString(hexIP)
	if err != nil {
		return nil, fmt.Errorf("cannot parse address field in socket line %q", hexIP)
	}
	switch len(byteIP) {
	case 4:
		return net.IP{byteIP[3], byteIP[2], byteIP[1], byteIP[0]}, nil
	case 16:
		i := net.IP{
			byteIP[3], byteIP[2], byteIP[1], byteIP[0],
			byteIP[7], byteIP[6], byteIP[5], byteIP[4],
			byteIP[11], byteIP[10], byteIP[9], byteIP[8],
			byteIP[15], byteIP[14], byteIP[13], byteIP[12],
		}
		return i, nil
	default:
		return nil, fmt.Errorf("unable to parse IP %s", hexIP)
	}
}

// Refernce: https://guanjunjian.github.io/2017/11/09/study-8-proc-net-tcp-analysis/
func parseProcNet(family, protocol uint8, path string) (sockets []Socket, err error) {
	var (
		file *os.File
		r    *bufio.Scanner
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()
	r = bufio.NewScanner(io.LimitReader(file, 1024*1024*2))
	header := make(map[int]string)
	for i := 0; r.Scan(); i++ {
		if i == 0 {
			header[0] = "sl"
			header[1] = "local_address"
			header[2] = "rem_address"
			header[3] = "st"
			header[4] = "queue"
			header[5] = "t"
			header[6] = "retrnsmt"
			header[7] = "uid"
			for index, field := range strings.Fields(r.Text()[strings.Index(r.Text(), "uid")+3:]) {
				header[8+index] = field
			}
		} else {
			socket := Socket{
				Family: strconv.FormatUint(uint64(family), 10),
				Type:   strconv.FormatUint(uint64(protocol), 10),
			}
			droped := false
			for index, key := range strings.Fields(r.Text()) {
				switch header[index] {
				case "local_address":
					fields := strings.Split(key, ":")
					if len(fields) != 2 {
						droped = true
						break
					}
					var sip net.IP
					sip, err = parseIP(fields[0])
					if err != nil {
						droped = true
						break
					}
					socket.SIP = sip.String()
					var port uint64
					port, err = strconv.ParseUint(fields[1], 16, 64)
					if err != nil {
						droped = true
						break
					}
					socket.SPort = strconv.Itoa(int(port))
				case "rem_address":
					fields := strings.Split(key, ":")
					if len(fields) != 2 {
						droped = true
						break
					}
					var dip net.IP
					dip, err = parseIP(fields[0])
					if err != nil {
						droped = true
						break
					}
					socket.DIP = dip.String()
					var port uint64
					port, err = strconv.ParseUint(fields[1], 16, 64)
					if err != nil {
						droped = true
						break
					}
					socket.DPort = strconv.Itoa(int(port))
				case "st":
					st, err := strconv.ParseUint(key, 16, 64)
					if err != nil {
						continue
					}
					if protocol == unix.IPPROTO_UDP && st != 7 {
						droped = true
						break
					}

					if protocol == unix.IPPROTO_TCP && (st != 10 && st != 1) {
						droped = true
						break
					}
					socket.State = strconv.FormatUint(st, 10)
				case "uid":
					uid, err := strconv.ParseUint(key, 0, 64)
					if err != nil {
						continue
					}
					socket.UID = strconv.FormatUint(uid, 10)
				case "inode":
					inode, err := strconv.ParseUint(key, 0, 64)
					if err != nil {
						continue
					}
					socket.Inode = strconv.FormatUint(uint64(inode), 10)
				default:
				}
			}
			if !droped && len(socket.DIP) != 0 && len(socket.SIP) != 0 && socket.State != "0" {
				sockets = append(sockets, socket)
			}
		}
	}
	return
}

// add limitation of this
func FromProc() (sockets []Socket, err error) {
	tcpSocks, err := parseProcNet(unix.AF_INET, unix.IPPROTO_TCP, "/proc/net/tcp")
	if err != nil {
		return
	}
	sockets = append(sockets, tcpSocks...)
	tcp6Socks, err := parseProcNet(unix.AF_INET6, unix.IPPROTO_TCP, "/proc/net/tcp6")
	if err == nil {
		sockets = append(sockets, tcp6Socks...)
	}
	udpSocks, err := parseProcNet(unix.AF_INET, unix.IPPROTO_UDP, "/proc/net/udp")
	if err == nil {
		sockets = append(sockets, udpSocks...)
	}
	udp6Socks, err := parseProcNet(unix.AF_INET6, unix.IPPROTO_UDP, "/proc/net/udp6")
	if err == nil {
		sockets = append(sockets, udp6Socks...)
	}
	inodeMap := make(map[string]int)
	for index, socket := range sockets {
		if socket.Inode != "0" {
			inodeMap[socket.Inode] = index
		}
	}
	return
}

func FromNetlink() (sockets []Socket, err error) {
	var udpSockets, udp6Sockets, tcpSockets, tcp6Sockets []Socket
	udpSockets, err = parseNetlink(unix.AF_INET, unix.IPPROTO_UDP, 7)
	if err != nil {
		return
	}
	sockets = append(sockets, udpSockets...)
	udp6Sockets, err = parseNetlink(unix.AF_INET6, unix.IPPROTO_UDP, 7)
	if err != nil {
		return
	}
	// TCP - sockets for both established & listen, any better for state? dig out this
	sockets = append(sockets, udp6Sockets...)
	tcpSockets, err = parseNetlink(unix.AF_INET, unix.IPPROTO_TCP, 1)
	if err != nil {
		return
	}
	sockets = append(sockets, tcpSockets...)
	tcpSockets, err = parseNetlink(unix.AF_INET, unix.IPPROTO_TCP, 10)
	if err != nil {
		return
	}
	sockets = append(sockets, tcpSockets...)
	tcp6Sockets, err = parseNetlink(unix.AF_INET6, unix.IPPROTO_TCP, 1)
	if err != nil {
		return
	}
	sockets = append(sockets, tcp6Sockets...)
	tcp6Sockets, err = parseNetlink(unix.AF_INET6, unix.IPPROTO_TCP, 10)
	if err != nil {
		return
	}
	sockets = append(sockets, tcp6Sockets...)
	return
}

// To learn the way osquery get sockets, we go through the source code of osquery
// 1. Collect all sockets from from /proc/<pid>/fd and search for the links of
//    type of socket:[<inode>], and we get the relationship of pid - inode(socket)
// 2. Get <pid> ns/net -> inode, execute step 3 every time once a new inode is found
// 3. Get & parse the tcp/tcp6/udp/udp6 from /net/ of every pid.
//
// https://github.com/osquery/osquery/pull/608, as metioned in this pull request.
// netlink is somehow faster than /proc/ way.

func main() {
	sockets, err := FromNetlink()
	if err != nil {
		sockets, _ = FromProc()
	}
	inodeMap := make(map[string]int)
	for index, socket := range sockets {
		fmt.Println(index, socket)
		if socket.Inode != "0" {
			inodeMap[socket.Inode] = index
		}
	}
}
