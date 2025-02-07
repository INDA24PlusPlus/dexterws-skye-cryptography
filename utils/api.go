package utils

type File struct {
	Data  []byte
	Nonce [12]byte
}

type Response struct {
	File File
	Hash b512
	// Required proof nodes
	// Root hash
}

type UploadRequest struct {
    File File
    Id  uint8
}

type ValidateHashRequest struct {
    Id  uint8
}

type UpdateHashPayload struct {
	NewHash  b512 //Xored with old Hash
	Checksum b512 // Sha256(OldHash+NewHash)
}

type ValidateHashResponse struct {
    Hash b512
}

