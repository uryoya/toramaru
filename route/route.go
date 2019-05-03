package route

import (
	"errors"
	"regexp"
	"strings"
)

type Route struct {
	Location string
	Host     string
}

// TODO: ngnxみたいなルール作る
func Parse(s string) (*Route, error) {
	location, host, err := splitAsLocationHost(s)
	if err != nil {
		return nil, err
	}
	return &Route{Location: location, Host: host}, nil
}

func splitAsLocationHost(s string) (location string, host string, err error) {
	buf := strings.Split(s, ">")
	switch {
	case len(buf) < 2:
		err = errors.New("delimiter `>` must use once")
	case len(buf) > 2:
		err = errors.New("delimiter `>` must use only once")
	default:
		location = buf[0]
		host = buf[1]
	}
	return location, host, err
}

func (r *Route) Match(path string) bool {
	matched, _ := regexp.MatchString(r.Location+`.*`, path)
	return matched
}
