package main

import (
	//"crypto/rand"
	//"crypto/rsa"
	//"crypto/x509"
	//"encoding/pem"
	"golang.org/x/crypto/ssh"
	//"io/ioutil"
	"fmt"
)

/*

You probably want to use ssh.ParseAuthorizedKey from the ssh package to load the key:

https://godoc.org/golang.org/x/crypto/ssh#ParseAuthorizedKey

That will give you a public key which you can call ssh.FingerprintLegacyMD5 on in order to get the fingerprint (assuming here you want the md5).

https://godoc.org/golang.org/x/crypto/ssh#FingerprintLegacyMD5 https://godoc.org/golang.org/x/crypto/ssh#FingerprintSHA256
*/

func main() {
	// Read a key from a file in authorized keys file line format
	// This could be an rsa.pub file or a line from authorized_keys
	pubKeyBytes := []byte(`ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCzMETD7LIf/P0IjgjpiTqgrrxsu1UfCfMhD8htNLgek6vIWWKVCJOajHHVMUu4ce5kA2s5Rh9EzolJ/IWRXEghx/SpW+/CP4Dhn7Q909UJuMQy2sN74dO0xpdD1tBimSHYmPkxs6XzxqKDwJYV77d1sZZCwNXvw8UEsBijK3B/dgFHSUGnX2jrTWSbIALVlbp9P3x3i0ypXK84XY8FIhPWduZ6fFlbUb14aTyJgQgw1oghGYOFhv/B48A/t9F3cE15xXqxDKsyWRDxnoxJaJD9iyEfa91NbqO6u1sb+ByLEE/i7C0UGJryqifcDMM0nDEF89RG+DjrSOb7x43u+83ixXYFx+eIfqwvAvXO+SVxPjX5yPRN8x0ybi5TiAszjioIMGRx0hUi74ugZApDTYsNv45G6cpNEEaiLZAR9qGUpoUnGGYtMJSmSg4mRO84QJtXcah3vnOAG5X83KgwDWYDqmhyewxG7kCaY/pybP+pSV1QCx146QTJN7jL3nNZN+70m3/TQ64p4vxCnlwOd8chopnLpPe6G8BgFpxRr4lAsfDe5nN9xTZFR0+2TfdqkRrplWmV4JoGuFTuz8VOfzbnvPwWDgxPhfGC9bn32ZnRM7O2syR+YT4BbcEU7epVk6pLmSFq7lBspWuIgJGawfUIUl02BPen2nMdoPrZRLeMpQ==`)

	// Parse the key, other info ignored
	pk, _, _, _, err := ssh.ParseAuthorizedKey(pubKeyBytes)
	if err != nil {
		panic(err)
	}

	// Get the fingerprint
	f := ssh.FingerprintLegacyMD5(pk)

	// Print the fingerprint
	fmt.Printf("%s\n", f)
}
