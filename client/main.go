package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"dws-sk-fs/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

func uploadFile() error {
    var file string
    fmt.Print("Enter file path: ")
    fmt.Scanln(&file)
    f, err := os.Open(file)
    if err != nil {
        fmt.Println(err)
        return errors.New("File not found")
    }
    defer f.Close()
    var password string
    fmt.Print("Enter encryption password: ")
    fmt.Scanln(&password)
    key := []byte(password)
    fileInfo, err := f.Stat()
    if err != nil {
        fmt.Println(err)
        return errors.New("File not found")
    }
    data := make([]byte, fileInfo.Size())
    _, err = f.Read(data)
    if err != nil {
        fmt.Println(err)
        return errors.New("File not found")
    }
    encrypted_file, err := encryptFile(data, key)
    if err != nil {
        fmt.Println(err)
        return errors.New("Encryption failed")
    }
    var id uint8
    fmt.Print("Enter file id: ")
    fmt.Scanln(&id)
    uploadRequest := utils.UploadRequest {
        File: encrypted_file,
        Id: id,
    }
    // Send uploadRequest to server
    json_uploadRequest, err := json.Marshal(uploadRequest)
    if err != nil {
        fmt.Println(err)
        return errors.New("Failed to reach server")
    }
    http.Post("http://localhost:8080/upload", "application/json", bytes.NewBuffer(json_uploadRequest))
    return nil
}

func encryptFile(data []byte, key []byte) (utils.File, error) {
    var nonce [12]byte
    _, err := rand.Read(nonce[:])
    if err != nil {
        return utils.File{}, err
    }
    key_aes := make([]byte, 32)
    copy(key_aes, key)
    aes, err := aes.NewCipher(key_aes)
    if err != nil {
        return utils.File{}, err
    }
    gcm, err := cipher.NewGCM(aes)
    if err != nil {
        return utils.File{}, err
    }
    ciphertext := gcm.Seal(nil, nonce[:], data, nil)
    return utils.File {
        Data: ciphertext,
        Nonce: nonce,
    }, nil
}

func downloadFile() error {
    var id uint8
    fmt.Print("Enter file id: ")
    fmt.Scanln(&id)
    var password string
    fmt.Print("Enter password: ")
    fmt.Scanln(&password)
    key := []byte(password)
    response, err := http.Get("http://localhost:8080/download?id=" + string(id))
    if err != nil {
        return errors.New("Failed to reach server")
    }
    var downloadResponse utils.Response
    jsonDownloadResponse := json.NewDecoder(response.Body)
    jsonDownloadResponse.Decode(&downloadResponse)
    fmt.Println("File downloaded: " + string(downloadResponse.File.Data))
    decrypted_file, err := decryptFile(downloadResponse.File, key)
    if err != nil {
        return errors.New("Decryption failed")
    }
    file, err := os.Create("downloaded.dat")
    if err != nil {
        fmt.Println(err)
        return errors.New("Failed to create file")
    }
    file.Write(decrypted_file)
    return nil
}

func decryptFile(file utils.File, key []byte) ([]byte, error) {
    key_aes := make([]byte, 32)
    copy(key_aes, key)
    aes, err := aes.NewCipher(key_aes)
    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(aes)
    if err != nil {
        return nil, err
    }
    data, err := gcm.Open(nil, file.Nonce[:], file.Data, nil)
    if err != nil {
        return nil, err
    }
    fmt.Println("Decrypted data: " + string(data))
    return data, nil
}


func commandLoop() {
    for {
        fmt.Print("Enter command: ")
        var command string
        fmt.Scanln(&command)
        switch command {
        case "upload":
            uploadFile()
        case "download":
            downloadFile()
        case "exit":
            return
        default:
            fmt.Println("Invalid command")
            fmt.Println("Commands: upload, download, exit")
        }
    }
}

func main() {
    merkle := utils.MerkleTree{}
    merkle.Instantiate(8)
    commandLoop()
}
