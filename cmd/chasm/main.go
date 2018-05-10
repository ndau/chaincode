package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	in := os.Stdin
	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		in = f
	}

	var buf bytes.Buffer
	tee := io.TeeReader(in, &buf)

	sn, err := ParseReader("", tee)
	if err != nil {
		log.Fatal(describeErrors(err, buf.String()))
	}
	b := sn.(Script).bytes()
	fmt.Println(hex.Dump(b))
}
