// This is rather a simple http.Client wrapper that comply with ConnectionFunc type function
package resilienceconnect

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	acceptedContentTypes = "application/json | application/vnd.api+json"
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
// ConnectWith assumes the external resource is a restful service
func (h *JsonapiClient) ConnectWith(request Requestor, output interface{}) (Responder, error) {
	jreq := NewJsonRequestWrapper()
	_ = request.Request(jreq)
	contentType := jreq.HttpRequest.Header.Get("Content-Type")
	if contentType == "" {
		return nil, fmt.Errorf("header \"Content-Type\" is required, accepts(%s)", acceptedContentTypes)
	}
	if !h.validContentType(contentType) {
		return nil, fmt.Errorf("invalid content type (%s) accepts(%s)", contentType, acceptedContentTypes)
	}

	response, err := h.httpClient.Do(jreq.HttpRequest)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	err = json.NewDecoder(response.Body).Decode(output)
	if err != nil {
		b, _ := ioutil.ReadAll(response.Body)
		if len(b) > 0 {
			str := string(b)
			herr := NewHttpError()
			herr.Set(str)

			return nil, herr
		}
		return nil, err
	}

	jsonResponder := NewJsonApiClientResponder()
	jsonResponder.setHttpStatus(response.StatusCode)
	jsonResponder.setHttpMessage(response.Status)

	return jsonResponder, nil
}

func (h *JsonapiClient) validContentType(contentType string) bool {
	split := strings.Split(acceptedContentTypes, " | ")
	for _, s := range split {
		if contentType == s {
			return true
		}
	}
	return false
}

type JsonApiClientResponder struct {
	httpStatus  int
	httpMessage string
}

func NewJsonApiClientResponder() *JsonApiClientResponder {
	return new(JsonApiClientResponder)
}

func (c *JsonApiClientResponder) setHttpStatus(status int) {
	c.httpStatus = status
}

func (c *JsonApiClientResponder) StatusCode() int {
	return c.httpStatus
}

func (c *JsonApiClientResponder) setHttpMessage(message string) {
	c.httpMessage = message
}

func (c *JsonApiClientResponder) StatusMessage() string {
	return c.httpMessage
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
