package graphs

import (
	"context"
	"errors"
	"goflow/commons"
	"reflect"

	newrelic "github.com/newrelic/go-agent"
)

func newNode(ctx context.Context, nodeID int, name string) *Node {
	return &Node{
		ctx,
		false,
		false,
		nodeID,
		name,
		commons.Unvalidated,
		newEdges(),
		newContainer(ctx),
		nil,
	}
}

type Node struct {
	ctx                            context.Context
	consumesTransactionInput       bool
	contributesToTransactionOutput bool
	id                             int
	name                           string
	nodeState                      commons.NodeState
	edges                          *edges
	container                      *container
	validator                      commons.NodeValidator
}

func (node *Node) MyID() int {
	return node.id
}

func (node *Node) MyType() string {
	return reflect.TypeOf(node).String()
}

func (node *Node) MyName() string {
	return node.name
}

func (node *Node) Describe() map[string]interface{} {
	descriptionMap := make(map[string]interface{})
	descriptionMap["my type"] = node.MyType()
	descriptionMap["node name"] = node.name
	descriptionMap["node id"] = node.id
	descriptionMap["node state"] = node.nodeState.String()
	descriptionMap["non-nil enigma context"] = node.ctx != nil
	descriptionMap["consumes transaction input?"] = node.IsTransactionInputConsumer()
	descriptionMap["contributes to transaction output?"] = node.IsTransactionOutputProducer()
	descriptionMap["edge details"] = node.edges.Describe()
	descriptionMap["container details"] = node.container.Describe()
	return descriptionMap
}

func (node *Node) GetContainer() *container {
	return node.container
}

func (node *Node) AttachValidator(validator commons.NodeValidator) {
	node.validator = validator
}

func (node *Node) AddDependant(dependantNodes ...*Node) {
	var edg *edge
	for _, dependantNode := range dependantNodes {
		edg = newEdge()
		edg.setWeight(0.0)
		edg.setEdgeNode(node)
		node.edges.addNextEdge(edg)

		edg.setOtherEdgeNode(dependantNode)
		dependantNode.edges.addPreviousEdge(edg)
	}

	return
}

func (node *Node) IsTransactionInputConsumer() bool {
	return node.consumesTransactionInput
}

func (node *Node) MarkConsumesTransactionInput() {
	node.consumesTransactionInput = true
}

func (node *Node) IsTransactionOutputProducer() bool {
	return node.contributesToTransactionOutput
}

func (node *Node) MarkProducesTransactionOutput() {
	node.contributesToTransactionOutput = true
}

func (node *Node) markValidatedButUnexecutable() {
	goFlowLogger := commons.GetLogger()

	if goFlowLogger != nil {
		goFlowLogger.LogDebug(
			node.ctx,
			"validated and marked unexecutable",
			goFlowLogger.String("node name", node.MyName()),
		)
	}

	node.nodeState = commons.ValidatedUnexecutable
}

func (node *Node) markValidatedButUnexecutableViaNormalisation() {
	goFlowLogger := commons.GetLogger()

	if goFlowLogger != nil {
		goFlowLogger.LogDebug(
			node.ctx,
			"validated and marked unexecutable via normalisation",
			goFlowLogger.String("node name", node.MyName()),
		)
	}

	node.nodeState = commons.ValidatedUnexecutableViaNormalisation
}

func (node *Node) markValidatedAndExecutable() {
	goFlowLogger := commons.GetLogger()

	if goFlowLogger != nil {
		goFlowLogger.LogDebug(
			node.ctx,
			"validated and marked executable",
			goFlowLogger.String("node name", node.MyName()),
		)
	}

	node.nodeState = commons.ValidatedExecutable
}

func (node *Node) markPickedForExecution() {
	node.nodeState = commons.PickedForExecution
}

func (node *Node) markExecutionStarted() {
	node.nodeState = commons.ExecutionStarted
}

func (node *Node) markExecutionFinished() {
	node.nodeState = commons.ExecutionFinished
}

func (node *Node) MarkPicked() {
	node.markPickedForExecution()
}

func (node *Node) IsPicked() bool {
	return node.nodeState == commons.PickedForExecution
}

func (node *Node) MarkExecuting() {
	node.markExecutionStarted()
}

func (node *Node) IsExecuting() bool {
	return node.nodeState == commons.ExecutionStarted
}

func (node *Node) MarkExecuted() {
	node.markExecutionFinished()
}

func (node *Node) HasExecuted() bool {
	return node.nodeState == commons.ExecutionFinished
}

