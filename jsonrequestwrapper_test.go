package resilienceconnect

import (
	"net/http"
	"testing"
	"reflect"
	"encoding/json"
	"github.com/zibilal/sqltyping"
)

const (
	success = "\u2713"
	failed  = "\u2717"
)

func TestNewJsonRequestWrapper(t *testing.T) {
	t.Log("Testing NewJsonRequestWrapper, struct input")
	{

		input := struct {
			OrderId string `json:"order_id"`
			Status  string `json:"status"`
		}{}

		jreq := NewJsonRequestWrapper()

		if jreq == nil {
			t.Fatalf("%s expected jreq object not nil", failed)
		}

		jreq.BodyRequest(http.MethodGet, "http://example.com", input)

		if !jreq.Valid() {
			t.Fatalf("%s expected valid request", failed)
		}
	}

	t.Log("Testing NewJsonRequestWrapper, map input")
	{

		input := map[string]interface{}{
			"order_id": "12345",
			"status": "accepted",
		}

		jreq := NewJsonRequestWrapper().BodyRequest(http.MethodGet,"http://example.com", input)
		if !jreq.Valid() {
			t.Logf("%s expected request is not valid", failed)
		}

	}
}

func TestJsonRequestWrapper_AddHeader(t *testing.T) {
	t.Log("Testing JsonRequestWrapper AddHeader method")
	{
		jreq := NewJsonRequestWrapper()

		if jreq == nil {
			t.Fatalf("%s expected jreq object not nil", failed)
		}

		input := struct {
			OrderId string `json:"order_id"`
			Status  string `json:"status"`
		}{"1123443", "applied"}

		jreq.BodyRequest(http.MethodGet,"http://example.com", input).AddHeader("X-Api-Key", "11234123413414").AddHeader("signature", "basasdfasdfffgg")

		expected := http.Header(map[string][]string{
			"X-Api-Key":{"11234123413414"},
			"Signature":{"basasdfasdfffgg"},
		})

		if  reflect.DeepEqual(jreq.HttpRequest.Header, expected) {
			t.Logf("%s expected header", success)
		} else {
			t.Log(expected, jreq.HttpRequest.Header)
			t.Fatalf("%s expected header equal %v", failed, jreq.HttpRequest.Header)
		}
	}
}

func TestJsonRequestWrapper_Request(t *testing.T) {
	t.Log("Testing JsonRequestWrapper Request method")
	{
		jreq := NewJsonRequestWrapper()

		input := struct {
			OrderId string `json:"order_id"`
			Status  string `json:"status"`
		}{"1123443", "applied"}
		jreq.BodyRequest(http.MethodGet,"http://example.com", input).AddHeader("X-Api-Key", "11234123413414").AddHeader("signature", "basasdfasdfffgg")

		zeInput := struct {
			OrderId string `json:"order_id"`
			Status  string `json:"status"`
		}{}
		_ = json.NewDecoder(jreq.HttpRequest.Body).Decode(&zeInput)

		if sqltyping.IsEmpty(zeInput) {
			t.Fatalf("%s expected request body not empty", failed)
		}

		if reflect.DeepEqual(input, zeInput) {
			t.Logf("%s expected input and zeInput are equal", success)
		} else {
			t.Fatalf("%s expected input and zeInput are equal", failed)
		}
	}

	t.Log("Testing JsonRequestWrapper Request method and testing AsCurlCommand function")
	{
		jreq := NewJsonRequestWrapper()

		input := struct {
			OrderId string `json:"order_id"`
			Status  string `json:"status"`
		}{"1123443", "applied"}
		jreq.BodyRequest(http.MethodGet,"http://example.com", input).AddHeader("X-Api-Key", "11234123413414").AddHeader("signature", "basasdfasdfffgg")

		zeInput := struct {
			OrderId string `json:"order_id"`
			Status  string `json:"status"`
		}{}
		_ = json.NewDecoder(jreq.HttpRequest.Body).Decode(&zeInput)

		if sqltyping.IsEmpty(zeInput) {
			t.Fatalf("%s expected request body not empty", failed)
		}

		if reflect.DeepEqual(input, zeInput) {
			t.Logf("%s expected input and zeInput are equal", success)
		} else {
			t.Fatalf("%s expected input and zeInput are equal", failed)
		}

		curlCommand, err := jreq.AsCurlCommand()
		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		}

		if curlCommand == "" {
			t.Fatalf("%s expected curlCommand not empty", failed)
		}

		t.Logf("%s curl command: %s", success, curlCommand)
	}
}
