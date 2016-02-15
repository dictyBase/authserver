package commands

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	"gopkg.in/codegangsta/cli.v1"
)

// Generate RSA public and private keys in PEM format
func GenerateKeys(c *cli.Context) {
	// validate
	if !c.IsSet("public") {
		log.Fatal("public key output file is not provided")
	}
	if !c.IsSet("private") {
		log.Fatal("private key output file is not provided")
	}

	// open files
	prvWriter, err := os.Create(c.String("private"))
	defer prvWriter.Close()
	if err != nil {
		log.Fatalf("unable to create private key file %q\n", err)
	}
	pubWriter, err := os.Create(c.String("public"))
	defer pubWriter.Close()
	if err != nil {
		log.Fatalf("unable to create public key file %q\n", err)
	}

	// generate and write to files
	private, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("error in generating private key %q\n", err)
	}
	if err := private.Validate(); err != nil {
		log.Fatalf("error in validating private key %q\n", err)
	}
	prvPem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(private),
	}
	public := private.Public()
	pubCont, err := x509.MarshalPKIXPublicKey(public)
	if err != nil {
		log.Fatalf("unable to marshall private key %q\n", err)
	}
	pubPem := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubCont,
	}
	if err := pem.Encode(prvWriter, prvPem); err != nil {
		log.Fatalf("unable to write private key %q\n", err)
	}
	if err := pem.Encode(pubWriter, pubPem); err != nil {
		log.Fatalf("unable to write public key %q\n", err)
	}

}
