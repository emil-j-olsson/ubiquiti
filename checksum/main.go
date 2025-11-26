package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

const Version = "v1.0.0"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("Checksum Generator %s", Version)
		os.Exit(0)
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
	hash := sha256.Sum256(data)
	checksum := hex.EncodeToString(hash[:])
	fmt.Print(checksum)
	os.Exit(0)
}
