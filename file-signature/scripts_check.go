package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ScriptCheck(fileName string) bool {
	file, err := os.Open(fileName)
	if err != nil {
		return false
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#!") && strings.Contains(line, "python") {
			fmt.Printf("Match found: Python Script\n")
			return true
		}
		// Sometimes we will have to do a full scan for the contents or content sniffing example
		if strings.HasPrefix(line, "import ") || strings.HasPrefix(line, "def ") || strings.Contains(line, "print(") {
			fmt.Printf("Match found: Python Script\n")
			return true
		}
		// Can add additonal checks for different file type
	}

	return false
}
