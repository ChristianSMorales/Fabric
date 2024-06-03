package cryptoMaterial

import (
	"FABRIC/tools"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

var root, _ = os.Getwd()
var path = filepath.Join(root, "Assets", "cda")

func genKeys() (priv string, pub string) {
	//Generate RSA 2048 priv key lenght
	//private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	// Encode private key to PEM format
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	//public key
	publicKey := privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}
	//Encode the public key to PEM format
	publicKeyPEM := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	privkey := pem.EncodeToMemory(privateKeyPEM)
	pubkey := pem.EncodeToMemory(publicKeyPEM)

	return string(privkey), string(pubkey)
}

// provincia-canton-parroquia-RecintoElectoral
func genCert(pathFolder string) {
	privKey, pubKey := genKeys()
	cda := tools.GetCDA(path)
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:            []string{"Ecuador"},
			Organization:       []string{"Concejo Nacional Electoral"},
			OrganizationalUnit: []string{cda.Recinto},
			Locality:           []string{cda.Canton},
			Province:           []string{cda.Provincia},
			StreetAddress:      []string{cda.Parroquia},
			CommonName:         cda.ID,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0), // válido por 1 año
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	finalPath := filepath.Join(root, pathFolder)

	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, pubKey, privKey)

	// Codificar el certificado en formato PEM
	certOut, err := os.Create(filepath.Join(finalPath, "new_cert.pem"))
	if err != nil {
		fmt.Println("Error al crear el archivo new_cert.pem:", err)
		return
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	if err != nil {
		fmt.Println("Error al crear el certificado:", err)
		return
	}

	// Write private key to a file
	if err := os.WriteFile(filepath.Join(finalPath, "private_key.pem"), []byte(privKey), 0600); err != nil {
		panic(err)
	}
	// Write public key to a file
	if err := os.WriteFile(filepath.Join(finalPath, "public_key.pem"), []byte(pubKey), 0644); err != nil {
		panic(err)
	}
}
func cryptMaterial(ipmac string) (privkey string, pubkey string, cert string) {

}
