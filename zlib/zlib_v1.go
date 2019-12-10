package main
 
import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
)
 
func main() {
	var in bytes.Buffer
	b := []byte(`abcdefghijklmnopqrstuvwxyz`)
	w := zlib.NewWriter(&in)
	w.Write(b)
	w.Close()

	fmt.Println(in.Len(),in.Cap(),len(b)) 

	var out bytes.Buffer
	r, _ := zlib.NewReader(&in)
	io.Copy(&out, r)
	fmt.Println(out.String())
 
}
