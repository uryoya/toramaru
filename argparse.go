package main

import (
	"errors"
	"strconv"

	"github.com/uryoya/toramaru/route"
)

type Toramaru struct {
	Port    int
	Routes  []route.Route
	Help    bool
	Version bool
}

const EOA = "__EOA__" // end of args
func argparse(args []string) (toramaru *Toramaru, err error) {
	// 優先するオプション
	switch {
	case isHelp(args):
		return &Toramaru{Help: true}, nil
	case isVersion(args):
		return &Toramaru{Version: true}, nil
	}

	// 引数が必要なオプション
	// 今現在起動に必用なオプションは引数が必要なもののみなのでEOAを置いて対処
	toramaru = &Toramaru{Help: false, Version: false}
	args = append(args, EOA)
	for i := 1; i < len(args)-1; i += 2 {
		opt := args[i]
		arg := args[i+1]

		if arg == EOA {
			return nil, errors.New("invalid options")
		}

		switch opt {
		case "-p", "--port":
			toramaru.Port, err = strconv.Atoi(arg)
			if err != nil {
				return nil, errors.New("port can not convert to int")
			}

		case "-r", "--route-proxy":
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

func isHelp(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}

func isVersion(args []string) bool {
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			return true
		}
	}
	return false
}
