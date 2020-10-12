package apirest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"reflect"
	"testing"
)

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
		start    string
		final    string
		client   http.Client
		buf      *bytes.Buffer
		want     string
		fileName string
	}{
		{name: "Request Sucess",
			url:      "http://localhost:5000/",
			client:   http.Client{},
			buf:      bytes.NewBuffer([]byte("GRU,BRC,10\nBRC,SCL,5\nGRU,CDG,75\nGRU,SCL,20\nGRU,ORL,56\nORL,CDG,5\nSCL,ORL,20\n")),
			want:     "File Uploaded successfully\n",
			fileName: "test.csv",
		},
	}

	for _, tt := range myStructTestResponse {
		t.Run(tt.name, func(t *testing.T) {
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile("file", tt.fileName)

			part.Write(tt.buf.Bytes())

			writer.Close()

			request, _ := http.NewRequest(http.MethodPost, tt.url, body)
			request.Header.Set("Content-Type", writer.FormDataContentType())

			response, _ := tt.client.Do(request)
			got, _ := ioutil.ReadAll(response.Body)

			AssertResponsebody(t, string(got), string(tt.want))
		})
	}
}
func TestRequestConsult(t *testing.T) {
	myStructTestResponse := []struct {
		name   string
		url    string
		start  string
		final  string
		client http.Client
		json   []byte
		want   []byte
	}{
		{name: "Request Sucess",
			url:    "http://localhost:5000/consult",
			client: http.Client{},
			json:   []byte(`{"routes":"ORL-CDG"}`),
			want:   []byte(`{"start":"ORL","sprintroutes":"CDG ","cost":5}` + "\n"),
		},
	}

	for _, tt := range myStructTestResponse {
		t.Run(tt.name, func(t *testing.T) {

			request, _ := http.NewRequest(http.MethodPost, tt.url, bytes.NewBuffer(tt.json))
			request.Header.Set("Content-Type", "application/json")

			response, _ := tt.client.Do(request)
			got, _ := ioutil.ReadAll(response.Body)

			AssertResponsebody(t, string(got), string(tt.want))
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
