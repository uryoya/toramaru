package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"

	"github.com/uryoya/toramaru/route"
)

const help string = `USAGE: toramaru [options...]

  OPTIONS:
    -p --port        [PORT]  -- -p 8080
    -r --route_proxy [ROUTE_PROXY]
    -h --help  -- show this help

  ROUTE_PROXY: [LOCATION]>[HOST]
    example:
    "/path/to/location>localhost:8070"
`

type Toramaru struct {
	Port   int
	Routes []route.Route
	Help   bool
}

const EOA = "__EOA__" // end of args
func argparse(args []string) (toramaru *Toramaru, err error) {
	toramaru = &Toramaru{Help: false}
	args = append(args, EOA)
	for i := 1; i < len(args)-1; i += 2 {
		opt := args[i]
		arg := args[i+1]
		switch {
		case opt == "-h" || opt == "--help":
			toramaru.Help = true

		case opt == "-p" && arg != EOA:
			toramaru.Port, err = strconv.Atoi(arg)
			if err != nil {
				return nil, errors.New("port can not convert to int")
			}

		case opt == "-r" && arg != EOA:
			r, err := route.Parse(arg)
			if err != nil {
				return nil, err
			}
			toramaru.Routes = append(toramaru.Routes, *r)

		default:
			return nil, errors.New("invalid options")
		}
	}
	return toramaru, nil
}

func main() {
	toramaru, err := argparse(os.Args)
	switch {
	case err != nil:
		fmt.Println(err)
		os.Exit(-1)
	case toramaru.Help:
		fmt.Print(help)
		os.Exit(0)
	}

	director := func(request *http.Request) {
		log.Printf("%s [%s] %s", request.Proto, request.Method, request.URL.Path)
		request.URL.Scheme = "http"

		for _, route := range toramaru.Routes {
			if route.Match(request.URL.Path) {
				request.URL.Host = route.Host
				break
			}
		}
	}

	rp := &httputil.ReverseProxy{Director: director}
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", toramaru.Port),
		Handler: rp,
	}

	log.Printf("toramaru running on: %d\n", toramaru.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
