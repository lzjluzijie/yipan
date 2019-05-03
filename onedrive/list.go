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
	ID     string
	Name   string
	Size   int64
	Folder Folder
}

type Folder struct {
	ChildCount int
}

type ListChildrenResponse struct {
	Children []Child `json:"value"`
}

func ListChildren(id, path string) (files []File, err error) {
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

	files = make([]File, 0)
	for _, child := range listChildrenResponse.Children {
		if child.Folder.ChildCount != 0 {
			fs, err := ListChildren(child.ID, path+"/"+child.Name)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			files = append(files, fs...)
		} else {
			url, err := Share(child.ID)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			files = append(files, File{
				ID:   child.ID,
				Name: child.Name,
				Path: path + "/" + child.Name,
				Size: child.Size,
				URL:  url,
			})
		}
	}
	return
}
