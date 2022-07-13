package controllers

import (
	"fmt"
	"goflow/commons"
	"goflow/cores/executors"
	"goflow/cores/graphs"
)

type ControllerJobPacket struct {
	ReplyChannel chan error
	Graph        *graphs.TopologyGraph
}

func Control(job *ControllerJobPacket, controllerWorkerLimitChannel chan struct{}, executorWorkerLimitChannel chan struct{}) {
	job.Graph.SetStateStart()
	defer job.Graph.SetStateStop()

	var executableID int
	var totalExecutablesFinished int
	var totalExecutablesToBePicked int
	var totalNumberOfExecutablesRanUnsucessfully int

	var executable commons.Executable
	var executorWorkerReply interface{}
	var executorWorkerStatus commons.GoRoutineStatus

	var goFlowLogger *commons.Logger
	goFlowLogger = commons.GetLogger()

	totalExecutablesFinished = 0
	totalNumberOfExecutablesRanUnsucessfully = 0
	totalExecutablesToBePicked = job.Graph.NumberOfNodes()

	// TODO: optimise capacity
	executorWorkersReplyChannel := make(chan interface{}, 10)
	executorStatusMapping := make(map[int]commons.GoRoutineStatus)
	executableBroadcastChannelMapping := make(map[int]chan struct{})

	executorWorkersReplyChannel <- job.Graph.GetInputHandlerNode()
	for {
		select {
		case executorWorkerReply = <-executorWorkersReplyChannel:
			switch reply := executorWorkerReply.(type) {
			case commons.GoRoutineStatus:
				executorWorkerStatus = reply
				switch executorWorkerStatus {
				case commons.Executing:
					// TODO: add logging maybe?
				case commons.Executed:
					totalExecutablesFinished++
					if totalExecutablesFinished == totalExecutablesToBePicked {
						break
					}
				}

			case commons.Executable:
				executable = reply
				executableID = executable.MyID()
				if goFlowLogger != nil {
					goFlowLogger.LogDebug(
						nil,
						"received next executable",
						goFlowLogger.Int("executable id", executableID),
						goFlowLogger.String("executable name", executable.MyName()),
					)
				}

				_, ok := executorStatusMapping[executableID]
				if !ok {
					executable.MarkPicked()
					executorStatusMapping[executableID] = commons.Scheduled

					executorWorkerLimitChannel <- struct{}{}

					// insert a broadcast channel for the executable
					// if already present, then skip
					executableBroadcastChannel, ok := executableBroadcastChannelMapping[executableID]
					if !ok {
						executableBroadcastChannel = make(chan struct{})
						executableBroadcastChannelMapping[executableID] = executableBroadcastChannel
					}

					// find all the executable(s) on which this executable depends
					dependsOnExecutables := executable.DependsOnExecutables()
					dependsOnExecutablesBroadcastChannels := make([]chan struct{}, 0)
					for _, dependsOnExecutable := range dependsOnExecutables {
						_, ok := executableBroadcastChannelMapping[dependsOnExecutable.MyID()]
						if !ok {
							if goFlowLogger != nil {
								goFlowLogger.LogDebug(
									nil,
									"executable doesn't have a broadcast channel, making a new one",
									goFlowLogger.Int("executable id", executable.MyID()),
									goFlowLogger.String("executable name", executable.MyName()),
									goFlowLogger.Int("depends on executable with id", dependsOnExecutable.MyID()),
									goFlowLogger.String("depends on executable with name", dependsOnExecutable.MyName()),
								)
							}
							executableBroadcastChannelMapping[dependsOnExecutable.MyID()] = make(chan struct{})
						}
						dependsOnExecutablesBroadcastChannels = append(dependsOnExecutablesBroadcastChannels, executableBroadcastChannelMapping[dependsOnExecutable.MyID()])
					}

					// start the executor
					go executors.ExecuteExecutable(executable, executableBroadcastChannel,
						dependsOnExecutablesBroadcastChannels, executorWorkersReplyChannel, executorWorkerLimitChannel)
				} else {
					if goFlowLogger != nil {
						goFlowLogger.LogDebug(
							nil,
							"next executable already under execution",
							goFlowLogger.Int("executable id", executable.MyID()),
							goFlowLogger.String("executable name", executable.MyName()),
						)
					}
				}
			case error:
				totalNumberOfExecutablesRanUnsucessfully++
			}

		}
		if totalExecutablesFinished == totalExecutablesToBePicked {
			break
		}
	}

	if totalNumberOfExecutablesRanUnsucessfully > 0 {
		err := fmt.Errorf("some executables executed with error, total executables(%d) - executables which executed unsuccessfully(%d)", totalExecutablesToBePicked, totalNumberOfExecutablesRanUnsucessfully)
		job.ReplyChannel <- err
		<-controllerWorkerLimitChannel
		return
	} else {
		if goFlowLogger != nil {
			goFlowLogger.LogDebug(
				nil,
				"graph controlled without any error",
				goFlowLogger.Int("total executable", totalExecutablesToBePicked),
				goFlowLogger.Int("executable which ran successfully", totalExecutablesToBePicked-totalNumberOfExecutablesRanUnsucessfully),
			)
		}
	}

	job.ReplyChannel <- nil
	<-controllerWorkerLimitChannel

	return
}
