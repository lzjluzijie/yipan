package onedrive

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type File struct {
	ID   string
	Name string
	Path string
	Size int64
	URL  string
}

type Child struct {
	ID   string
	Name string
	Size int64
}

type ListChildrenResponse struct {
	Children []Child `json:"value"`
}

func ListChildren(id, path string) (err error) {
	req, err := NewRequest("GET", fmt.Sprintf("https://graph.microsoft.com/v1.0/me/drive/items/%s/children", id), nil)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	listChildrenResponse := &ListChildrenResponse{}
	err = json.Unmarshal(data, listChildrenResponse)
	if err != nil {
		return
	}

	log.Println(listChildrenResponse)
	return
}
