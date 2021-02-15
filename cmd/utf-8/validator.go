package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"unicode/utf8"
)

// This program aims to check a file if it's encoded properly in UTF-8 or not.

// An encoding is invalid if it is incorrect UTF-8, encodes a rune that is out of range,
// or is not the shortest possible UTF-8 encoding for the value. No other validation is
// performed.

var (
	filePath   string
	vMeta      bool
	runeMapper bool
)

const (
	usageFilePath   = "<string>: mention filename"
	usageVMeta      = "<bool>: enable verbose offset mode to print line, line number, offset in case of error (first error)"
	usageRuneMapper = "<bool>: enable mapper mode to print occurrence of every character (rune) on successful parsing"
)

const (
	defaultFilePath   = ""
	defaultVMeta      = false
	defaultRuneMapper = false
)

func fatalln(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func basicValidator(file *os.File) error {
	var mp map[string]int64
	if runeMapper {
		mp = make(map[string]int64)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)

	for scanner.Scan() {
		if runeMapper {
			mp[scanner.Text()]++
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if runeMapper {
		for k, v := range mp {
			fmt.Println("Rune: ", k, ", Count: ", v)
		}
	}
	return nil
}

func advancedValidator(file *os.File) error {
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	var mp map[string]int64
	if runeMapper {
		mp = make(map[string]int64)
	}

	debugLMap := make(map[int64]bool)
	size := 0
	for start := 0; start < len(buf); start += size {
		var r rune
		if r, size = utf8.DecodeRune(buf[start:]); r == utf8.RuneError && size == 1 {
			line, lineNo, err := byteLineFinder(file, int64(start))
			if err != nil {
				return fmt.Errorf("Invalid UTF-8 encoding at offset: %d (bytes)", start)
			}
			return fmt.Errorf("Invalid UTF-8 encoding at:\nByteOffset: %d\nLine Number: %d\nLine: %s", start, lineNo, line)
		}

		if r == utf8.RuneError {
			line, lineNo, _ := byteLineFinder(file, int64(start))
			if _, ok := debugLMap[lineNo]; !ok {
				debugLMap[lineNo] = true
				fmt.Printf("\nLine no: %d\n", lineNo)
				fmt.Printf("Line: %s\n", line)
			}
		}

		if runeMapper {
			mp[string(r)]++
		}
	}

	if runeMapper {
		for k, v := range mp {
			fmt.Println("Rune: ", k, ", Count: ", v)
		}
	}
	return nil
}

func byteLineFinder(file *os.File, find int64) (string, int64, error) {
	if _, err := file.Seek(0, 0); err != nil {
		return "", 0, err
	}

	scanner := bufio.NewScanner(file)
	line := int64(1)
	bytesRead := int64(0)

	for scanner.Scan() {
		b := scanner.Text()
		offset := bytesRead + int64(len(b))
		if bytesRead <= find && find <= offset {
			return b, line, nil
		}
		bytesRead = offset + 1
		line++
	}

	if err := scanner.Err(); err != nil {
		return "", 0, err
	}
	return "", line, nil
}

func main() {
	flag.StringVar(&filePath, "file", defaultFilePath, usageFilePath)
	flag.BoolVar(&vMeta, "v", defaultVMeta, usageVMeta)
	flag.BoolVar(&runeMapper, "m", defaultRuneMapper, usageRuneMapper)
	flag.Parse()

	if filePath == "" {
		log.Fatalln(flag.ErrHelp.Error() + ". Try (-h) or (--help) flag")
	}

	file, err := os.Open(filePath)
	fatalln(err)

	switch vMeta {
	case true:
		err := advancedValidator(file)
		fatalln(err)
	case false:
		err := basicValidator(file)
		fatalln(err)
	}
}
