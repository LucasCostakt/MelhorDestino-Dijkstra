package dijkstra

//Interface contendo as funções de Dijkstra onde o Graph vai implementada-lá
type Store interface {
	CreatBestDestiny() (int, string, string)
	GetNodeEdges(Node) []Edge
	NewCostTable(Node) map[Node]int
	Dijkstra(string, Node) (int, map[string]int)
	Destinations() string
	AddNode(Node)
	AddEdge(Node, Node, int)
	DestinationExist(string, string) bool
}

type Destination struct {
	Start string `json:"start"`
	Final string `json:"Final"`
	Cost  int    `json:"Cost"`
}

type Graph struct {
	Edges []Edge
	Nodes []Node
}

type Edge struct {
	Parent Node
	Child  Node
	Cost   int
}

type Node struct {
	Name string
}
