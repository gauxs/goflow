package executors

// import (
// 	goflow "INGO-goflow"
// 	"context"
// 	"reflect"
// )

// func NewOptimalExecutor() *OptimalExecutor {
// 	return &OptimalExecutor{
// 		nil,
// 		nil,
// 		nil,
// 		nil,
// 	}
// }

// type OptimalExecutor struct {
// 	broadcastChannel chan struct{}
// 	notifierChannels []chan struct{}
// 	executables      []goflow.Executable
// 	outputChannel    chan goflow.Executable
// }

// func (oe *OptimalExecutor) MyType() string {
// 	return reflect.TypeOf(oe).String()
// }

// func (oe *OptimalExecutor) MyName() string {
// 	return ""
// }

// func (oe *OptimalExecutor) AddExecutable(executables ...goflow.FullParallelExecutable) {
// 	oe.executables = executables
// }

// func (oe *OptimalExecutor) SetNotifierChannels(notifierChans []chan struct{}) {
// 	oe.notifierChannels = notifierChans
// }

// func (oe *OptimalExecutor) SetBroadcastChannel(broadcastChannel chan struct{}) {
// 	oe.broadcastChannel = broadcastChannel
// }

// func (oe *OptimalExecutor) SetOutputChannel(outputChannel chan goflow.FullParallelExecutable) {
// 	oe.outputChannel = outputChannel
// }

// func (oe *OptimalExecutor) Clear() {
// 	oe.executables = nil
// }

// func (oe *OptimalExecutor) Run(ctx context.Context) int {
// 	// var ok bool
// 	// var newRelicTxn newrelic.Transaction

// 	// val := ctx.Value(goflow.NewRelicTransactionGoFlowContextKey())
// 	// if val != nil {
// 	// 	newRelicTxn, ok = val.(newrelic.Transaction)
// 	// 	if !ok {
// 	// 		newRelicTxn = nil
// 	// 	}
// 	// }

// 	go oe.manageExecution(ctx)
// }

// func (oe *OptimalExecutor) manageExecution(ctx context.Context) {
// 	defer close(oe.broadcastChannel)

// 	// enigmaLogger.LogDebug(
// 	// 	fullParallelExecutor.EnigmaContext,
// 	// 	fullParallelExecutor.MyName(),
// 	// 	enigmaLogger.String("executable name", fullParallelExecutor.executable.MyName()),
// 	// 	enigmaLogger.Int("executable id", int(fullParallelExecutor.executable.MyID())),
// 	// 	enigmaLogger.Int("number of notifiers", len(fullParallelExecutor.notifiersChannels)),
// 	// )

// 	for _, notifierChannel := range oe.notifierChannels {
// 		select {
// 		case <-notifierChannel:
// 		}
// 	}

// 	for _, executable := range oe.executables {
// 		executable.MarkExecuting()
// 		// TODO: Handle error here to know that this was executed successfully
// 		executable.SyncExecute(ctx)
// 		executable.MarkExecuted()

// 		dependantExecutables := executable.DependantExecutables()
// 		for _, dependantExecutable := range dependantExecutables {
// 			oe.outputChannel <- dependantExecutable
// 		}
// 	}

// 	oe.outputChannel <- nil
// }
