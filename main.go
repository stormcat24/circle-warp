package main
import (
	"github.com/stormcat24/circle-warp/server"
	"github.com/stormcat24/circle-warp/config"
	"path"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
//	"flag"
)

func main() {

	base := path.Base("circle-warp")

	app := kingpin.New(base, "circle-warp")
	app.Version("0.1.0")

	conf := config.Config{}
//	app.Flag("port", "circle-warp listen port").
//		Default(":7777").
//		Short('p').
//		OverrideDefaultFromEnvar("WARP_PORT").
//		PlaceHolder("WARP_PORT").
//		StringVar(&conf.WarpPort)

	app.Flag("circleci-token", "CircleCI access token").
		Short('t').
		Required().
		OverrideDefaultFromEnvar("CIRCLECI_TOKEN").
		PlaceHolder("CIRCLECI_TOKEN").
		StringVar(&conf.CircleciToken)

	app.Flag("circleci-host", "CircleCI host name").
		Default("circleci.com").
		Short('h').
		OverrideDefaultFromEnvar("CIRCLECI_HOST").
		PlaceHolder("CIRCLECI_HOST").
		StringVar(&conf.CircleciHost)

	kingpin.MustParse(app.Parse(os.Args[1:]))

//	flag.Set("bind", conf.WarpPort)

	server.InitApiClient(&conf)
	server.Serve(&conf)
}

