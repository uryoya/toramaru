package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

const version string = "0.1.0-SNAPSHOT"
const help string = `USAGE: toramaru [OPTIONS...]

  OPTIONS:
    -p, --port [PORT]                使用するポートを指定
    -r, --route-proxy [ROUTE_PROXY]  リバースプロキシを指定

    -h, --help                       このヘルプを表示
    -v, --version                    バージョンを表示

    ROUTE_PROXY: [LOCATION]>[HOST]

  example:
		toramaru -p 8080 -r "/api>localhost:8071" -r "/>localhost:8070"
`

func main() {
	toramaru, err := argparse(os.Args)
	switch {
	case err != nil:
		fmt.Println(err)
		os.Exit(-1)
	case toramaru.Help:
		fmt.Print(help)
		os.Exit(0)
	case toramaru.Version:
		fmt.Println(version)
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

	fmt.Println("proxies:")
	for _, route := range toramaru.Routes {
		fmt.Printf("%s => %s\n", route.Location, route.Host)
	}
	fmt.Printf("toramaru running on: %d\n", toramaru.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
