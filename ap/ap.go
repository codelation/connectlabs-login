package ap

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/pat"
)

//AP holds information about connecting with cloudtrax APs
type AP struct {
	Secret   string
	Sessions SessionStorage
}

//SessionStorage defines methods needed to privide session storage
type SessionStorage interface {
	FindSession(session *Session, sessionID string) error
	UpdateSession(session Session, req *Request) error
}

//AddRoutes adds the ap authorize route to the router
func (ap *AP) AddRoutes(p *pat.Router) {
	p.Get("/auth/auth.html", ap.handleAPRequest)
}

func (ap *AP) handleAPRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("handling ap request")
	err := r.ParseForm()
	if err != nil {
		log.Println("error parsing request form")
		return
	}

	request := &Request{}
	response := &Response{}
	session := &Session{}

	sessions := ap.Sessions

	request.ParseForm(r.Form)
	if err = sessions.FindSession(session, request.Session); err != nil {

	}

	//Get the new response authorization
	response.ResponseAuthorization, err = GenerateRA(response.ResponseCode,
		request.RequestAuthorization, ap.Secret)

	if err != nil {
		//nothing will work after this, should we do something here?
		log.Printf("error occured while generating the response authenticator:\n%s", err.Error())
	}

	switch request.RequestType {
	case RequestTypeStatus:
		//TODO: always reject for now, but eventually this could be used to re-up credentials
		response.ResponseCode = RejectCode
		response.BlockedMessage = "Your session has expired."
	case RequestTypeAccounting:
		response.ResponseCode = OKCode
	case RequestTypeLogin:
		response.ResponseCode = AcceptCode
		response.Seconds = 3600
		response.Download = 2000
		response.Upload = 800
	}

	sessions.UpdateSession(*session, request)

	err = response.Execute(&w)

	if err != nil {
		log.Printf("error while handling Accounting Request response: %s\n", err.Error())
	}

	j, _ := json.Marshal(
		struct {
			Session  string
			Mac      string
			Download uint
			Upload   uint
		}{
			Session:  session.Session,
			Mac:      session.Device,
			Download: session.Download,
			Upload:   session.Upload,
		})
	log.Printf("Session: %s\n", j)
}

//GenerateRA takes the response CODE, the (un-decoded) RA field, and the site secret,
//and generates the Response Authentication token.
//NOTE: I don't like this method, it will be updated/changed/mamed at some point.
func GenerateRA(code string, ra string, secret string) (string, error) {
	var buffer bytes.Buffer
	var err error
	hasher := md5.New()

	decodedRa, err := hex.DecodeString(ra)
	if err != nil {
		return "", fmt.Errorf(
			"An error has occured while decoding the hex string.\n%s", err.Error())
	}
	buffer.WriteString(code)
	buffer.WriteString(string(decodedRa))
	buffer.WriteString(secret)
	_, err = hasher.Write(buffer.Bytes())
	if err != nil {
		return "", fmt.Errorf(
			"An error has occured while writing to the md5 hasher.\n %s", err.Error())
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
