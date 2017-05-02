package main

import "net/http"

type Session interface {
	Get(res *http.ResponseWriter, req *http.Request) (Session, error)
}
