package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
)

func main() {
	savePrivateFileTo := "./id_rsa_test"
	savePublicFileTo := "./id_rsa_test.pub"
	bitSize := 4096

	privateKey, err := generatePrivateKey(bitSize)
	if err != nil {
		log.Fatal(err.Error())
	}

	publicKeyBytes, err := generatePublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatal(err.Error())
	}

	privateKeyBytes := encodePrivateKeyToPEM(privateKey)

	err = writeKeyToFile(privateKeyBytes, savePrivateFileTo)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = writeKeyToFile([]byte(publicKeyBytes), savePublicFileTo)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// generatePrivateKey creates a RSA Private Key of specified byte size
func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	log.Println("Private Key generated")
	return privateKey, nil
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

// generatePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func generatePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	log.Println("Public key generated")
	return pubKeyBytes, nil
}

// writePemToFile writes keys to a file
func writeKeyToFile(keyBytes []byte, saveFileTo string) error {
	err := ioutil.WriteFile(saveFileTo, keyBytes, 0600)
	if err != nil {
		return err
	}

	log.Printf("Key saved to: %s", saveFileTo)
	return nil
}



func GenerateKey(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
    private, err := rsa.GenerateKey(rand.Reader, bits)
    if err != nil {
        return nil, nil, err
    }
    return private, &private.PublicKey, nil

}

func EncodePrivateKey(private *rsa.PrivateKey) []byte {
    return pem.EncodeToMemory(&pem.Block{
        Bytes: x509.MarshalPKCS1PrivateKey(private),
        Type:  "RSA PRIVATE KEY",
    })
}

func EncodePublicKey(public *rsa.PublicKey) ([]byte, error) {
    publicBytes, err := x509.MarshalPKIXPublicKey(public)
    if err != nil {
        return nil, err
    }
    return pem.EncodeToMemory(&pem.Block{
        Bytes: publicBytes,
        Type:  "PUBLIC KEY",
    }), nil
}

//EncodeSSHKey
func EncodeSSHKey(public *rsa.PublicKey) ([]byte, error) {
    publicKey, err := ssh.NewPublicKey(public)
    if err != nil {
        return nil, err
    }
    return ssh.MarshalAuthorizedKey(publicKey), nil
}

func MakeSSHKeyPair() (string, string, error) {

    pkey, pubkey, err := GenerateKey(2048)
    if err != nil {
        return "", "", err
    }

    pub, err := EncodeSSHKey(pubkey)
    if err != nil {
        return "", "", err
    }

    //glog.Info("privateKey=[%s]\n pubKey=[%s]",string(EncodePrivateKey(pkey)),string(pub))
    return string(EncodePrivateKey(pkey)), string(pub), nil
}
