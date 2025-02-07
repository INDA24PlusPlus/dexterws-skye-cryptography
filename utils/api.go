package utils

type File struct {
	Data  []byte
	Nonce [12]byte
}

type Response struct {
	File File
	Hash u512
	// Required proof nodes
	// Root hash
}

type UploadRequest struct {
    File File
    Id  uint8
}
