package utils

type File struct {
	Data  []byte
	Nonce uint64
}

type Response struct {
	File File
	Hash u512
	// Required proof nodes
	// Root hash
}
