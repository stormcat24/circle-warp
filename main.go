package main
import (
	"github.com/stormcat24/circle-warp/server"
	"github.com/stormcat24/circle-warp/config"
	"path"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

func main() {

	base := path.Base("circle-warp")

	app := kingpin.New(base, "circle-warp")
	app.Version("0.0.1")

	conf := &config.Config{}
	flag(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))

	server.Serve(conf)
}

func flag(app *kingpin.Application) {

}
