package graphs

// func (node *Node) SyncExecute() (wasExecutionSuccessful bool) {
// 	// defer newrelic.StartSegment(newRelicTransaction, node.MyName()).End()

// 	// egm := enigma.GetEnigma()
// 	// enigmaLogger := egm.GetLogger()
// 	// enigmaErrorCollection := egm.GetErrorCollection()

// 	err := node.Prepare()
// 	if err != nil {
// 		// enigmaLogger.LogError(
// 		// 	node.EnigmaContext,
// 		// 	enigmaErrorCollection.GetEnigmaErrorCategory(enigma_constants.GRAPH_NODE_PREPERATION),
// 		// 	enigmaLogger.String("reason", "unable to prepare node for sync execute"),
// 		// 	enigmaLogger.String("node name", node.MyName()),
// 		// 	enigmaLogger.Error(err),
// 		// )
// 		return false
// 	}

// 	container := node.GetContainer()
// 	if container != nil {
// 		// actorConstructionSegment := newrelic.StartSegment(newRelicTransaction, node.MyName()+" - Actor Construction")
// 		err = container.constructActor()
// 		// actorConstructionSegment.End()
// 		if err != nil {
// 			// not logging error since actor would have logged its error
// 			// just inform executor that execution had some error
// 			return false
// 		}

// 		// enigmaLogger.LogDebug(
// 		// 	node.EnigmaContext,
// 		// 	"actor successfully constructed",
// 		// 	enigmaLogger.String("node name", node.MyName()),
// 		// )

// 		// actorActingSegment := newrelic.StartSegment(newRelicTransaction, node.MyName()+" - Actor Acting")
// 		err := container.askActorToAct()
// 		// actorActingSegment.End()
// 		if err != nil {
// 			// not logging error since actor would have logged its error
// 			// just inform executor that execution had some error
// 			return false
// 		}

// 		// enigmaLogger.LogDebug(
// 		// 	node.EnigmaContext,
// 		// 	"actor successfully acted",
// 		// 	enigmaLogger.String("node name", node.MyName()),
// 		// )
// 	} else {
// 		// not logging error since actor would have logged its error
// 		// just inform executor that execution had some error
// 		return false
// 	}

// 	return true
// }
