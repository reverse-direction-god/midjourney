package until

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func Post(url string, s []byte) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "*", bytes.NewBuffer(s))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	bys, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bys, nil
}
