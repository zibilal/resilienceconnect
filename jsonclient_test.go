package resilienceconnect

import (
	"encoding/json"
	"github.com/zibilal/sqltyping"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHttpClient_ConnectWith(t *testing.T) {
	t.Log("Testing ConnectWith")
	{
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := struct {
				Name  string `json:"full_name"`
				Email string `json:"email"`
			}{
				"Test Name", "test@example.com",
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			b, _ := json.Marshal(data)
			w.Write(b)
		}))

		defer ts.Close()

		jsonReq := NewJsonRequestWrapper()

		data := struct {
			Name  string `json:"full_name"`
			Email string `json:"email"`
		}{}
		jsonReq.BodyRequest(http.MethodGet, ts.URL, nil).AddHeader("Content-Type", "application/json")
		hclient := NewJsonapiClient(4)
		responder, err := hclient.ConnectWith(jsonReq, &data)

		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		}

		if sqltyping.IsEmpty(data) {
			t.Fatalf("%s expected data is not empty", failed)
		} else {
			t.Logf("%s value %+v", success, data)
		}

		t.Logf("%s Response code: %d", success, responder.StatusCode())
		t.Logf("%s Response message: %s", success, responder.StatusMessage())
	}

	t.Log("Testing ConnectWith, testing timeout")
	{
		jsonReq := NewJsonRequestWrapper()

		data := struct {
			Name  string `json:"full_name"`
			Email string `json:"email"`
		}{}
		jsonReq.BodyRequest(http.MethodGet, "http://www.foodboook.com/", nil).AddHeader("Content-Type", "application/json")
		hclient := NewJsonapiClient(3)
		now := time.Now()
		_, err := hclient.ConnectWith(jsonReq, &data)
		idata := time.Since(now)
		if err != nil {
			t.Logf("%s expected error not nil", success)
		}
		idata1 := int(idata)
		tmp := time.Duration(2)
		tmp2 := tmp * time.Second
		idata2 := int(tmp2)
		t.Log("Idata1", idata1, idata)
		t.Log("Idata2", idata2, 2*time.Second)

		if idata1 > idata2 {
			t.Logf("%s expected timeout not far from 3 seconds, got %d seconds", success, idata)
		}

	}

	t.Log("Testing ConnectWith, with nil http request")
	{
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		}))
		defer ts.Close()

		jsonReq := NewJsonRequestWrapper()
		hclient := NewJsonapiClient(2)
		data := struct {
			Name  string `json:"full_name"`
			Email string `json:"email"`
		}{}
		jsonReq.BodyRequest(http.MethodGet, ts.URL, nil).AddHeader("Content-Type", "application/json")
		_, err := hclient.ConnectWith(jsonReq, &data)
		if err != nil {
			t.Logf("%s expected error not nil, err: %s", success, err.Error())
		}
	}

	t.Log("Testing ConnectWith with empty server")
	{
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := struct {
				Name  string `json:"full_name"`
				Email string `json:"email"`
			}{"Example Name", "email@example.com"}
			json.NewEncoder(w).Encode(data)
		}))

		defer ts.Close()

		jsonReq := NewJsonRequestWrapper()

		data := struct {
			Name  string `json:"full_name"`
			Email string `json:"email"`
		}{}
		jsonReq.BodyRequest(http.MethodGet, ts.URL, nil).AddHeader("Content-Type", "application/json")
		hclient := NewJsonapiClient(2)
		_, err := hclient.ConnectWith(jsonReq, &data)

		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		}

		if sqltyping.IsEmpty(data) {
			t.Fatalf("%s expected data is not empty", failed)
		} else {
			t.Logf("%s value %+v", success, data)
		}
	}
}
