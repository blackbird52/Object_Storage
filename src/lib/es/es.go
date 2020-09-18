package es

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"os"
	"strings"
)

type Metadata struct {
	Name    string
	Version int
	Size    int64
	Hash    string
}

func getMetadata(name string, versionId int) (meta Metadata, err error) {
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d/_source", os.Getenv("ES_SERVER"), name, versionId)
	response, err := http.Get(url)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("fail to get %s_%d: %d", name, versionId, response.StatusCode)
		return
	}
	result, _ := ioutil.ReadAll(response.Body)
	_ = json.Unmarshal(result, &meta)
	return
}

type hit struct {
	Source Metadata `json:"_source"`
}

type searchResult struct {
	Hits struct {
		Total int
		Hits  []hit
	}
}

func SearchLatestVersion(name string) (meta Metadata, err error) {
	url := fmt.Sprintf("http://%s/metadata/_search?q=name:%s&size=1&sort=version:desc",
		os.Getenv("ES_SERVER"), url2.PathEscape(name))
	response, err := http.Get(url)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("fail to search last metadata: %d", response.StatusCode)
		return
	}
	result, _ := ioutil.ReadAll(response.Body)
	var sr searchResult
	_ = json.Unmarshal(result, &sr)
	if len(sr.Hits.Hits) != 0 {
		meta = sr.Hits.Hits[0].Source
	}
	return
}

func GetMetadata(name string, version int) (Metadata, error) {
	if version == 0 {
		return SearchLatestVersion(name)
	}
	return getMetadata(name, version)
}

func PutMetadata(name string, version int, size int64, hash string) error {
	doc := fmt.Sprintf(`{"name":"%s","version":%d,"size":%d,"hash":"%s"}`, name, version, size, hash)
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d?op_type=create", os.Getenv("ES_SERVER"), name, version)
	request, _ := http.NewRequest("PUT", url, strings.NewReader(doc))
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode == http.StatusConflict {
		return PutMetadata(name, version+1, size, hash)
	}
	if response.StatusCode != http.StatusCreated {
		result, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf("fail to put metadata: %d %s", response.StatusCode, string(result))
	}
	return nil
}

func AddVersion(name string, hash string, size int64) error {
	meta, err := SearchLatestVersion(name)
	if err != nil {
		return err
	}
	return PutMetadata(name, meta.Version+1, size, hash)
}

func SearchAllVersions(name string, from int, size int) ([]Metadata, error) {
	url := fmt.Sprintf("http://%s/metadata/_search?sort=name,version&from=%d&size=%d",
		os.Getenv("ES_SERVER"), from, size)
	if name != "" {
		url += "&q=name:" + name
	}
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	metas := make([]Metadata, 0)
	result, _ := ioutil.ReadAll(response.Body)
	var sr searchResult
	_ = json.Unmarshal(result, &sr)
	for _, hit := range sr.Hits.Hits {
		metas = append(metas, hit.Source)
	}
	return metas, nil
}
