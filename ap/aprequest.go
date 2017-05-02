package ap

import "net/url"

//RequestTypeLogin string representation
const RequestTypeLogin string = "login"

//RequestTypeStatus string representation
const RequestTypeStatus string = "status"

//RequestTypeAccounting string representation
const RequestTypeAccounting string = "acct"

//Request has all variables expected in requests from an AP
type Request struct {
	RequestType          string
	RequestAuthorization string
	MacAddress           string
	Username             string
	Password             string
	NodeAddress          string //Mac address of the AP it's connected to
	IPV4Address          string
	Session              string
	Download             string
	Upload               string
	Seconds              string
}

//ParseForm reads information from the form post key/value pair array
func (r *Request) ParseForm(v url.Values) {
	r.RequestType = v.Get("type")
	r.RequestAuthorization = v.Get("ra")
	r.MacAddress = v.Get("mac")
	r.Username = v.Get("username")
	r.Password = v.Get("password")
	r.NodeAddress = v.Get("node")
	r.IPV4Address = v.Get("ipv4")
	r.Session = v.Get("session")
	r.Download = v.Get("download")
	r.Upload = v.Get("upload")
	r.Seconds = v.Get("seconds")
}
