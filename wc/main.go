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
	"github.com/jeandeaual/go-locale"
)

type WcArguments struct {
	Filename string
	Bytes bool
	Lines bool
	Words bool
	Chars bool
}

var (
	countBytes = flag.Bool("c", false, "display number of bytes")
	countLines = flag.Bool("l", false, "display number of lines")
	countWords = flag.Bool("w", false, "display number of words")
	countChars = flag.Bool("m", false, "display number of characters")
)

func main() {
	arguments := parseArguments()

	flag.Parse()

	if !*countBytes && !*countLines && !*countWords && !*countChars {
		*countBytes = true
		*countLines = true
		*countWords = true
	}

	byteCount, lineCount, wordCount, charCount, err := count(countBytes, countLines, countWords, countChars, arguments.Filename)
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
	if *countChars {
		fmt.Printf("  %d", charCount)
	}
	fmt.Printf(" %s\n", arguments.Filename)
}

func parseArguments() *WcArguments {
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
	// TODO: 
	//   * enable flags in WcArguments
	//   * if -c and -m are enabled, whichever is the latest (right-most) one supersedes the other
	//   * get rid of flags package use
	//   * determine locale

	userLocales, _ := locale.GetLocales()
	fmt.Println("Locales: %v", userLocales)
	return &WcArguments{
		Filename: filename,
	}
}

func count(countBytes *bool, countLines *bool, countWords *bool, countChars *bool, filename string) (int, int, int, int, error) {
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
	charCount := 0
	for {
		currentByte, err := reader.ReadByte() // or reader.ReadRune() if counting chars
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
		if *countChars {
		}

		previousByte = currentByte
	}

	return byteCount, lineCount, wordCount, charCount, err
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
