package api

import (
	"log"
	"net/http"
	"strings"
)

var AuthorizeToken string

func AuthorizeAPI(res http.ResponseWriter, req *http.Request) {
	log.Println("Checking API authorization header")
	auth := req.Header.Get("Authorization")
	sections := strings.Split(auth, " ")

	if auth == "" ||
		len(sections) < 1 ||
		strings.ToLower(sections[0]) != "token" ||
		strings.ToLower(sections[1]) != AuthorizeToken {
		http.Error(res, "not authorized", http.StatusUnauthorized)
		return
	}

	//only call if authorized above
	GetUserByMAC(res, req)
}

func GetUserByMAC(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "not found", http.StatusNotFound)
}
