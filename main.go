// OSS API
//
// Schemes: http
// Host: 0.0.0.0:port
// BasePath: /api/v4/oss
// Version 0.1
// Contact: Christophe Buffard
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta

package main

import (
	"fmt"
	"github.com/Go-routine-4995/ossfrontend/adapter/controllers"
	"github.com/Go-routine-4995/ossfrontend/adapter/gateway"
	"github.com/Go-routine-4995/ossfrontend/adapter/repository/filedb"
	"github.com/Go-routine-4995/ossfrontend/authentication"
	"github.com/Go-routine-4995/ossfrontend/service"
	"gopkg.in/yaml.v2"
	"os"
)

const (
	config  = "conf.yml"
	version = 0.01
)

type Config struct {
	Service struct {
		Port       string `yaml:"port"`
		Host       string `yaml:"host"`
		CertFile   string `yaml:"certfile"`
		KeyFile    string `yaml:"keyfile"`
		CaCertFile string `yaml:"cacertfile"`
	} `yaml:"service"`
	Database struct {
		PubKey string `yaml:"pubKey"`
	} `yaml:"database"`
	Server struct {
		Url string `yaml:"nats"`
	} `yaml:"server"`
}

func main() {

	var (
		conf string
		srv  interface{}
	)

	fmt.Println("Starting OSS front end v", version)
	args := os.Args
	if len(args) < 2 {
		conf = config
	} else {
		conf = args[1]
	}

	cfg := openFile(conf)

	// crate the gateway
	b := gateway.NewBroker(cfg.Server.Url)

	// create a service
	srv = service.NewService(b)

	// create a logger
	// srv = logging.NewLoggingService(srv)

	// create a new storere
	sk := filedb.NewFileDB(cfg.Database.PubKey)
	// create a new authentication
	a := authentication.NewAuthentication(sk)

	api := controllers.NewApiServer(srv, cfg.Service.Port, a)
	api.Start()
}

func openFile(s string) Config {
	f, err := os.Open(s)
	if err != nil {
		processError(err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		processError(err)
	}

	return cfg
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}
