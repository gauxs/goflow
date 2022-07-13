package controllers

// func NewOptimalController(graphChan chan *graphs.TopologyGraph, executableChan chan goflow.Executable) *OptimalController {
// 	return &OptimalController{
// 		graphChannel:      graphChan,
// 		executableChannel: executableChan,
// 	}
// }

// type OptimalController struct {
// 	graphChannel      chan *graphs.TopologyGraph
// 	executableChannel chan goflow.Executable
// }

// func (oc *OptimalController) MyType() string {
// 	return reflect.TypeOf(oc).String()
// }

// func (oc *OptimalController) MyName() string {
// 	return ""
// }

// func (oc *OptimalController) Run() {
// 	for {
// 		select {
// 		case graph := <-oc.graphChannel:
// 			// control graph
// 			oc.Control(graph)
// 		}
// 	}
// }

// func (oc *OptimalController) Control(collector goflow.FullParallelExecutablesCollector) error {
// 	var numberOfExecutablesRemaining, totalNumberOfExecutablesRanSuccessfully int
// 	var receivedExecutable goflow.FullParallelExecutable

// 	collector.SetStateStart()

// 	firstEligibleExecutable := collector.FirstEligibleFullParallelExecutable()
// 	if firstEligibleExecutable == nil {
// 		return errors.New("nil first full parallel executable")
// 	}

// 	numberOfExecutablesRemaining = 0
// 	totalNumberOfExecutablesRanSuccessfully = 0

// 	executbaleIsExecutionHandledMap := make(map[int64]bool)
// 	executableBroadcastChannelMap := make(map[int64]chan struct{})
// 	executableOutputChannel := make(chan goflow.FullParallelExecutable, 1)

// 	executableOutputChannel <- firstEligibleExecutable
// 	for {
// 		select {
// 		case receivedExecutable = <-executableOutputChannel:
// 			if receivedExecutable == nil {
// 				numberOfExecutablesRemaining--
// 				if numberOfExecutablesRemaining == 0 {
// 					break
// 				}
// 			} else {
// 				// enigmaLogger.LogDebug(
// 				// 	fullParallelController.EnigmaContext,
// 				// 	"received next executable",
// 				// 	enigmaLogger.String("controller name", fullParallelController.MyName()),
// 				// 	enigmaLogger.Int("executable id", int(receivedExecutable.MyID())),
// 				// 	enigmaLogger.String("executable name", receivedExecutable.MyName()),
// 				// )

// 				if !executbaleIsExecutionHandledMap[receivedExecutable.MyID()] {
// 					receivedExecutable.MarkPicked()
// 					fullParallelController.Start(receivedExecutable, executableBroadcastChannelMap, executableOutputChannel)
// 					executbaleIsExecutionHandledMap[receivedExecutable.MyID()] = true
// 					numberOfExecutablesRemaining++
// 					totalNumberOfExecutablesRanSuccessfully++
// 				} else {
// 					// enigmaLogger.LogDebug(
// 					// 	fullParallelController.EnigmaContext,
// 					// 	"next executable already under execution",
// 					// 	enigmaLogger.String("controller name", fullParallelController.MyName()),
// 					// 	enigmaLogger.Int("executable id", int(receivedExecutable.MyID())),
// 					// 	enigmaLogger.String("executable name", receivedExecutable.MyName()),
// 					// )
// 				}
// 			}
// 		}
// 		if numberOfExecutablesRemaining == 0 {
// 			break
// 		}
// 	}

// 	collector.SetStateStop()
// 	if len(executbaleIsExecutionHandledMap) != totalNumberOfExecutablesRanSuccessfully {
// 		return fmt.Errorf("some executables didn't execute, total executables(%d) - successfully executed executables(%d)", len(executbaleIsExecutionHandledMap), totalNumberOfExecutablesRanSuccessfully)
// 	} else {
// 		// enigmaLogger.LogDebug(
// 		// 	fullParallelController.EnigmaContext,
// 		// 	fullParallelController.MyName(),
// 		// 	enigmaLogger.Int("total executable", len(executbaleIsExecutionHandledMap)),
// 		// 	enigmaLogger.Int("executable which ran successfully", totalNumberOfExecutablesRanSuccessfully),
// 		// )
// 	}
// 	return nil
// }

// func (oc *OptimalController) Start(fullParallelExecutable goflow.FullParallelExecutable, executableBroadcastChannelMap map[int64]chan struct{}, executableOutputChannel chan enigma_commons.FullParallelExecutable) {
// 	chanList := make([]chan struct{}, 0)

// 	fullParallelExecutor := enigma_cores_executor_full_parallel.NewFullParallelExecutor(fullParallelController.EnigmaContext)
// 	fullParallelExecutor.AddExecutable(fullParallelExecutable)
// 	dependsOnExecutables := fullParallelExecutable.DependsOnExecutables()
// 	for _, dependsOnExecutable := range dependsOnExecutables {
// 		_, ok := executableBroadcastChannelMap[dependsOnExecutable.MyID()]
// 		if !ok {
// 			// enigmaLogger.LogDebug(
// 			// 	fullParallelController.EnigmaContext,
// 			// 	"executable doesn't have a broadcast channel, making a new one",
// 			// 	enigmaLogger.String("controller name", fullParallelController.MyName()),
// 			// 	enigmaLogger.Int("executable id", int(fullParallelExecutable.MyID())),
// 			// 	enigmaLogger.String("executable name", fullParallelExecutable.MyName()),
// 			// 	enigmaLogger.Int("depends on executable with id", int(dependsOnExecutable.MyID())),
// 			// 	enigmaLogger.String("depends on executable with name", dependsOnExecutable.MyName()),
// 			// )
// 			executableBroadcastChannelMap[dependsOnExecutable.MyID()] = make(chan struct{})
// 		}
// 		chanList = append(chanList, executableBroadcastChannelMap[dependsOnExecutable.MyID()])
// 	}
// 	fullParallelExecutor.SetInputChannels(chanList)

// 	if _, ok := executableBroadcastChannelMap[fullParallelExecutable.MyID()]; !ok {
// 		executableBroadcastChannelMap[fullParallelExecutable.MyID()] = make(chan struct{})
// 	}

// 	fullParallelExecutor.SetBroadcastChannel(executableBroadcastChannelMap[fullParallelExecutable.MyID()])
// 	fullParallelExecutor.SetOutputChannel(executableOutputChannel)
// 	fullParallelExecutor.Run()

// 	return
// }
