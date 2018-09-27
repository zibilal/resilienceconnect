// This is rather a simple http.Client wrapper that comply with ConnectionFunc type function
package resilienceconnect

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

// HttpClient wrap the functionality to connect to http server
// This wraps http.Client object
type HttpClient struct {
	httpClient *http.Client
}

// Create new HttpClient object
// timout are in seconds
func NewHttpClient(timeout int) *HttpClient {
	h := new(HttpClient)
	h.httpClient = &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}
	return h
}

// ConnectWith connect to an external resource specified in url argument
// url is the external resource path
// request specify the request body and header specific to the external resource contract api
// ConnectWith assumes the external resource is an restful service
func (h *HttpClient) ConnectWith(request Requestor, output interface{}) error {
	jreq := NewJsonRequestWrapper()
	request.Request(jreq)
	response, err := h.httpClient.Do(jreq.HttpRequest)
	if err != nil {
		return err
	}
	if response == nil {
		return errors.New("unable to get http response")
	}
	if response.Body == nil {
		return errors.New("return an empty response body")
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

type HttpError struct {
	value string
}

func NewHttpError() *HttpError {
	return new(HttpError)
}

func (h *HttpError) Set(str string) {
	h.value = str
}

func (h *HttpError) Error() string {
	return h.value
}
