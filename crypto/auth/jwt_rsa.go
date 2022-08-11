package main

//generate rsa key pair
//openssl genrsa -out id_rsa 4096
//openssl rsa -in id_rsa -pubout -out id_rsa.pub

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWT struct {
	privateKey []byte
	publicKey  []byte
}

func NewJWT(privateKey []byte, publicKey []byte) JWT {
	return JWT{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (j JWT) Create(ttl time.Duration, content interface{}) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(j.privateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["dat"] = content             // Our custom data.
	claims["exp"] = now.Add(ttl).Unix() // The expiration time after which the token must be disregarded.
	claims["iat"] = now.Unix()          // The time at which the token was issued.
	claims["nbf"] = now.Unix()          // The time before which the token must be disregarded.

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func (j JWT) Validate(token string) (interface{}, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(j.publicKey)
	if err != nil {
		return "", fmt.Errorf("validate: parse key: %w", err)
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("validate: invalid")
	}

	return claims["dat"], nil
}

func main() {
	prvKey, err := ioutil.ReadFile("id_rsa")
	if err != nil {
		log.Fatalln(err)
	}
	pubKey, err := ioutil.ReadFile("id_rsa.pub")
	if err != nil {
		log.Fatalln(err)
	}

	jwtToken := NewJWT(prvKey, pubKey)

	// 1. Create a new JWT token.
	tok, err := jwtToken.Create(time.Hour, "abcdefghijklmnopqrst")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("TOKEN:", tok)

	// 2. Validate an existing JWT token.
	content, err := jwtToken.Validate(tok)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("CONTENT:", content)
}
