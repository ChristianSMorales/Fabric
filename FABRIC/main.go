/*package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func main() {
	// Generate RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// Encode private key to PEM format
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	// Write private key to a file
	privateKeyFile, err := os.Create("private.pem")
	if err != nil {
		panic(err)
	}
	defer privateKeyFile.Close()
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		panic(err)
	}

}
*/

package main

import (
	"FABRIC/tools"
	"fmt"
	"os"
	"path/filepath"
)

func main() {

	root, _ := os.Getwd()
	path := filepath.Join(root, "Assets", "cda")
	/*
		var provincias map[int]string
		var cantones map[int]string
		var parroquias map[int]string
		var recintos map[int]string
		println(root)
		provincias, cantones, parroquias, recintos = tools.BuildDictionary(path)

		fmt.Println("Provincias:", provincias[1])
		fmt.Println("Cantones:", cantones[5])
		fmt.Println("Parroquias:", parroquias[5])
		fmt.Println("Recintos:", recintos[1])

	*/
	fmt.Println(tools.GetCDA(path))
}
