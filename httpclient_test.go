package resilienceconnect

import (
	"encoding/json"
	"github.com/zibilal/sqltyping"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHttpClient_ConnectTo(t *testing.T) {
	t.Log("Testing ConnectWith")
	{
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := struct {
				Name  string `json:"full_name"`
				Email string `json:"email"`
			}{
				"Test Name", "test@example.com",
			}
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			b, _ := json.Marshal(data)
			w.Write(b)
		}))

		defer ts.Close()

		jsonReq := NewJsonRequestWrapper()

		data := struct {
			Name  string `json:"full_name"`
			Email string `json:"email"`
		}{}
		jsonReq.BodyRequest(http.MethodGet, ts.URL, nil)
		hclient := NewHttpClient(4)
		err := hclient.ConnectWith(jsonReq, &data)

		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		}

		if sqltyping.IsEmpty(data) {
			t.Fatalf("%s expected data is not empty", failed)
		} else {
			t.Logf("%s value %+v", success, data)
		}
	}

	t.Log("Testing ConnectWith, testing timeout")
	{
		jsonReq := NewJsonRequestWrapper()

		data := struct {
			Name  string `json:"full_name"`
			Email string `json:"email"`
		}{}
		jsonReq.BodyRequest(http.MethodGet, "http://www.foodboook.com/", nil)
		hclient := NewHttpClient(3)
		now := time.Now()
		err := hclient.ConnectWith(jsonReq, &data)
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
}
