package executors

// import (
// 	goflow "INGO-goflow"
// 	"reflect"
// 	"sync"
// )

// func NewSubOptimalExecutor() (*SubOptimalExecutor, error) {
// 	return &SubOptimalExecutor{
// 		make([]goflow.Executable, 0),
// 	}, nil
// }

// type SubOptimalExecutor struct {
// 	executables []goflow.Executable
// }

// func (soe *SubOptimalExecutor) MyType() string {
// 	return reflect.TypeOf(soe).String()
// }

// func (soe *SubOptimalExecutor) MyName() string {
// 	return ""
// }

// func (soe *SubOptimalExecutor) AddExecutable(executables ...goflow.Executable) {
// 	soe.executables = append(soe.executables, executables...)
// }

// func (soe *SubOptimalExecutor) ClearExecutable() {
// 	soe.executables = make([]goflow.Executable, 0)
// }

// func (soe *SubOptimalExecutor) Run() int {
// 	var waitGRP sync.WaitGroup
// 	var numberOfExecutbale int = len(soe.executables)
// 	var numberOfNodeExecutedSuccessfully int = numberOfExecutbale
// 	var executablesExecutionStatusChan = make(chan bool, numberOfExecutbale)

// 	waitGRP.Add(numberOfExecutbale)
// 	for index := 0; index < numberOfExecutbale; index++ {
// 		soe.executables[index].MarkExecuting()
// 		go soe.executables[index].AsyncExecute(&waitGRP, executablesExecutionStatusChan, semiParallelExecutor.GetNewRelicTransaction().NewGoroutine())
// 		soe.executables[index].MarkExecuted()
// 	}

// 	waitGRP.Wait()
// 	close(executablesExecutionStatusChan)
// 	for exectedSuccessfully := range executablesExecutionStatusChan {
// 		if !exectedSuccessfully {
// 			numberOfNodeExecutedSuccessfully--
// 		}
// 	}

// 	return numberOfNodeExecutedSuccessfully
// }
