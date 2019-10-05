package main

import (
	"net/http"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var routes = []route{
	route{
		"Index",
		"GET",
		"/",
		index,
	},
	route{
		"MessageIndex",
		"GET",
		"/messages",
		messageIndex,
	},
	route{
		"MessageCreate",
		"POST",
		"/messages",
		MessageCreate,
	},
	route{
		"MessageShow",
		"GET",
		"/messages/{messageId}",
		messageShow,
	},
}
