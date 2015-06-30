package server

import (
	"github.com/stormcat24/circle-warp/config"
	"github.com/zenazn/goji"
)

func Serve(conf *config.Config) {

	goji.Get("/:org/:repo/*", Rewrite)
	goji.Serve()
}