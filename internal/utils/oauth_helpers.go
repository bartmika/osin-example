package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/openshift/osin"
)

func DownloadAccessToken(url string, auth *osin.BasicAuth, output map[string]interface{}) error {
	// download access token
	preq, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	if auth != nil {
		preq.SetBasicAuth(auth.Username, auth.Password)
	}

	pclient := &http.Client{}
	presp, err := pclient.Do(preq)
	if err != nil {
		return err
	}
	defer presp.Body.Close()

	if presp.StatusCode != 200 {
		log.Println("DownloadAccessToken|err|res:", presp)
		return errors.New("Invalid status code")
	}

	jdec := json.NewDecoder(presp.Body)
	err = jdec.Decode(&output)
	return err
}
