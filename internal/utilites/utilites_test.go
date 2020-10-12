package utilites

import (
	"fmt"
	"reflect"
	"testing"

	d "github.com/MelhorDestino/internal/dijkstra"
)

func TestSplitRoutes(t *testing.T) {
	myStructTestResponse := []struct {
		name string
		send string
		want []string
	}{
		{name: "Split string Sucess",
			send: "ORL-CDG",
			want: []string{"ORL", "CDG"},
		},
	}

	for _, tt := range myStructTestResponse {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitRoutes(tt.send)
			AssertResponsebody(t, got, tt.want)
		})
	}
}
func TestOpenFile(t *testing.T) {
	myStructTestResponse := []struct {
		name     string
		fileName string
		want     []string
	}{
		{name: "Open file Sucess",
			fileName: "test.csv",
			want:     []string{"GRU,BRC,10", "BRC,SCL,5", "GRU,CDG,75", "GRU,SCL,20", "GRU,ORL,56", "ORL,CDG,5", "SCL,ORL,20"},
		},
	}

	for _, tt := range myStructTestResponse {
		t.Run(tt.name, func(t *testing.T) {
			got := OpenFile(tt.fileName)
			AssertResponsebody(t, got, tt.want)
		})
	}
}

func TestConvertFileInDestinations(t *testing.T) {
	want := []d.Destination{
		{Start: "GRU",
			Final: "BRC",
			Cost:  10},
		{Start: "BRC",
			Final: "SCL",
			Cost:  5},
		{Start: "GRU",
			Final: "CDG",
			Cost:  75},
	}
	myStructTestResponse := []struct {
		name string
		send []string
	}{
		{name: "Convert File In Destinations Sucess",
			send: []string{"GRU,BRC,10", "BRC,SCL,5", "GRU,CDG,75"},
		},
	}

	for _, tt := range myStructTestResponse {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertFileInDestinations(tt.send)
			AssertResponsebody(t, got, want)
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
