package ap

import (
	"html/template"
	"net/http"
)

//RejectCode string representation
const RejectCode string = "REJECT"

//AcceptCode string representation
const AcceptCode string = "ACCEPT"

//OKCode string representation
const OKCode string = "OK"

const responseTemplate string = `"CODE" "{{.ResponseCode}}"
"RA" "{{.ResponseAuthorization}}"
{{if eq .ResponseCode "ACCEPT" -}}
"SECONDS" "{{.Seconds}}"
"DOWNLOAD" "{{.Download}}"
"UPLOAD" "{{.Upload}}"
{{- else if eq .ResponseCode "REJECT" -}}
"BLOCKED_MSG" "{{.BlockedMessage}}"
{{end}}`

//Response is the generated object being sent back to the Access Point
type Response struct {
	ResponseCode          string
	Request               *Request
	ResponseAuthorization string
	Seconds               int32
	Download              int32
	Upload                int32
	BlockedMessage        string
	Secret                string
}

//Execute the APResposne object, and write the response to the io writer
func (r *Response) Execute(w http.ResponseWriter) error {
	var err error

	//Pull in our template
	t := template.Must(template.New("response").Parse(responseTemplate))

	if err != nil {
		return err
	}

	//Execute the response template, and write to the response
	err = t.Execute(w, *r)
	if err != nil {
		return err
	}

	//return nil; no news is good news.
	return nil
}
