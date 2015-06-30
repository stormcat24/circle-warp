package server

import (
	"github.com/zenazn/goji/web"
	"net/http"
	"github.com/stormcat24/circle-warp/circleci"
	"fmt"
	"strings"
	"io/ioutil"
	"github.com/stormcat24/circle-warp/config"
)

var client circleci.ApiClient

func InitApiClient(conf *config.Config) {
	client = *circleci.NewApiClient(conf)
}

func Rewrite(c web.C, w http.ResponseWriter, r *http.Request) {

	organization := c.URLParams["org"]
	reponame := c.URLParams["repo"]

	repoFullname := fmt.Sprintf("%s/%s", organization, reponame)
	resourcePath := strings.Replace(r.RequestURI, fmt.Sprintf("/%s", repoFullname), "", 1)

	build, err := client.GetLatestSuccessBuild(repoFullname, "master")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Resource Not Found"))
		return
	}

	artifacts, err := client.GetArtifacts(repoFullname, build.BuildNum)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Resource Not Found"))
		return
	}

	resources := map[string]circleci.BuildArtifact{}

	for _, a := range *artifacts {
		if strings.HasPrefix(a.PrettyPath, "$CIRCLE_ARTIFACTS/") {
			key := strings.Replace(a.PrettyPath, "$CIRCLE_ARTIFACTS", "", 1)
			resources[key] = a
		}
	}


	if target, ok := resources[resourcePath]; ok {
		resp, _ := http.Get(target.Url)
		defer resp.Body.Close()

		content, _ := ioutil.ReadAll(resp.Body)
		ct := resp.Header.Get("Content-Type")
		if len(ct) > 0 {
			w.Header().Add("Content-Type", ct)
		}

		w.Write([]byte(content))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Resource Not Found"))
	}

}