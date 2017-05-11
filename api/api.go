package api

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/ryanhatfield/connectlabs-login/server"
	"github.com/ryanhatfield/connectlabs-login/sso"
)

type API struct {
	AuthorizeToken string
	Users          sso.UserStorage
}

func returnObject(res http.ResponseWriter, obj interface{}) {
	if j, err := json.MarshalIndent(obj, "", "  "); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	} else {
		res.Header().Set("Content-Type", "application/json")
		res.Write(j)
	}
}

func (api *API) AuthorizeAPI(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Checking API authorization header")
		auth := r.Header.Get("Authorization")
		sections := strings.Split(auth, " ")

		if auth == "" ||
			len(sections) < 1 ||
			strings.ToLower(sections[0]) != "token" ||
			strings.ToLower(sections[1]) != api.AuthorizeToken {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (api *API) GetUserByMAC(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(res, "unable to parse incoming form", http.StatusBadRequest)
		return
	}

	mac := req.Form.Get(":mac")
	if _, err := net.ParseMAC(mac); err != nil {
		http.Error(res, "unable to parse mac address", http.StatusBadRequest)
		return
	}

	user := &sso.User{}
	if err := api.Users.FindUserByDevice(mac, user); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	} else if user == nil {
		http.Error(res, "user not found", http.StatusNotFound)
		return
	}

	returnObject(res, user)
}

func (api *API) GetUserByID(res http.ResponseWriter, req *http.Request) {

	if err := req.ParseForm(); err != nil {
		http.Error(res, "unable to parse incoming form", http.StatusBadRequest)
		return
	}

	stringID := req.Form.Get(":id")

	tempID, err := strconv.ParseUint(stringID, 0, 32)
	if err != nil {
		http.Error(res, "unable to parse user id", http.StatusBadRequest)
		return
	}

	id := uint(tempID)
	user := &sso.User{}

	err = api.Users.FindUserByID(id, user)
	if err == server.ErrModelNotFound {
		http.Error(res, "unable to find user", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(res, "unable to find user", http.StatusInternalServerError)
		return
	}

	returnObject(res, user)
}
