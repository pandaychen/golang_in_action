package main

/*
	MD5SUM Generate,with lowest memory alloc
	use dd if=/dev/zero of=test bs=1M count=1000 to generate test file
*/

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func md5sum(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()
	srchandle := bufio.NewReader(f)

	dsthandle := md5.New()

	//read-from and then write-to
	_, err = io.Copy(dsthandle, srchandle)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", dsthandle.Sum(nil))

}

func main() {
	value := md5sum("./bigfile")
	fmt.Println(value)
}
