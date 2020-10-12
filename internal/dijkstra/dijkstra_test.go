package dijkstra

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCreatBestDestination(t *testing.T) {
	send := []Destination{
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
		name  string
		g     Graph
		start string
		final string
		want1 int
		want2 string
		want3 string
	}{
		{name: "Creat Best Destination Sucess",
			g:     Graph{},
			start: "BRC",
			final: "SCL",
			want1: 5,
			want2: "BRC",
			want3: "SCL ",
		},
	}

	for _, tt := range myStructTestResponse {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2, got3 := tt.g.CreatBestDestination(tt.start, tt.final, send)
			AssertResponsebody(t, got1, tt.want1)
			AssertResponsebody(t, got2, tt.want2)
			AssertResponsebody(t, got3, tt.want3)
		})
	}

}

func TestDestinationExist(t *testing.T) {
	g := Graph{}
	d := Node{Name: "BRC"}
	f := Node{Name: "SCL"}
	g.AddEdge(d, f, 5)

	myStructTestResponse := []struct {
		name  string
		start string
		final string
		want  bool
	}{
		{name: "Destination Exist Sucess",
			start: "BRC",
			final: "SCL",
			want:  true,
		},
	}

	for _, tt := range myStructTestResponse {
		t.Run(tt.name, func(t *testing.T) {
			got := g.DestinationExist(tt.start, tt.final)
			AssertResponsebody(t, got, tt.want)
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
