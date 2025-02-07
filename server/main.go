package main

import (
	"bufio"
	"dws-sk-fs/utils"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"strconv"
	"crypto/sha512"
)

// Upload
// Downlaoad
const (
	UPLOAD = iota + 1
	DOWNLOAD
)

type FileDownload struct {
	FileIdLength int
	FileId       string
	FileLength   int
	File         []byte
}

// Om fil ej finns skicka 0
// Om fil finns skicka 1
// Skicka storlek på filen
// Skicka filen
func handleDownload(conn net.Conn) {
	filename, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	filename = strings.TrimSpace(filename)

	file, err := os.Open(filename)
	// If file not found, send byte 0 to client
	// If file found, send byte 1 to client
	if err != nil {
		conn.Write([]byte{0})
		return
	}
	conn.Write([]byte{1})

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fileSize := fileInfo.Size()
	// Send file size to client as 8 bytes
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(fileSize))

	conn.Write(b)
	io.Copy(conn, file)
}

func handleHashValidationRequest(conn net.Conn) {
	filename, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	filename = strings.TrimSpace(filename)

	file, err := os.ReadFile(filename)
	// If file not found, send byte 0 to client
	// If file found, send byte 1 to client
	if err != nil {
		conn.Write([]byte{0})
		return
	}
	conn.Write([]byte{1})
	b := sha512.Sum512(file)
	conn.Write(b[:])
	
	// Send Hash of file
}

func handleUpload(conn net.Conn) {
	// Read 8 bit id
	id, err := bufio.NewReader(conn).ReadByte()
	if err != nil {
		log.Fatal(err)
	}
	// Make file
	file, err := os.Create(fmt.Sprintf("fs/%d", id))
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(file, conn)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		// Read 1 byte to determine if client wants to upload or download
		_, err := bufio.NewReader(conn).ReadByte()
		if err != nil {
			log.Fatal(err)
		}
		handleUpload(conn)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request, m *utils.MerkleTree) {
    request := utils.UploadRequest{}
    jsonRequest := json.NewDecoder(r.Body)
    jsonRequest.Decode(&request)
    file, err := os.Create(fmt.Sprintf("fs/%d.dat", request.Id))
    if err != nil {
        log.Fatal(err)
    }
    file.Write(request.File.Data)
    file, err = os.Create(fmt.Sprintf("fs/%d.nonce", request.Id))
    if err != nil {
        log.Fatal(err)
    }
    file.Write(request.File.Nonce[:])
}

func downloadHandler(w http.ResponseWriter, r *http.Request, m *utils.MerkleTree) {
    id := r.URL.Query().Get("id")
    file, err := os.Open(fmt.Sprintf("fs/%s.dat", id))
    if err != nil {
        log.Fatal(err)
    }
    fileInfo, err := file.Stat()
    if err != nil {
        log.Fatal(err)
    }
    fileData := make([]byte, fileInfo.Size())
    file.Read(fileData)
    nonce, err := os.Open(fmt.Sprintf("fs/%s.nonce", id))
    if err != nil {
        log.Fatal(err)
    }
    var nonceData [12]byte
    uint_id,_:= strconv.Atoi(id)
    nonce.Read(nonceData[:])
    response := utils.Response{
        File: utils.File{
            Data:  fileData,
            Nonce: nonceData,
        },
        Hash: m.ValidateHash(uint8(uint_id), sha512.Sum512(fileData)),
    }
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        log.Fatal(err)
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonResponse)
}

func validationHandler(w http.ResponseWriter, r *http.Request, m *utils.MerkleTree) {	
    request := utils.ValidateHashRequest{}
    jsonRequest := json.NewDecoder(r.Body)
    jsonRequest.Decode(&request)
    file, err := os.Open(fmt.Sprintf("fs/%d.dat", request.Id))
    if err != nil {
        log.Fatal(err)
    }
    fileInfo, err := file.Stat()
    if err != nil {
        log.Fatal(err)
    }
    fileData := make([]byte, fileInfo.Size())
    file.Read(fileData)
    response := utils.ValidateHashResponse{
        Hash: m.ValidateHash(request.Id, sha512.Sum512(fileData)),
    }
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        log.Fatal(err)
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonResponse)
}

func addFilesToMerkleTree(m *utils.MerkleTree){
    f, err := os.Open("fs")
    if err != nil {
        fmt.Println(err)
        return
    }
    files, err := f.Readdir(0)
    if err != nil {
        fmt.Println(err)
        return
    }

    for _, v := range files {
	if strings.Split(v.Name(), ".")[1]=="dat" {
		fileInfo, err := f.Stat()
   		 if err != nil {
        		log.Fatal(err)
    		}
   		fileData := make([]byte, fileInfo.Size())
    		f.Read(fileData)
		uint_id,_:= strconv.Atoi(strings.Split(v.Name(), ".")[0])
		m.SetLeaf(uint8(uint_id),sha512.Sum512(fileData))
	}
        fmt.Println(strings.Split(v.Name(), ".")[1], v.IsDir())
    }
    m.CalcHash()
}
func main() {
	//m := utils.MerkleTree{}
	if _, err := os.Stat("fs"); os.IsNotExist(err) {
		os.Mkdir("fs", 0755)
	}
	//ln, err := net.Listen("tcp", ":8000")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("Listening on port 8000")
	//conn, err := ln.Accept()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for {
	//	handleConnection(conn, m)
	//}
	mt := &utils.MerkleTree{}
	mt.Instantiate(8)
	addFilesToMerkleTree(mt)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		uploadHandler(w, r, mt)
	})
	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		downloadHandler(w, r, mt)
	})
	http.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		validationHandler(w, r, mt)
	})

	http.ListenAndServe(":8080", nil)
}
