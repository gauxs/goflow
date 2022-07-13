package graphs

import (
	"context"

	"errors"
	"goflow/commons"
	"reflect"

	newrelic "github.com/newrelic/go-agent"
)

func NewTopologyGraph(ctx context.Context, graphName string) *TopologyGraph {
	if graphName == "" {
		// graph name is mandatory
		return nil
	}

	tg := &TopologyGraph{
		ctx,
		graphName,
		commons.Prerun,
		0,
		make([]*Node, 0, commons.DefaultSliceCapacity),
		-1,
		-1,
	}

	tg.newInputHandlerNode(commons.InputHandlerNodeName)
	tg.newOutputHandlerNode(commons.OutputHandlerNodeName)

	return tg
}

type TopologyGraph struct {
	ctx                  context.Context
	name                 string
	graphState           commons.GraphState
	numberOfNodes        int
	allNodes             []*Node
	queryInputHandlerID  int
	queryOutputHandlerID int
}

func (topologyGraph *TopologyGraph) MyType() string {
	return reflect.TypeOf(topologyGraph).String()
}

func (topologyGraph *TopologyGraph) MyName() string {
	return topologyGraph.name
}

func (topologyGraph *TopologyGraph) Describe() map[string]interface{} {
	descriptionMap := make(map[string]interface{})
	descriptionMap["my type"] = topologyGraph.MyType()
	descriptionMap["graph Name"] = topologyGraph.name
	descriptionMap["number of nodes in graph"] = topologyGraph.numberOfNodes
	descriptionMap["input handler node id"] = topologyGraph.queryInputHandlerID
	descriptionMap["output handler node id"] = topologyGraph.queryOutputHandlerID
	descriptionMap["non-nil enigma context"] = topologyGraph.ctx != nil
	nodeDetails := make([]map[string]interface{}, 0, topologyGraph.numberOfNodes)
	for _, node := range topologyGraph.allNodes {
		nodeDetails = append(nodeDetails, node.Describe())
	}
	descriptionMap["node details"] = nodeDetails
	return descriptionMap
}

func (topologyGraph *TopologyGraph) NumberOfNodes() int {
	return topologyGraph.numberOfNodes
}

func (topologyGraph *TopologyGraph) setState(newGraphState commons.GraphState) {
	topologyGraph.graphState = newGraphState
}

func (topologyGraph *TopologyGraph) NewNode(nodeName string) *Node {
	if nodeName == "" {
		// avoid this
		nodeName = "Unnamed node"
	}

	node := newNode(topologyGraph.ctx, topologyGraph.numberOfNodes, nodeName)
	topologyGraph.numberOfNodes++
	topologyGraph.allNodes = append(topologyGraph.allNodes, node)

	return node
}

func (topologyGraph *TopologyGraph) newInputHandlerNode(nodeName string) *Node {
	if nodeName == "" {
		// avoid this
		nodeName = "Unnamed node"
	}

	topologyGraph.queryInputHandlerID = topologyGraph.numberOfNodes
	node := newNode(topologyGraph.ctx, topologyGraph.numberOfNodes, nodeName)
	topologyGraph.numberOfNodes++
	topologyGraph.allNodes = append(topologyGraph.allNodes, node)

	return node
}

func (topologyGraph *TopologyGraph) newOutputHandlerNode(nodeName string) *Node {
	if nodeName == "" {
		// avoid this
		nodeName = "Unnamed node"
	}

	topologyGraph.queryOutputHandlerID = topologyGraph.numberOfNodes
	node := newNode(topologyGraph.ctx, topologyGraph.numberOfNodes, nodeName)
	topologyGraph.numberOfNodes++
	topologyGraph.allNodes = append(topologyGraph.allNodes, node)

	return node
}

func (topologyGraph *TopologyGraph) getNode(index int) *Node {
	return topologyGraph.allNodes[index]
}

func (topologyGraph *TopologyGraph) GetContext() context.Context {
	return topologyGraph.ctx
}

func (topologyGraph *TopologyGraph) GetInputHandlerNode() *Node {
	return topologyGraph.getNode(topologyGraph.queryInputHandlerID)
}

func (topologyGraph *TopologyGraph) GetOutputHandlerNode() *Node {
	return topologyGraph.getNode(topologyGraph.queryOutputHandlerID)
}

// TODO: optmise this, no need to traverse all nodes everytime, use colors
func (topologyGraph *TopologyGraph) nextSetOfNodesInTopologicalOrder() []*Node {
	var nextNodes []*Node

	for index := 0; index < len(topologyGraph.allNodes); index++ {
		node := topologyGraph.allNodes[index]
		if node.nodeState == commons.ValidatedExecutable ||
			node.nodeState == commons.Unvalidated { //TODO: remove 2nd OR condition once validation is enabled
			addThisNode := true
			for _, edge := range node.edges.previousEdges {
				neighborNode := edge.getMyNeighborNode(node)
				if !(neighborNode.nodeState == commons.ExecutionFinished ||
					neighborNode.nodeState == commons.ValidatedUnexecutable ||
					neighborNode.nodeState == commons.ValidatedUnexecutableViaNormalisation) {
					addThisNode = false
					break
				}
			}

			if addThisNode {
				node.MarkPicked()
				nextNodes = append(nextNodes, node)
			}
		}
	}

	return nextNodes
}

