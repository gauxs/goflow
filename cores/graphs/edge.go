package graphs

import (
	"context"
	"goflow/commons"
	"reflect"
)

func newEdge() *edge {
	return &edge{
		edgeNode:      nil,
		weight:        0.0,
		otherEdgeNode: nil,
	}
}

type edge struct {
	ctx           context.Context
	edgeNode      *Node
	weight        float64
	otherEdgeNode *Node
}

func (e *edge) MyType() string {
	return reflect.TypeOf(e).String()
}

func (e *edge) MyName() string {
	return ""
}

func (e *edge) Describe() map[string]interface{} {
	descriptionMap := make(map[string]interface{})
	descriptionMap["my type"] = e.MyType()
	descriptionMap["edge weight"] = e.weight
	descriptionMap["node name"] = e.edgeNode.MyName()
	descriptionMap["other node name"] = e.otherEdgeNode.MyName()
	return descriptionMap
}

func (e *edge) setWeight(weight float64) {
	e.weight = weight
}

func (e *edge) setEdgeNode(node *Node) {
	e.edgeNode = node
}

func (e *edge) setOtherEdgeNode(node *Node) {
	e.otherEdgeNode = node
}

func (e *edge) getMyNeighborNode(currentNode *Node) *Node {
	if e.edgeNode == currentNode {
		return e.otherEdgeNode
	} else {
		return e.edgeNode
	}
}

func newEdges() *edges {
	return &edges{
		previousEdges: make([]*edge, 0, commons.DefaultSliceCapacity),
		nextEdges:     make([]*edge, 0, commons.DefaultSliceCapacity),
	}
}

type edges struct {
	nextEdges     []*edge
	previousEdges []*edge
}

func (edgs *edges) MyType() string {
	return reflect.TypeOf(edgs).String()
}

func (edgs *edges) MyName() string {
	return ""
}

func (edgs *edges) Describe() map[string]interface{} {
	descriptionMap := make(map[string]interface{})
	descriptionMap["number of next edges"] = len(edgs.nextEdges)
	nextEdgesDescriptionList := make([]map[string]interface{}, 0, len(edgs.nextEdges))
	for _, edge := range edgs.nextEdges {
		nextEdgesDescriptionList = append(nextEdgesDescriptionList, edge.Describe())
	}
	descriptionMap["number of previous edges"] = len(edgs.previousEdges)
	previousEdgesDescriptionList := make([]map[string]interface{}, 0, len(edgs.previousEdges))
	for _, edge := range edgs.previousEdges {
		previousEdgesDescriptionList = append(previousEdgesDescriptionList, edge.Describe())
	}
	return descriptionMap
}

func (edgs *edges) addNextEdge(e *edge) {
	edgs.nextEdges = append(edgs.nextEdges, e)
}

func (edgs *edges) addPreviousEdge(e *edge) {
	edgs.previousEdges = append(edgs.previousEdges, e)
}
