package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// APIENDPOINT is file in currentdirectory where simplemediconverter looks for the API endpoint for slack
// You need to get this slack app configured inside your workspace and the URL you get is a secret
// Learn more: https://api.slack.com/messaging/webhooks
const APIENDPOINT = "slackapiendpoint.config"

type slackclient struct {
	httpclient  *http.Client
	apiendpoint string
}

func (cli *slackclient) prepareClient() {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}

	filecontent, err := ioutil.ReadFile(APIENDPOINT)
	if err != nil {
		log.Printf("Failed to read slack file: %s\n", APIENDPOINT)
		log.Printf("Error : %s\n", err.Error())
	}

	_, err = url.Parse(string(filecontent))
	if err != nil {
		log.Printf("Failed to extract URL..check contents of %s to see if it is real URL\n", APIENDPOINT)
	}

	cli.apiendpoint = strings.TrimRight(string(filecontent), "\n")
	cli.httpclient = &http.Client{Transport: tr, Timeout: 30 * time.Second}
}

func (cli *slackclient) connect() error {
	return nil
}

func performHTTPPost(cli *slackclient, message string) error {
	requestBody, err := json.Marshal(map[string]string{"text": message})

	if err != nil {
		log.Printf("\nperformHTTPPost problem with json marshaling\n")
		return err
	}

	resp, err := cli.httpclient.Post(cli.apiendpoint,
		"application/json",
		bytes.NewBuffer(requestBody))

	if err != nil {
		log.Printf("\nperformHTTPPost problem with httpost: %s\n", err.Error())
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Unable to read Response\n")
			return err
		}

		log.Printf("Response : %s\n", string(body))
		return err
	}

	log.Printf("Message : %s was successfully delivered\n", message)

	return nil
}

func (cli *slackclient) notifyMessage(message string) error {
	log.Printf("\nNotifyMessage called with %s\n", message)

	performHTTPPost(cli, message)

	return nil
}

func (cli *slackclient) disConnect() error {
	return nil
}
