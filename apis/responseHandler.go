package apis

import (
	handlers "github.com/danmaina/HttpResponse"
	"net/http"
)

// returnResponse is a package encapsulated function that formats responses in a standard Json format.
func returnResponse(status int, err error, body interface{}, res http.ResponseWriter) {
	_ = handlers.Response{
		Status: status,
		Error:  err,
		Body:   body,
	}.ReturnResponse(res)
}
