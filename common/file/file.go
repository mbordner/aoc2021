package file

import (
	"errors"
	"os"
	"strings"
)

func GetLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	bytesread, err := file.Read(buffer)
	if err != nil {
		return nil, err
	}

	if bytesread != int(filesize) {
		return nil, errors.New("didn't read all of the file")
	}

	rows := strings.Split(string(buffer), "\n")

	return rows, nil
}
