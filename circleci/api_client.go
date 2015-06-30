package circleci

import (
	"net/http"
	"bytes"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/stormcat24/circle-warp/config"
	"errors"
)

type ApiClient struct {
	ApiHost *string
	Token *string
	client *http.Client
}

type ProjectBuild struct {
	BuildUrl string `json:"build_url"`
	Branch   string `json:"branch"`
	Username string `json:"username"`
	Reponame string `json:"reponame"`
	Status   string `json:"status"`
	BuildNum uint32 `json:"build_num"`
}

type BuildArtifact struct {
	Url        string `json:"url"`
	NodeIndex  int    `json:"node_index"`
	PrettyPath string `json:"pretty_path"`
	Path       string `json:"path"`
}

func NewApiClient(conf *config.Config) *ApiClient {

	instance := &ApiClient{
		ApiHost: &conf.CircleciHost,
		Token: &conf.CircleciToken,
		client: &http.Client{},
	}

	return instance
}

func (self *ApiClient) GetLatestSuccessBuild(repo string, branch string) (*ProjectBuild, error) {

	var endpoint string
	if len(branch) == 0 || branch == "master" {
		endpoint = fmt.Sprintf("/project/%s", repo)
	} else {
		endpoint = fmt.Sprintf("/project/%s/tree/%s", repo, branch)
	}

	params := &map[string]string{
		"limit": "1",
		"filter": "success",
	}

	req, err := self.newRequest("GET", endpoint, params)
	if err != nil {
		return nil, err
	}


	resp, err := self.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.New(fmt.Sprintf("api error: status code %d", resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)

	var builds []ProjectBuild
	json.Unmarshal(body, &builds)

	if len(builds) > 0 {
		return &builds[0], nil
	} else {
		return nil, nil
	}

}

func (self *ApiClient) GetArtifacts(repo string, buildNum uint32) (*[]BuildArtifact, error){

	endpoint := fmt.Sprintf("/project/%s/%d/artifacts", repo, buildNum)

	req, err := self.newRequest("GET", endpoint, &map[string]string{})

	resp, err := self.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.New(fmt.Sprintf("api error: status code %d", resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)

	var artifacts []BuildArtifact

	errJson := json.Unmarshal([]byte(body), &artifacts)

	return &artifacts, errJson
}

func (self *ApiClient) newRequest(method string, api string, param *map[string]string) (*http.Request, error) {

	var url bytes.Buffer

	url.WriteString("https://")
	url.WriteString(*self.ApiHost)
	url.WriteString("/api/v1")
	url.WriteString(api)
	url.WriteString("?circle-token=")
	url.WriteString(*self.Token)

	for name, value := range *param {
		url.WriteString("&")
		url.WriteString(name)
		url.WriteString("=")
		url.WriteString(value)
	}

	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")

	return req, nil
}