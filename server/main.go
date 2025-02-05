package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
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


func handleConnection(conn net.Conn) {
    defer conn.Close()
    for {
        handleDownload(conn)
    }
}

func main() {
    ln, err := net.Listen("tcp", ":8000")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Listening on port 8000")
    conn, err := ln.Accept()
    if err != nil {
        log.Fatal(err)
    }
    for {
        handleConnection(conn)
    }
}
