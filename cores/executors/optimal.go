package executors

import (
	"goflow/commons"
)

// TODO: reply in proper format using reply packet
type ExecutorReplyPacket struct {
}

func ExecuteExecutable(executable commons.Executable, broadcastChannel chan struct{},
	dependantOnBroadcastChannels []chan struct{}, replyChannel chan interface{},
	executorWorkerLimitChannel chan struct{}) {
	defer close(broadcastChannel)

	var goFlowLogger *commons.Logger
	goFlowLogger = commons.GetLogger()

	if goFlowLogger != nil {
		goFlowLogger.LogDebug(
			nil,
			"received executable",
			goFlowLogger.String("executable name", executable.MyName()),
			goFlowLogger.Int("executable id", executable.MyID()),
			goFlowLogger.Int("number of notifiers", len(dependantOnBroadcastChannels)),
		)
	}

	// TODO: will have to send with executable id to mark status 'executing'
	// replyChannel <- commons.Executing
	for _, dependantOnBroadcastChannel := range dependantOnBroadcastChannels {
		select {
		case <-dependantOnBroadcastChannel:
		}
	}

	if goFlowLogger != nil {
		goFlowLogger.LogDebug(
			nil,
			"dependant-on executables completed, ready for execution",
			goFlowLogger.String("executable name", executable.MyName()),
			goFlowLogger.Int("executable id", executable.MyID()),
		)
	}

	executable.MarkExecuting()
	err := executable.Execute()
	if err != nil {
		replyChannel <- err
	}
	executable.MarkExecuted()

	dependantExecutables := executable.DependantExecutables()
	for _, dependantExecutable := range dependantExecutables {
		replyChannel <- dependantExecutable
	}

	replyChannel <- commons.Executed
	<-executorWorkerLimitChannel

	return
}
