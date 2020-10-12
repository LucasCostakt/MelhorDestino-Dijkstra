package web

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

type testClient struct {
	Response string
}

func TestCreateHttp(t *testing.T) {
	myStructTestResponse := []struct {
		name string
		want *httpServer
	}{
		{name: "Sucess", want: new(httpServer)},
	}

	for _, tt := range myStructTestResponse {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateHttp()

			AssertResponsebody(t, got, tt.want)
		})
	}
}
func TestNewRoutes(t *testing.T) {
	myStructTestResponse := []struct {
		name string
		want *httpServer
	}{
		{name: "Sucess", want: new(httpServer)},
	}

	for _, tt := range myStructTestResponse {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRoutes(tt.want)
			AssertResponsebody(t, got, tt.want)
		})
	}
}

func TestRequestSave(t *testing.T) {
	myStructTestResponse := []struct {
		name     string
		url      string
		fileName string
		client   http.Client
		want     string
		buf      *bytes.Buffer
	}{
		{name: "Request Sucess",
			url:      "http://localhost:5000/",
			fileName: "test.csv",
			client:   http.Client{},
			want:     "File Uploaded successfully\n",
			buf:      bytes.NewBuffer([]byte("GRU,BRC,10\nBRC,SCL,5\nGRU,CDG,75\nGRU,SCL,20\nGRU,ORL,56\nORL,CDG,5\nSCL,ORL,20\n")),
		},
	}

	for _, tt := range myStructTestResponse {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := NewfileUploadRequest(tt.url, tt.buf, tt.fileName)
			response, _ := tt.client.Do(request)
			got, _ := ioutil.ReadAll(response.Body)
			AssertResponsebody(t, string(got), tt.want)
		})
	}
}
func TestRequestConsult(t *testing.T) {
	myStructTestResponse := []struct {
		name   string
		url    string
		start  string
		final  string
		client httpClient
		want   string
	}{
		{name: "Request Sucess",
			url:    "http://localhost:5000/",
			start:  "ORL",
			final:  "CDG",
			client: &http.Client{},
			want:   "Total de custo $5, e as rotas são: ORL -> CDG ",
		},
	}

	for _, tt := range myStructTestResponse {
		t.Run(tt.name, func(t *testing.T) {
			got := requestConsult(tt.start, tt.final, tt.client)
			AssertResponsebody(t, string(got), tt.want)
		})
	}
}
func TestFormatRoute(t *testing.T) {
	myStructTestResponse := []struct {
		name  string
		start string
		final string
		want  string
	}{
		{name: "Request Sucess",
			start: "Orl",
			final: "cDg",
			want:  "ORL-CDG",
		},
	}

	for _, tt := range myStructTestResponse {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatRoute(tt.start, tt.final)
			AssertResponsebody(t, string(got), tt.want)
		})
	}
}

func AssertResponsebody(t *testing.T, got, expectedResponse interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, expectedResponse) {
		str1 := fmt.Sprintf("%v", got)
		str2 := fmt.Sprintf("%v", expectedResponse)
		t.Errorf("body is wrong, got %q expectedResponse %q\n", str1, str2)
	}
}
