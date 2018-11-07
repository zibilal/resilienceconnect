package resilienceconnect

import (
	"testing"
	"net/http"
	"fmt"
	"errors"
)

var countCall int

func TestHttpConnect_Connect_Retry(t *testing.T) {
	t.Log("Testing ResilienceConnect Connect")
	{
		connector := NewResilienceConnector()

		bodyRequest := struct {
			Name string `json:"full_name"`
			Email string `json:"email"`
		}{
			"Testing Name", "testing@example.com",
		}

		jsonRequest := NewJsonRequestWrapper()
		jsonRequest.BodyRequest(http.MethodGet, "http://example.com", &bodyRequest)

		jOptions := ConnectorOption{}
		jOptions.Put(IsRestrying, true)
		jOptions.Put(Retry, 2)
		jOptions.Put(ConnectorFunc, ConnectionFunc(ConnectResource1) )
		jOptions.Put(Wait, 1)

		response := struct {
			Name string `json:"full_name"`
			Email string `json:"email"`
			Status string `json:"status"`
		}{}

		err := connector.Connect(jsonRequest, jOptions, &response)
		if err != nil {
			t.Logf("%s expected error not nil", success)
		} else {
			t.Fatalf("%s expected error not nil, got nil error", failed)
		}
		if countCall == 2 {
			t.Logf("%s expected function is called 2 times", success)
		} else {
			t.Fatalf("%s expected function is called 2 times, actually called %d times", failed, countCall)
		}

		t.Log("Count call", countCall)
	}
}

func ConnectResource1(request Requestor, output interface{}) error {

	fmt.Println("Connect resources, called", countCall)
	countCall++

	return errors.New("testing error")
}

func TestHttpConnect_Connect_Backoff(t *testing.T) {
	t.Log("Testing ResilienceConnect Connect")
	{
		connector := NewResilienceConnector()

		bodyRequest := struct {
			Name string `json:"full_name"`
			Email string `json:"email"`
		}{
			"Testing Name", "testing@example.com",
		}

		jsonRequest := NewJsonRequestWrapper()
		jsonRequest.BodyRequest(http.MethodGet, "http://example.com", &bodyRequest)

		jOptions := ConnectorOption{}
		jOptions.Put(IsBackingOff, true)
		jOptions.Put(ConnectorFunc, ConnectionFunc(ConnectResource2) )
		jOptions.Put(Wait, 1)

		response := struct {
			Name string `json:"full_name"`
			Email string `json:"email"`
			Status string `json:"status"`
		}{}

		err := connector.Connect(jsonRequest, jOptions, &response)
		if err != nil {
			t.Fatalf("%s expected error not nil, got %s", failed, err.Error())
		}
		if countCall == 4{
			t.Logf("%s expected function is called 4 times", success)
		} else {
			t.Fatalf("%s expected function is called 4 times, actually called %d times", failed, countCall)
		}

		t.Log("Count call", countCall)
	}
}

func ConnectResource2( request Requestor, output interface{}) error {

	fmt.Println("Connect resources, called", countCall)
	countCall++
	if countCall < 4{
		return errors.New("testing error")
	} else {
		return nil
	}
}
