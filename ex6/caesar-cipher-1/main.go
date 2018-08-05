package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Complete the caesarCipher function below.
func caesarCipher(s string, k int32) string {
	ret := ""
	for _, c := range s {
		switch {
		case c >= 'a' && c <= 'z':
			r := (c-'a'+k)%26 + 'a'
			ret += string(r)
		case c >= 'A' && c <= 'Z':
			r := (c-'A'+k)%26 + 'A'
			ret += string(r)
		default:
			ret += string(c)
		}
	}
	return ret
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1024*1024)

	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 1024*1024)

	_, err = strconv.ParseInt(readLine(reader), 10, 64)
	checkError(err)

	s := readLine(reader)

	k64, err := strconv.ParseInt(readLine(reader), 10, 64)
	checkError(err)
	k := int32(k64)

	result := caesarCipher(s, k)

	fmt.Fprintf(writer, "%s\n", result)

	writer.Flush()
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
