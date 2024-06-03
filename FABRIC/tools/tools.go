package tools

import (
	"FABRIC/Assets"
	"bufio"
	"fmt"
	"github.com/muonsoft/validation/validate"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var GRecinto Assets.Recinto

// Check if the IP address is valid
func IsValidIP(ip string) bool {
	ipState := validate.IP(ip, validate.DenyPrivateIP())
	// nil -> IP is public
	if ipState != nil {
		if strings.Contains(ipState.Error(), "prohibited") {
			return true
		}
	}
	return false
}

// Check if the MAC addres is valid
func IsValidMAC(mac string) bool {
	macRegex := regexp.MustCompile(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`)
	//Check if format is right
	if !macRegex.MatchString(mac) {
		return false
	}
	_, macState := net.ParseMAC(mac)

	//Check if mac is valid
	if macState != nil {
		return false
	}
	return true
}

// Build the dictionary
func BuildDictionary() (provincias map[int]string, cantones map[int]string, parroquias map[int]string, recintos map[int]string) {
	root, _ := os.Getwd()
	path := filepath.Join(root, "Assets")
	provincias = make(map[int]string)
	cantones = make(map[int]string)
	parroquias = make(map[int]string)
	recintos = make(map[int]string)

	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		archivo, error := os.Open(filepath.Join(path, e.Name()))

		if error != nil {
			log.Fatal(error)
		}
		switch e.Name() {
		case "provincias.txt":
			scanner := bufio.NewScanner(archivo)
			for scanner.Scan() {
				linea := scanner.Text()
				partes := strings.Split(linea, ": ")
				clave, _ := strconv.Atoi(partes[0])
				valor := partes[1]
				provincias[clave] = valor
			}
		case "cantones.txt":
			scanner := bufio.NewScanner(archivo)
			for scanner.Scan() {
				linea := scanner.Text()
				partes := strings.Split(linea, ": ")
				clave, _ := strconv.Atoi(partes[0])
				valor := partes[1]
				cantones[clave] = valor
			}
		case "parroquias.txt":
			scanner := bufio.NewScanner(archivo)
			for scanner.Scan() {
				linea := scanner.Text()
				partes := strings.Split(linea, ": ")
				clave, _ := strconv.Atoi(partes[0])
				valor := partes[1]
				parroquias[clave] = valor
			}
		case "recintoElectoral.txt":
			scanner := bufio.NewScanner(archivo)
			for scanner.Scan() {
				linea := scanner.Text()
				partes := strings.Split(linea, ": ")
				clave, _ := strconv.Atoi(partes[0])
				valor := partes[1]
				recintos[clave] = valor
			}
		}
	}
	return provincias, cantones, parroquias, recintos
}

// Parsing CDA name->Provincia-canton-parroquia-recinto
func GetCDA(parentPath string) (recinto Assets.Recinto) {
	err := filepath.Walk(parentPath, ProccesFile)
	if err != nil {
		log.Fatal(err)
	}
	return GRecinto
}

// Proces every file in the tree cda folders
func ProccesFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !info.IsDir() {
		archivo, error := os.Open(path)
		if error != nil {
			log.Fatal(error)
		}
		scanner := bufio.NewScanner(archivo)
		for scanner.Scan() {
			linea := scanner.Text()
			partes := strings.Split(linea, ";")
			ip := partes[0]
			mac := partes[1]
			if IsValidIP(ip) && IsValidMAC(mac) {
				GRecinto = queryDictionaries(info.Name(), ip, mac)
			}

		}
	}
	return nil
}

// Quey dictionaries for the values of codes
func queryDictionaries(filename string, ip string, mac string) (recinto Assets.Recinto) {
	partes := strings.Split(filename, ".")

	var provincias map[int]string
	var cantones map[int]string
	var parroquias map[int]string
	var recintos map[int]string

	provincias, cantones, parroquias, recintos = BuildDictionary()

	iprovincia, _ := strconv.Atoi(partes[0])
	icanton, _ := strconv.Atoi(partes[1])
	iparroquia, _ := strconv.Atoi(partes[2])
	irecinto, _ := strconv.Atoi(partes[3])

	recinto = Assets.Recinto{
		ID:        ip + mac,
		Provincia: provincias[iprovincia],
		Canton:    cantones[icanton],
		Parroquia: parroquias[iparroquia],
		Recinto:   recintos[irecinto],
	}
	return recinto
}

// INIT struct based on cda.txt file
func MakeFolders(docPath string) {
	archivo, error := os.Open(filepath.Join(docPath))
	if error != nil {
		log.Fatal(error)
	}
	scanner := bufio.NewScanner(archivo)
	for scanner.Scan() {
		root, _ := os.Getwd()
		path := filepath.Join(root, "Assets", "cda")
		linea := scanner.Text()
		partes := strings.Split(linea, ".")
		fmt.Println(linea)
		for _, part := range partes {
			path = filepath.Join(path, part)
			fmt.Println(path)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				err := os.MkdirAll(path, 0755)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		pathfile := filepath.Join(path, linea)
		file, _ := os.Create(pathfile)
		_, _ = file.Write([]byte("192.168.100.1;20:20:20:20:20:20"))
		defer file.Close()

	}
}

//os.Open() opens a file f.Close() closes it. io.ReadAll()