func (node *Node) gatherForInboundDatastore() error {
	if node.container == nil {
		return errors.New("node has nil container")
	}

	goFlowLogger := commons.GetLogger()

	numOfPrevEdges := len(node.edges.previousEdges)
	if numOfPrevEdges > 0 {
		for index := 0; index < numOfPrevEdges; index++ {
			edge := node.edges.previousEdges[index]
			neighborNode := edge.getMyNeighborNode(node)
			if neighborNode != nil {
				if neighborNode.container != nil {
					dataFromOutboundBuffer := neighborNode.GetContainer().getFromOutboundDataStore()
					for _, data := range dataFromOutboundBuffer {
						err := node.GetContainer().pushIntoInboundDataStore(data)
						if err != nil {
							return commons.WrapError(err, "unable to add data in node's inbound data store")
						}
					}
				} else {
					if goFlowLogger != nil {
						goFlowLogger.LogWarn(
							node.ctx,
							"nil container for previous neighbor",
							goFlowLogger.String("node name", node.MyName()),
							goFlowLogger.String("neighbor node name", neighborNode.MyName()),
						)
					}
				}
			} else {
				if goFlowLogger != nil {
					goFlowLogger.LogWarn(
						node.ctx,
						"nil node for previous edge",
						goFlowLogger.String("node name", node.MyName()),
					)
				}
			}
		}
	} else {
		if goFlowLogger != nil {
			goFlowLogger.LogWarn(
				node.ctx,
				"no previous edges",
				goFlowLogger.String("node name", node.MyName()),
			)
		}
	}

	return nil
}

func (node *Node) Prepare() error {
	if node.ctx != nil {
		value := node.ctx.Value(commons.NewRelicTransactionGoFlowContextKey())
		newRelicTransaction, ok := value.(newrelic.Transaction)
		if ok && newRelicTransaction != nil {
			defer newrelic.StartSegment(newRelicTransaction, node.MyName()+" - Node Preperation").End()
		}
	}

	err := node.gatherForInboundDatastore()
	if err != nil {
		return commons.WrapError(err, "unable to gather data from previous nodes")
	}

	return nil
}

func (node *Node) Execute() error {
	var ok bool
	var newRelicTransaction newrelic.Transaction
	if node.ctx != nil {
		value := node.ctx.Value(commons.NewRelicTransactionGoFlowContextKey())
		newRelicTransaction, ok = value.(newrelic.Transaction)
		if ok && newRelicTransaction != nil {
			defer newrelic.StartSegment(newRelicTransaction, node.MyName()).End()
		}
	}

	goFlowLogger := commons.GetLogger()

	err := node.Prepare()
	if err != nil {
		if goFlowLogger != nil {
			goFlowLogger.LogError(
				node.ctx,
				err,
				goFlowLogger.String("reason", "unable to prepare node for sync execute"),
				goFlowLogger.String("node name", node.MyName()),
			)
		}
		return err
	}

	container := node.GetContainer()
	if container != nil {
		var actorConstructionSegment *newrelic.Segment
		if newRelicTransaction != nil {
			actorConstructionSegment = newrelic.StartSegment(newRelicTransaction, node.MyName()+" - Actor Construction")
		}
		err = container.constructActor()
		if actorConstructionSegment != nil {
			actorConstructionSegment.End()
		}

		if err != nil {
			// not logging error since actor would have logged its error
			// just inform executor that execution had some error
			return err
		}

		if goFlowLogger != nil {
			goFlowLogger.LogDebug(
				node.ctx,
				"actor successfully constructed",
				goFlowLogger.String("node name", node.MyName()),
			)
		}

		var actorActingSegment *newrelic.Segment
		if newRelicTransaction != nil {
			actorActingSegment = newrelic.StartSegment(newRelicTransaction, node.MyName()+" - Actor Acting")
		}
		err := container.askActorToAct()
		if actorActingSegment != nil {
			actorActingSegment.End()
		}

		if err != nil {
			// not logging error since actor would have logged its error
			// just inform executor that execution had some error
			return err
		}

		if goFlowLogger != nil {
			goFlowLogger.LogDebug(
				node.ctx,
				"actor successfully acted",
				goFlowLogger.String("node name", node.MyName()),
			)
		}
	} else {
		// not logging error since actor would have logged its error
		// just inform executor that execution had some error
		return errors.New("container nil for node" + node.MyName())
	}

	return nil
}

func (node *Node) DependantExecutables() []commons.Executable {
	dependantExecutables := make([]commons.Executable, 0)
	for _, edge := range node.edges.nextEdges {
		dependantExecutables = append(dependantExecutables, edge.getMyNeighborNode(node))
	}

	return dependantExecutables
}

func (node *Node) DependsOnExecutables() []commons.Executable {
	dependsOnExecutables := make([]commons.Executable, 0)
	for _, edge := range node.edges.previousEdges {
		dependsOnExecutables = append(dependsOnExecutables, edge.getMyNeighborNode(node))
	}
	return dependsOnExecutables
}
