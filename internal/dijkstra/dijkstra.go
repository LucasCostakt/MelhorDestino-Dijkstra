package dijkstra

import (
	"fmt"
	"reflect"
	"sort"
)

//valor de infinito dividido por 2
const Infinity = (int(^uint(0)>>1) / 2)

type myIntMap map[int]string

//adiciona uma aresta no gráfico
func (g *Graph) AddEdge(parent, child Node, cost int) {
	edge := Edge{
		Parent: parent,
		Child:  child,
		Cost:   cost,
	}

	g.Edges = append(g.Edges, edge)
	g.AddNode(parent)
	g.AddNode(child)
}

//Adiciona um nó à lista de nós do gráfico somente se o nó ainda não foi adicionado
func (g *Graph) AddNode(node Node) {
	var isPresent bool
	for _, n := range g.Nodes {
		if reflect.DeepEqual(n, node) {
			isPresent = true
		}
	}
	if !isPresent {
		g.Nodes = append(g.Nodes, node)
	}
}

//Algoritimo de Dijkstra
//Retorna o caminho mais curto de startNode para todos os outros nós
func (g *Graph) Dijkstra(destino string, startNode Node) (int, map[string]int) {
	var sumSmallestEdge int
	//mapa das rotas
	routes := make(map[string]int)

	costTable := g.NewCostTable(startNode)

	var visited []Node

	//Percorre todos os nós
	for len(visited) != len(g.Nodes) {

		//Retorna o nó não visitado mais próximo (custo mais baixo) da tabela de custos
		node := getClosestNonVisitedNode(costTable, visited)

		//Marca o nó como visitado
		visited = append(visited, node)
		nodeEdges := g.GetNodeEdges(node)

		for _, edge := range nodeEdges {
			distanceToNeighbor := costTable[node] + edge.Cost

			//Se a distância acima for menor que a distância atualmente na tabela de custos para aquele vizinho então
			//Atualiza a tabela de custos para aquele vizinho
			if distanceToNeighbor < costTable[edge.Child] {
				costTable[edge.Child] = distanceToNeighbor
				//se a distancia para o vizinho for maior que a rota atualiza a melhor rota
				if distanceToNeighbor > routes[edge.Child.Name] {
					routes[edge.Child.Name] = distanceToNeighbor
				}
			}
		}
	}

	//Organiza a tabela de custo
	for node, cost := range costTable {
		if destino == node.Name {
			sumSmallestEdge = cost

		}
	}

	return sumSmallestEdge, routes
}

//Retorna uma tabela de custo inicializada, para o algoritmo Dijkstra temos que atribuir o menor de custo ao
//startNode (nó inicial no nosso caso é o destino/aeroporto de saída)
//então começar a distribuir valores de infinito para os outro nós do gráfico
func (g *Graph) NewCostTable(startNode Node) map[Node]int {
	costTable := make(map[Node]int)
	costTable[startNode] = 0
	for _, node := range g.Nodes {
		if node != startNode {
			costTable[node] = Infinity
		}
	}
	return costTable
}

//Retorna o map com indice trocado pelo valor
func exchangeIndex(routes map[string]int) map[int]string {
	sortRoutes := make(map[int]string)
	for node, indice := range routes {
		sortRoutes[indice] = node
	}
	return sortRoutes
}

//Ordena o slice com base na ordem das escalas de destinos
func (m myIntMap) sort() (index []int) {
	for k := range m {
		index = append(index, k)
	}
	sort.Ints(index)
	return
}

//Retorna a string de escalas no formato "GRU -> ORL"
func (m myIntMap) sprintMap(destino string) string {
	var sprint string
	for {
		for _, k := range m.sort() {
			if destino == m[k] {
				sprint += fmt.Sprintf("%s ", m[k])
				return sprint
			} else {
				sprint += fmt.Sprintf("%s -> ", m[k])
			}
		}
	}
}

//Retorna todas as arestas que começam com o nó especificado nesse caso todas as arestas vizinhas do nó
func (g *Graph) GetNodeEdges(node Node) (edges []Edge) {
	for _, edge := range g.Edges {
		if edge.Parent == node {
			edges = append(edges, edge)
		}
	}
	return edges
}

//Retorna o nó mais próximo com o custo mais baixo da costTable se o nó ainda não foi visitado
func getClosestNonVisitedNode(costTable map[Node]int, visited []Node) Node {
	type CostTableToSort struct {
		Node Node
		Cost int
	}
	var sorted []CostTableToSort

	//Verifica se o nó já foi visitado
	for node, cost := range costTable {
		var isVisited bool
		for _, visitedNode := range visited {
			if node == visitedNode {
				isVisited = true
			}
		}
		//Se não foi visitado adiciona ao slice de visitados
		if !isVisited {
			sorted = append(sorted, CostTableToSort{node, cost})
		}
	}

	// Precisamos do Nó com o menor custo da tabela de custos
	// Portanto, é importante classificá-lo
	// Enão classifica o mapa
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Cost < sorted[j].Cost
	})

	return sorted[0].Node
}

//vai fazer as chamadas e  conversões nessárias para que retorne o destino inicial os destinos na escala e o valor da passagem
func (g *Graph) CreatBestDestination(start string, final string, destination []Destination) (int, string, string) {
	//utilizado para c
	m := myIntMap{}

	//cria os nós nessários representados pelos destinos e faz uma aresta com o valor do custo sendo o valor de um destino a outro
	var startNode Node
	for _, str := range destination {
		d := Node{Name: str.Start}
		f := Node{Name: str.Final}
		if start == d.Name {
			startNode = d
		}
		g.AddEdge(d, f, str.Cost)
	}

	if exist := g.DestinationExist(start, final); exist != true {
		return 0, "destino não encontrado", ""
	}

	cost, routes := g.Dijkstra(final, startNode)
	sortRoutes := exchangeIndex(routes)
	m = sortRoutes
	sprintRoutes := m.sprintMap(final)
	return cost, startNode.Name, sprintRoutes
}

//verifica se o destinos partida e chegada existem
func (g *Graph) DestinationExist(start string, final string) bool {
	var isPresentStart bool
	var isPresentFinal bool
	for _, n := range g.Nodes {
		if reflect.DeepEqual(n.Name, start) {
			isPresentStart = true
		}
		if reflect.DeepEqual(n.Name, final) {
			isPresentFinal = true
		}
	}
	if !isPresentStart || !isPresentFinal {
		return false
	}
	return true
}
