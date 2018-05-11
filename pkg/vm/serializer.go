package vm

import (
	"encoding/json"
	"io"
)

// ChasmBinary defines the "binary" (assembled) format of a vm
type ChasmBinary struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Context string `json:"context"`
	Data    []byte `json:"data"`
}

// Serialize takes a stream of bytes (including the context marker) and sends it to
// a Writer in ChasmBinary format
func Serialize(name string, comment string, b []byte, w io.Writer) error {
	output := ChasmBinary{
		Name:    name,
		Comment: comment,
		Context: Contexts[b[0]],
		Data:    b,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(output)
}

// Deserialize takes a reader and extracts a ChasmBinary from it
func Deserialize(r io.Reader) (ChasmBinary, error) {
	var input ChasmBinary
	err := json.NewDecoder(r).Decode(&input)
	return input, err
}