func (topologyGraph *TopologyGraph) normalise(node *Node) {
	for _, edge := range node.edges.previousEdges {
		neighborNode := edge.getMyNeighborNode(node)
		if neighborNode != nil {
			normaliseBackwardNode(neighborNode)
		}
	}

	for _, edge := range node.edges.nextEdges {
		neighborNode := edge.getMyNeighborNode(node)
		if neighborNode != nil {
			normaliseForwardNode(neighborNode)
		}
	}

	return
}

func normaliseBackwardNode(node *Node) {
	totalContributingNextEdges := len(node.edges.nextEdges)
	if node.IsTransactionOutputProducer() {
		totalContributingNextEdges++
	}

	if totalContributingNextEdges > 1 {
		return
	}

	node.markValidatedButUnexecutableViaNormalisation()
	for _, edge := range node.edges.previousEdges {
		neighborNode := edge.getMyNeighborNode(node)
		normaliseBackwardNode(neighborNode)
	}

	return
}

func normaliseForwardNode(node *Node) {
	totalContributingPreviousEdges := len(node.edges.previousEdges)
	if node.IsTransactionInputConsumer() {
		totalContributingPreviousEdges++
	}

	if totalContributingPreviousEdges > 1 {
		return
	}

	node.markValidatedButUnexecutableViaNormalisation()
	for _, edge := range node.edges.nextEdges {
		neighborNode := edge.getMyNeighborNode(node)
		normaliseForwardNode(neighborNode)
	}

	return
}

func (topologyGraph *TopologyGraph) validate() error {
	var newRelicTransaction newrelic.Transaction
	if topologyGraph.ctx != nil {
		value := topologyGraph.ctx.Value(commons.NewRelicTransactionGoFlowContextKey())
		newRelicTransaction = value.(newrelic.Transaction)
		if newRelicTransaction != nil {
			defer newrelic.StartSegment(newRelicTransaction, topologyGraph.MyName()+" - Graph Validation").End()
		}
	}

	goFlowLogger := commons.GetLogger()
	for _, node := range topologyGraph.allNodes {
		if node.validator != nil {
			isValid, err := node.validator(node.ctx, topologyGraph.graphState)
			if err != nil {
				if goFlowLogger != nil {
					goFlowLogger.LogWarn(
						topologyGraph.ctx,
						"error while validating node, marking it executable",
						goFlowLogger.String("node name", node.MyName()),
					)
				}
				node.markValidatedAndExecutable()
			} else {
				if !isValid {
					node.markValidatedButUnexecutable()
				} else {
					node.markValidatedAndExecutable()
				}
			}
		} else {
			if goFlowLogger != nil {
				goFlowLogger.LogDebug(
					topologyGraph.ctx,
					"no validator for node, marking it executable",
					goFlowLogger.String("node name", node.MyName()),
				)
			}
			node.markValidatedAndExecutable()
		}
	}

	for _, node := range topologyGraph.allNodes {
		if node.nodeState == commons.ValidatedUnexecutable {
			topologyGraph.normalise(node)
		}
	}

	return nil
}

// TODO: validate all the node
func (topologyGraph *TopologyGraph) Construct(inputDataToGraph ...interface{}) error {
	var ok bool
	var newRelicTransaction newrelic.Transaction
	if topologyGraph.ctx != nil {
		value := topologyGraph.ctx.Value(commons.NewRelicTransactionGoFlowContextKey())
		newRelicTransaction, ok = value.(newrelic.Transaction)
		if ok && newRelicTransaction != nil {
			defer newrelic.StartSegment(newRelicTransaction, topologyGraph.MyName()+" - Graph Construction").End()
		}
	}

	inputHandlerNode := topologyGraph.GetInputHandlerNode()
	outputHandlerNode := topologyGraph.GetOutputHandlerNode()

	if inputHandlerNode == nil || outputHandlerNode == nil {
		return errors.New("inputHandler or outputHandler node nil")
	}

	for index, graphNode := range topologyGraph.allNodes {
		if index != topologyGraph.queryInputHandlerID && index != topologyGraph.queryOutputHandlerID {
			if len(graphNode.edges.previousEdges) == 0 {
				inputHandlerNode.AddDependant(graphNode)
			}

			if len(graphNode.edges.nextEdges) == 0 {
				graphNode.AddDependant(outputHandlerNode)
			}
		}
	}

	for _, data := range inputDataToGraph {
		err := inputHandlerNode.GetContainer().pushIntoInboundDataStore(data)
		if err != nil {
			return commons.WrapError(err, "unable to add data into input buffer of input handler node")
		}
	}

	// TODO: uncomment once validation done
	// err := topologyGraph.validate()
	// if err != nil {
	// 	enigma_utils.WrapError(err, "unable to validate graph")
	// }

	return nil
}

func (topologyGraph *TopologyGraph) GetOutput() []interface{} {
	return topologyGraph.GetOutputHandlerNode().GetContainer().getFromOutboundDataStore()
}

func (topologyGraph *TopologyGraph) SetStateStart() {
	topologyGraph.graphState = commons.Running
}

func (topologyGraph *TopologyGraph) SetStateStop() {
	topologyGraph.graphState = commons.Postrun
}
