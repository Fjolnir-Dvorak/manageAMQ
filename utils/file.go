package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

func CountLinesFromFilename(filename string) (int, error){
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println("ERROR: Could not open file")
		return 0, err
	}
	defer file.Close()
	return CountLinesFromFile(file)
}

func CountLinesFromFile(file *os.File) (int, error){
	defer file.Seek(0,0)
	file.Seek(0,0)
	reader := bufio.NewReader(file)

	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := reader.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil
		case err != nil:
			return count, err
		}
	}
}