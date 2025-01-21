package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/igadmg/goex/gx"
)

var (
	path_f = flag.String("path", "", "worikning dir")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of el:\n")
	fmt.Fprintf(os.Stderr, "\tel [flags]\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

var state int = StateStream
var firstStreamToken bool = true
var consumeNewLine bool = true

const (
	StateStream = iota
	StateCode
)

var state_end_tokens []string = []string{
	"<?",
	"?>",
}

func ReadStringToken(b *bufio.Reader, token string) (r string, err error) {
	bt := []byte(token)
	var br []byte

	for {
		br, err = b.ReadSlice(token[0])
		if errors.Is(err, io.EOF) {
			r += string(br)
			if len(r) > 0 {
				err = nil
			}
			return
		}
		if err != nil {
			return
		}

		var tok []byte
		tl := len(token) - 1
		if tl > 0 {
			tok, err = b.Peek(tl)
			if errors.Is(err, io.EOF) {
				r += string(br)
				if len(r) > 0 {
					err = nil
				}
				return
			}
			if err != nil {
				return
			}
			if !bytes.Equal(tok, bt[1:]) {
				br = append(br, gx.Must(b.ReadByte()))
				r += string(br)
				continue
			}
		}

		b.Discard(tl)
		r += string(br[:len(br)-1])
		return
	}
}

func read_files(dir string) (err error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, f := range files {
		if f.IsDir() {
			err = read_files(filepath.Join(dir, f.Name()))
			if err != nil {
				return
			}
		}

		if !strings.HasSuffix(f.Name(), ".el") {
			continue
		}

		filePath := filepath.Join(dir, f.Name())
		buffer := process_file(filePath)

		outputFilePath := filepath.Join(dir, strings.TrimSuffix(f.Name(), ".el"))
		outputFile := gx.Must(os.Create(outputFilePath))
		defer outputFile.Close()

		src, err := format.Source(buffer.Bytes())
		if err != nil {
			log.Printf("warning: internal error: invalid Go generated: %s", err)
			log.Printf("warning: compile the package to analyze the error")
			outputFile.Write(buffer.Bytes())
		}

		outputFile.Write(src)
	}

	return
}

func process_file(filePath string) bytes.Buffer {
	file := gx.Must(os.Open(filePath))
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.Buffer{}
	consumeNewLine = true
	firstStreamToken = true

	for {
		str, err := ReadStringToken(reader, state_end_tokens[state])
		if err != nil {
			break
		}

		switch state {
		case StateStream:
			if len(str) != 0 {
				if consumeNewLine {
					if str[0] == '\n' {
						str = str[1:]
						if str[0] == '\r' {
							str = str[1:]
						}
					} else if str[0] == '\r' {
						str = str[1:]
						if str[0] == '\n' {
							str = str[1:]
						}
					}
				}
				buffer.WriteString("wr.Write([]byte(`")
				str = strings.Join(strings.Split(str, "`"), "` + \"`\" + `")
				buffer.WriteString(str)
				buffer.WriteString("`))\n")
				consumeNewLine = str[len(str)-1] == '\n'
			} else {
				consumeNewLine = false
			}
			state = StateCode
		case StateCode:
			if firstStreamToken {
				prefix := ""
				for i, ch := range str {
					if !unicode.IsLetter(ch) {
						if unicode.IsSpace(ch) {
							str = str[i+1:]
						}
						break
					}
					prefix += string(ch)
				}
				firstStreamToken = false
			}
			if strings.HasPrefix(str, "=") {
				buffer.WriteString("wr.Write([]byte(fmt.Sprintf(\"%v\", ")
				buffer.WriteString(str[1:])
				buffer.WriteString(")))\n")
			} else {
				buffer.WriteString(str)
				buffer.WriteString("\n")
			}
			state = StateStream
		}
	}
	return buffer
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("el: ")
	flag.Usage = Usage
	flag.Parse()

	dir := "."
	err := read_files(dir)
	if err != nil {
		panic(err)
	}
}
