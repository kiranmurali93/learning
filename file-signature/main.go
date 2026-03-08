package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: sig <filename>")
		return
	}

	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	fmt.Println("File read")
	fmt.Println("File: ", fileName)

	found := false
	for _, s := range FileType {
		header := make([]byte, len(s.Signature))
		_, err := file.ReadAt(header, int64(s.Offset))
		if err != nil {
			continue
		}
		if s.CustomCheck != nil {
			matched := s.CustomCheck(fileName)
			if matched {
				found = true
				break
			}

		}

		if bytes.HasPrefix(header, s.Signature) {
			fmt.Printf("Match found: %s\n", s.Name)
			found = true
			break
		}
	}

	if !found {
		fmt.Println("Unknown file type")
	}
}
