package server

//
// import (
// 	"context"
// 	"encoding/gob"
// 	"encoding/json"
// 	"log"
// 	"net/http"
//
// 	"github.com/gorilla/sessions"
// )
//
// type UserContext struct {
// 	UserID      uint
// 	Mac         string
// 	Node        string
// 	RedirectURL string
// 	SSID        string
// }
//
// type ContextKey string
//
// const ContextKeySession ContextKey = "connectlabs_login"
// const ContextKeyUser ContextKey = "usercontext"
// const ContextKeyNode ContextKey = "nodeaddress"
// const ContextKeyMac ContextKey = "macaddress"
//
// var store sessions.Store
//
// type M map[string]interface{}
//
// func init() {
//
// 	gob.Register(new(ContextKey))
// 	gob.Register(&UserContext{})
// 	gob.Register(&M{})
// 	store = sessions.NewCookieStore([]byte("cookie-secret"))
// }
//
// func newContextWithUserContext(ctx context.Context, req *http.Request, res http.ResponseWriter) context.Context {
// 	//get saved context from cookie storage
// 	cs, err := store.Get(req, string(ContextKeySession))
//
// 	if err != nil {
// 		log.Printf("error getting cookie storage: %s\n", err.Error())
// 		return ctx
// 	}
//
// 	j, _ := json.MarshalIndent(cs.Values, "", "  ")
// 	log.Println(j)
//
// 	uc := &UserContext{}
// 	//get user context object
// 	jsonuc, ok := cs.Values[ContextKeyUser].(string)
// 	if !ok {
// 		log.Println("UserContext was not deserialized propperly")
// 		uc = &UserContext{}
// 	}
//
// 	log.Printf("json to unmarshal: %s\n", jsonuc)
//
// 	if err = json.Unmarshal([]byte(jsonuc), uc); err != nil {
// 		log.Printf("error when unmarshalling user context: %s\n", err.Error())
// 	}
//
// 	//helper function for getting url parameter and casting as string
// 	reqGet := func(key, def string, r *http.Request) string {
// 		v := r.Form.Get(key)
// 		if key == "" {
// 			return def
// 		}
// 		return v
// 	}
//
// 	req.ParseForm()
// 	uc.Node = reqGet("called", uc.Node, req)
// 	uc.Mac = reqGet("mac", uc.Mac, req)
// 	uc.SSID = reqGet("ssid", uc.SSID, req)
//
// 	cs, err = store.New(req, string(ContextKeySession))
// 	if err != nil {
// 		log.Printf("error getting new session: %s\n", err.Error())
// 	}
// 	j, _ = json.Marshal(uc)
//
// 	cs.Values[ContextKeyUser] = j
//
// 	j, _ = json.MarshalIndent(cs.Values, "", "  ")
// 	log.Println(j)
//
// 	err = cs.Save(req, res)
// 	if err != nil {
// 		log.Printf("error saving cookie: %s\n", err.Error())
// 	}
//
// 	return context.WithValue(ctx, ContextKeyUser, uc)
// }
//
// func RequestUserFromContext(ctx context.Context) *UserContext {
// 	return ctx.Value(ContextKeyUser).(*UserContext)
// }
//
// func UserContextMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
// 		ctx := newContextWithUserContext(req.Context(), req, res)
// 		next.ServeHTTP(res, req.WithContext(ctx))
// 	})
// }
