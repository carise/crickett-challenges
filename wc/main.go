package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"strings"
)

var (
	countBytes = flag.Bool("c", false, "display number of bytes")
	countLines = flag.Bool("l", false, "display number of lines")
	countWords = flag.Bool("w", false, "display number of words")
)

func main() {
	argsWithoutProg := os.Args[1:]
	log.Printf("args without program: %s\n", &argsWithoutProg)

	if len(argsWithoutProg) > 4 {
		log.Fatal("Too many arguments")
		os.Exit(1)
	}
	filename := argsWithoutProg[len(argsWithoutProg)-1]
	if strings.Compare(filename, "-c") == 0 || strings.Compare(filename, "-l") == 0 || strings.Compare(filename, "-w") == 0 {
		log.Fatal("Missing file name")
		os.Exit(1)
	}
	log.Printf("filename: %s\n", filename)

	flag.Parse()

	if !*countBytes && !*countLines && !*countWords {
		*countBytes = true
		*countLines = true
		*countWords = true
	}
	log.Printf("count bytes: %v\ncount lines: %v\ncount words: %v\n", *countBytes, *countLines, *countWords)

	byteCount, _, _, err := count(countBytes, countLines, countWords, filename)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("	%d	%s\n", byteCount, filename)
}

func count(countBytes *bool, countLines *bool, countWords *bool, filename string) (int, int, int, error) {
	fp, err := os.Open(filename)
	defer fp.Close()
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(fp)
	byteCount := 0
	if *countBytes {
		for {
			_, err := reader.ReadByte()
			if err != nil && err == io.EOF {
				break
			}
			byteCount++
		}
	}

	return byteCount, 0, 0, err
}
