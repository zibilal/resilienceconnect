// JSON Http request wrapper

package resilienceconnect

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/moul/http2curl"
)

// A JsonRequestWrapper represents an HTTP request which data
type JsonRequestWrapper struct {
	HttpRequest *http.Request
	// isSet will hold value of checking BodyRequest has already been executed successfully
	isSet bool
}

// NewHttpRequestWrapper assumes input is a valid struct data
// Method should have value the same as http verb words (GET|POST|DELETE)
func NewJsonRequestWrapper() *JsonRequestWrapper {
	return new(JsonRequestWrapper)
}

// This the wrapper of request.Header.Add functionality
// key is header key
// value is the value for the intended key
func (w *JsonRequestWrapper) AddHeader(key, value string) *JsonRequestWrapper {
	w.HttpRequest.Header.Add(key, value)
	return w
}

// Apply request body to supplied type
func (w *JsonRequestWrapper) BodyRequest(method, resource string, input interface{}) *JsonRequestWrapper {

	acceptedMethods := []string{
		http.MethodGet, http.MethodPost, http.MethodPut,
		http.MethodDelete, http.MethodPatch,
	}

	var foundMethod bool
	for i := 0; i < len(acceptedMethods); i++ {
		foundMethod = acceptedMethods[i] == method
		if foundMethod {
			break
		}
	}
	if !foundMethod {
		w.isSet = false
		return w
	}

	_, err := url.ParseRequestURI(resource)
	if err != nil {
		w.isSet = false
		return w
	}

	var body io.Reader
	if input != nil {
		b, err := json.Marshal(input)
		if err != nil {
			w.isSet = false
			return w
		}
		body = bytes.NewReader(b)
	}
	w.HttpRequest, err = http.NewRequest(method, resource, body)
	if err != nil {
		w.isSet = false
		return w
	}
	w.isSet = true

	return w
}

// Valid is used to check that this request had already wrap its request body
// This is for checking that BodyRequest function had already been executed successfully
func (w *JsonRequestWrapper) Valid() bool {
	return w.isSet
}

// Request will accept a pointer of JsonRequestWrapper
// dataRequest is of type *http.Request and will be filled with data that are currently
// hold by this object
func (w *JsonRequestWrapper) Request(dataRequest interface{}) error {
	jsonRequest, ok := dataRequest.(*JsonRequestWrapper)

	if !ok {
		return errors.New("please provide dataRequest with addressable value of type JsonRequestWrapper")
	}

	jsonRequest.HttpRequest = w.HttpRequest
	jsonRequest.isSet = w.isSet

	return nil
}

// Get request as string representing the runnable 'curl' command
// version of the request
func (w *JsonRequestWrapper) AsCurlCommand() (string, error) {
	if w.HttpRequest == nil {
		return "", errors.New("invalid state, http request still empty")
	}
	cmd, err := http2curl.GetCurlCommand(w.HttpRequest)
	if err != nil {
		return "", err
	}

	return cmd.String(), nil
}