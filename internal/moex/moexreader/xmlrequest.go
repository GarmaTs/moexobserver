package moexreader

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func GetXMLByRequest(url string) (io.Reader, error) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Println("RECOVERED GetXMLByRequest, url:", url)
	// 	}
	// }()

	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body: %v", err)
	}

	return bytes.NewReader(data), nil
}
