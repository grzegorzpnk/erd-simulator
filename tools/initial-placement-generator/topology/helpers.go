package topology

import (
	log "10.254.188.33/matyspi5/erd/initial-placement-generator/logger"
	"fmt"
	"io/ioutil"
	"net/http"
)

// func getHTTPRespBody(url string) (io.ReadCloser, error) {
func getHTTPRespBody(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		err := fmt.Errorf("HTTP GET failed for URL %s.\nError: %s\n", url, err)
		log.Errorf("%v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("HTTP GET returned status code %s for URL %s.\n", resp.Status, url)
		log.Errorf("%v", err)
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	return b, nil
}
