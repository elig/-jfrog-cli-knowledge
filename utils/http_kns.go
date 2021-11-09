package utils

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"io/ioutil"
	"net/http"
)

const (
	KNOW_SVC_ADDR = "https://knscli.jfrog.org"
)

type KnowResult struct {
	PostID      int64  `json:"post_id"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	PublishDate string `json:"publish_date"`
	ContentType string `json:"content_type"`
}

type KnowContent struct {
	KnowResult
	Content string `json:"content"`
}

func GetFacetsContent(endpoint string) (map[string]int, error) {

	var dat map[string]int

	resp, err := http.Get(KNOW_SVC_ADDR + endpoint)
	defer resp.Body.Close()

	if err != nil {
		log.Info(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Info(err)
	}

	if err := json.Unmarshal(body, &dat); err != nil {
		return nil, err
	}
	return dat, nil
}

func GetResultsContent(endpoint string) ([]KnowResult, error) {

	var results []KnowResult

	resp, err := http.Get(KNOW_SVC_ADDR + endpoint)
	defer resp.Body.Close()
	if err != nil {
		log.Info(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Info(err)
	}

	if err := json.Unmarshal(body, &results); err != nil {
		panic(err)
	}
	return results, nil
}

func GetContent(endpoint string) (KnowContent, error) {

	var content KnowContent

	resp, err := http.Get(KNOW_SVC_ADDR + endpoint)
	defer resp.Body.Close()
	if err != nil {
		log.Error(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	if err := json.Unmarshal(body, &content); err != nil {
		log.Error(err)
	}
	return content, nil
}
