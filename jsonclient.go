// This is rather a simple http.Client wrapper that comply with ConnectionFunc type function
package resilienceconnect

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// JsonapiClient wrap the functionality to connect to http server
// This wraps http.Client object
type JsonapiClient struct {
	httpClient *http.Client
}

// Create new JsonapiClient object
// timout are in seconds
func NewJsonapiClient(timeout int) *JsonapiClient {
	h := new(JsonapiClient)
	h.httpClient = &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}
	return h
}

// ConnectWith connect to an external resource specified in url argument
// url is the external resource path
// request specify the request body and header specific to the external resource contract api
// ConnectWith assumes the external resource is an restful service
func (h *JsonapiClient) ConnectWith(request Requestor, output interface{}) error {
	jreq := NewJsonRequestWrapper()
	request.Request(jreq)
	response, err := h.httpClient.Do(jreq.HttpRequest)
	if err != nil {
		return err
	}
	err = json.NewDecoder(response.Body).Decode(output)
	if err != nil {
		defer response.Body.Close()
		b, _ := ioutil.ReadAll(response.Body)
		if len(b) == 0 {
			str := string(b)
			herr := NewHttpError()
			herr.Set(str)

			return herr
		}
		return err
	}

	return nil
}

// This type can be use for wrapper html error, because external system is being load balanced with
// nginx or other load balancing system. Most of the time, if the external system can not handle
// the http error, the load balancing will taking care it, and mostly will send an html http error
type HttpError struct {
	value string
}

func NewHttpError() *HttpError {
	return new(HttpError)
}

// Set setup an error message specified in str parameter
func (h *HttpError) Set(str string) {
	h.value = str
}

// Error is the contract method to comply with error interface
func (h *HttpError) Error() string {
	return h.value
}
