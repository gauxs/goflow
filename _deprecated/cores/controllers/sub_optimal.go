package controllers

// import (
// 	graphs "INGO-goflow/cores/graphs"
// 	"fmt"
// )

// type SubOptimalController struct {
// 	chanCapacity int
// 	graphChannel chan *graphs.TopologyGraph
// }

// func NewSubOptimalController() (*SubOptimalController, error) {
// 	semiParallelExecutor, err := enigma_cores_executor_semi_parallel.NewSemiParallelExecutor(enigmaCTX)
// 	if err != nil {
// 		return nil, enigma_utils.WrapError(err, "unable to construct semi parallel executor")
// 	}

// 	return &SemiParallelController{
// 		enigmaCTX,
// 		semiParallelExecutor,
// 	}, nil
// }

// func (semiParallelController *SemiParallelController) MyName() string {
// 	return enigma_constants.SEMI_PARALLEL_CONTROLLER
// }

// func (semiParallelController *SemiParallelController) Control(collector enigma_commons.SemiParallelExecutablesCollector) error {
// 	collector.SetStateStart()

// 	var executionTurn int
// 	var totalNumberOfExecutables int
// 	var totalNumberOfNodeExecutedSuccessfully int
// 	for {
// 		executionTurn++
// 		collectionList := collector.NextEligibleSemiParallelExecutableCollection()
// 		if collectionList == nil || len(collectionList) == 0 {
// 			break
// 		}

// 		totalNumberOfExecutables += len(collectionList)
// 		semiParallelController.executor.AddExecutable(collectionList...)
// 		numberOfExecutablesExecutedSuccessfully := semiParallelController.executor.Run()
// 		totalNumberOfNodeExecutedSuccessfully += numberOfExecutablesExecutedSuccessfully
// 		semiParallelController.executor.ClearExecutable()

// 		enigmaLogger.LogDebug(
// 			semiParallelController.EnigmaContext,
// 			semiParallelController.MyName(),
// 			enigmaLogger.Int("execution turn number", executionTurn),
// 			enigmaLogger.Int("executables to be executed", len(collectionList)),
// 			enigmaLogger.Int("executables executed successfully", numberOfExecutablesExecutedSuccessfully),
// 		)
// 	}
// 	collector.SetStateStop()

// 	if totalNumberOfExecutables != totalNumberOfNodeExecutedSuccessfully {
// 		return fmt.Errorf("some executables didn't execute, total executables(%d) - successfully executed executables(%d)", totalNumberOfExecutables, totalNumberOfNodeExecutedSuccessfully)
// 	} else {
// 		enigmaLogger.LogDebug(
// 			semiParallelController.EnigmaContext,
// 			semiParallelController.MyName(),
// 			enigmaLogger.Int("total executable", totalNumberOfExecutables),
// 			enigmaLogger.Int("executable which ran successfully", totalNumberOfNodeExecutedSuccessfully),
// 		)
// 	}
// 	return nil
// }
