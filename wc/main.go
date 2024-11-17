package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
)

var (
	countBytes = flag.Bool("c", false, "display number of bytes")
	countLines = flag.Bool("l", false, "display number of lines")
	countWords = flag.Bool("w", false, "display number of words")
)

func main() {
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) > 4 {
		fmt.Printf("Too many arguments")
		os.Exit(1)
	}
	filename := argsWithoutProg[len(argsWithoutProg)-1]
	if strings.Compare(filename, "-c") == 0 || strings.Compare(filename, "-l") == 0 || strings.Compare(filename, "-w") == 0 {
		log.Fatal("Missing file name")
		os.Exit(1)
	}

	flag.Parse()

	if !*countBytes && !*countLines && !*countWords {
		*countBytes = true
		*countLines = true
		*countWords = true
	}

	byteCount, lineCount, wordCount, err := count(countBytes, countLines, countWords, filename)
	if err != nil {
		log.Fatal(err)
	}
	if *countLines {
		fmt.Printf("  %d", lineCount)
	}
	if *countWords {
		fmt.Printf("  %d", wordCount)
	}
	if *countBytes {
		fmt.Printf("  %d", byteCount)
	}
	fmt.Printf(" %s\n", filename)
}

func count(countBytes *bool, countLines *bool, countWords *bool, filename string) (int, int, int, error) {
	fp, err := os.Open(filename)
	defer fp.Close()
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(fp)
	var previousByte byte
	byteCount := 0
	lineCount := 0
	wordCount := 0
	for {
		currentByte, err := reader.ReadByte()
		if err != nil && err == io.EOF {
			break
		}
		if *countBytes {
			byteCount++
		}
		if *countLines {
			foundLineBreak := isLineBreak(currentByte, previousByte)
			if foundLineBreak {
				lineCount++
			}
		}
		if *countWords {
			if !unicode.IsSpace(rune(previousByte)) && unicode.IsSpace(rune(currentByte)) {
				wordCount++
			}
		}

		previousByte = currentByte
	}

	return byteCount, lineCount, wordCount, err
}

func isLineBreak(currentByte byte, previousByte byte) bool {
	// don't double-count if we are dealing with CRLF
	if previousByte == '\r' && currentByte == '\n' {
		return true
	}
	if previousByte != '\r' && currentByte == '\n' {
		return true
	}
	if previousByte == '\r' && currentByte != '\n' {
		return true
	}
	return false
}
