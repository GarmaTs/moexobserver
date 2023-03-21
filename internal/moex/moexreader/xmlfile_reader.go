package moexreader

import (
	"os"
)

func GetXMLFromFile(filename string) (*os.File, error) {
	reader, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return reader, nil
}
