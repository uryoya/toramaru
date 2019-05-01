package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strconv"
)

type Route struct {
	Host string
	Path string
}

type Toramaru struct {
	Port   int
	Routes []Route
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
			url, err := url.Parse("http://" + arg) // TODO: ngnxみたいなルール作る
			if err != nil || url.Host == "" || url.Path == "" {
				return nil, errors.New("route can not parse")
			}
			toramaru.Routes = append(toramaru.Routes, Route{url.Host, url.Path})

		default:
			return nil, errors.New("invalid options")
		}
	}
	return toramaru, nil
}

func help() string {
	return `USAGE: toramaru [options...]

	OPTIONS:
		-p [PORT]  -- -p 8080
		-r [ROUTE] -- -r "localhost:8070/a/" -r "localhost:8071/b/"
		-h --help  -- show this help
	`
}

func main() {
	toramaru, err := argparse(os.Args)
	switch {
	case toramaru.Help:
		fmt.Print(help())
		os.Exit(0)
	case err != nil:
		fmt.Println(err)
		os.Exit(-1)
	}

	director := func(request *http.Request) {
		log.Printf("%s [%s] %s", request.Proto, request.Method, request.URL.Path)
		request.URL.Scheme = "http"

		for _, route := range toramaru.Routes {
			matched, _ := regexp.MatchString(route.Path+`.*`, request.URL.Path)
			if matched {
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
