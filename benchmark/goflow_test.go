package goflow_test

import (
	"context"
	"fmt"
	"goflow"
	"goflow/commons"
	"goflow/cores/controllers"
	"goflow/cores/graphs"
	"goflow/templates/actors/transformers"
	"strconv"
	"sync"
	"testing"
	"time"
)

func BakeCakeInputHandlerActorExecutionLogic(
	ctx context.Context,
	sourceDataStore commons.OperableDataStore,
	destinationDataStore commons.OperableDataStore) error {

	goFlowLogger := commons.GetLogger()
	if goFlowLogger != nil {
		goFlowLogger.LogDebug(
			nil,
			"bake cake input handler done")
	}

	return nil
}

func BakeCakeOutputHandlerActorExecutionLogic(
	ctx context.Context,
	sourceDataStore commons.OperableDataStore,
	destinationDataStore commons.OperableDataStore) error {

	goFlowLogger := commons.GetLogger()
	if goFlowLogger != nil {
		goFlowLogger.LogDebug(
			nil,
			"bake cake output handler done")
	}

	return nil
}

func BakingProcess1(
	ctx context.Context,
	sourceDataStore commons.OperableDataStore,
	destinationDataStore commons.OperableDataStore) error {

	time.Sleep(100 * time.Millisecond)

	goFlowLogger := commons.GetLogger()
	if goFlowLogger != nil {
		goFlowLogger.LogDebug(
			nil,
			"bake cake process-1 done")
	}

	return nil
}

func BakingProcess2(
	ctx context.Context,
	sourceDataStore commons.OperableDataStore,
	destinationDataStore commons.OperableDataStore) error {

	time.Sleep(400 * time.Millisecond)

	goFlowLogger := commons.GetLogger()
	if goFlowLogger != nil {
		goFlowLogger.LogDebug(
			nil,
			"bake cake process-2 done")
	}

	return nil
}

func BakingProcess3(
	ctx context.Context,
	sourceDataStore commons.OperableDataStore,
	destinationDataStore commons.OperableDataStore) error {

	time.Sleep(50 * time.Millisecond)

	goFlowLogger := commons.GetLogger()
	if goFlowLogger != nil {
		goFlowLogger.LogDebug(
			nil,
			"bake cake process-3 done")
	}

	return nil
}

func BakingProcess4(
	ctx context.Context,
	sourceDataStore commons.OperableDataStore,
	destinationDataStore commons.OperableDataStore) error {

	time.Sleep(50 * time.Millisecond)

	goFlowLogger := commons.GetLogger()
	if goFlowLogger != nil {
		goFlowLogger.LogDebug(
			nil,
			"bake cake process-4 done")
	}

	return nil
}

func BakingProcess5(
	ctx context.Context,
	sourceDataStore commons.OperableDataStore,
	destinationDataStore commons.OperableDataStore) error {

	time.Sleep(50 * time.Millisecond)

	goFlowLogger := commons.GetLogger()
	if goFlowLogger != nil {
		goFlowLogger.LogDebug(
			nil,
			"bake cake process-5 done")
	}

	return nil
}

func BakingProcess6(
	ctx context.Context,
	sourceDataStore commons.OperableDataStore,
	destinationDataStore commons.OperableDataStore) error {

	time.Sleep(50 * time.Millisecond)

	goFlowLogger := commons.GetLogger()
	if goFlowLogger != nil {
		goFlowLogger.LogDebug(
			nil,
			"bake cake process-6 done")
	}

	return nil
}

func BakingProcess7(
	ctx context.Context,
	sourceDataStore commons.OperableDataStore,
	destinationDataStore commons.OperableDataStore) error {

	time.Sleep(50 * time.Millisecond)

	goFlowLogger := commons.GetLogger()
	if goFlowLogger != nil {
		goFlowLogger.LogDebug(
			nil,
			"bake cake process-7 done")
	}

	return nil
}

func BakingProcess8(
	ctx context.Context,
	sourceDataStore commons.OperableDataStore,
	destinationDataStore commons.OperableDataStore) error {

	time.Sleep(100 * time.Millisecond)

	goFlowLogger := commons.GetLogger()
	if goFlowLogger != nil {
		goFlowLogger.LogDebug(
			nil,
			"bake cake process-8 done")
	}

	return nil
}

var totalResponseTime int64

func BenchmarkBakeCake(b *testing.B) {
	// setup
	totalResponseTime = 0
	noOfJobs := 10000
	controllerWrkrLimit := 100
	executorWrkrLimit := 20000
	jobChannelCapacity := 10000

	var waitGRP sync.WaitGroup
	goflw := goflow.Init(controllerWrkrLimit, executorWrkrLimit, jobChannelCapacity)
	// goFlowLogger := goflw.EnableLogging()
	// goFlowLogger.SetLogEnvironment(commons.Production)
	// goFlowLogger.AddFileDescriptor(commons.StandardOutput)
	// goFlowLogger.AddAdditionalFields("APP", "MyAppName")
	// err := goFlowLogger.Build()
	// if err != nil {
	// 	b.Error(err)
	// }

	goflw.Start()
	waitGRP.Add(noOfJobs)
	ctrlChan := goflw.GetControlChannel()
	respTimeSlice := make([]int64, noOfJobs+1)
	b.ResetTimer()
	for jobNo := 1; jobNo <= noOfJobs; jobNo++ {
		go MakeJobRequest(b, jobNo, &waitGRP, ctrlChan, respTimeSlice)
	}

	waitGRP.Wait()
	b.StopTimer()
	for _, val := range respTimeSlice {
		totalResponseTime = totalResponseTime + val
	}

	fmt.Println(fmt.Sprintf("ART: %f", float64(totalResponseTime/int64(noOfJobs))/float64(1000)))
	fmt.Println("b.N: ", b.N)

	return
}

func TestBakeCake(t *testing.T) {
	// setup
	totalResponseTime = 0
	noOfJobs := 50000
	controllerWrkrLimit := 60000
	executorWrkrLimit := 200000
	jobChannelCapacity := 50000

	var waitGRP sync.WaitGroup
	goflw := goflow.Init(controllerWrkrLimit, executorWrkrLimit, jobChannelCapacity)
	// goFlowLogger := goflw.EnableLogging()
	// goFlowLogger.SetLogEnvironment(commons.Production)
	// goFlowLogger.AddFileDescriptor(commons.StandardOutput)
	// goFlowLogger.AddAdditionalFields("APP", "MyAppName")
	// err := goFlowLogger.Build()
	// if err != nil {
	// 	b.Error(err)
	// }

	goflw.Start()
	waitGRP.Add(noOfJobs)
	ctrlChan := goflw.GetControlChannel()
	respTimeSlice := make([]int64, noOfJobs+1)
	for jobNo := 1; jobNo <= noOfJobs; jobNo++ {
		go MakeJobRequest(nil, jobNo, &waitGRP, ctrlChan, respTimeSlice)
	}

	waitGRP.Wait()
	for _, val := range respTimeSlice {
		totalResponseTime = totalResponseTime + val
	}

	fmt.Println(fmt.Sprintf("Avg response time: %f", float64(totalResponseTime/int64(noOfJobs))/float64(1000)))

	return
}

func MakeJobRequest(b *testing.B, jobNo int, wg *sync.WaitGroup, controlChnl chan *controllers.ControllerJobPacket, respTimeSlice []int64) {
	defer (*wg).Done()
	ctx := context.WithValue(context.Background(),
		commons.UUIDGoFlowContextKey(), "UUID-BakeCakeBasic"+strconv.Itoa(jobNo))
	bakeCake := graphs.NewTopologyGraph(ctx, "BakeCake"+strconv.Itoa(jobNo))

	inputNode := bakeCake.GetInputHandlerNode()
	sldTransformer := transformers.NewSLDTransformer(BakeCakeInputHandlerActorExecutionLogic)
	inpContainer := inputNode.GetContainer()
	inpContainer.SetActor(sldTransformer)

	outputNode := bakeCake.GetOutputHandlerNode()
	sldTransformer = transformers.NewSLDTransformer(BakeCakeOutputHandlerActorExecutionLogic)
	outContainer := outputNode.GetContainer()
	outContainer.SetActor(sldTransformer)

	bakingProcess1 := bakeCake.NewNode("bakingProcess1")
	sldTransformer = transformers.NewSLDTransformer(BakingProcess1)
	bakingProcess1Container := bakingProcess1.GetContainer()
	bakingProcess1Container.SetActor(sldTransformer)

	bakingProcess2 := bakeCake.NewNode("bakingProcess2")
	sldTransformer = transformers.NewSLDTransformer(BakingProcess2)
	bakingProcess2Container := bakingProcess2.GetContainer()
	bakingProcess2Container.SetActor(sldTransformer)

	bakingProcess3 := bakeCake.NewNode("bakingProcess3")
	sldTransformer = transformers.NewSLDTransformer(BakingProcess3)
	bakingProcess3Container := bakingProcess3.GetContainer()
	bakingProcess3Container.SetActor(sldTransformer)

	bakingProcess4 := bakeCake.NewNode("bakingProcess4")
	sldTransformer = transformers.NewSLDTransformer(BakingProcess4)
	bakingProcess4Container := bakingProcess4.GetContainer()
	bakingProcess4Container.SetActor(sldTransformer)

	bakingProcess5 := bakeCake.NewNode("bakingProcess5")
	sldTransformer = transformers.NewSLDTransformer(BakingProcess5)
	bakingProcess5Container := bakingProcess5.GetContainer()
	bakingProcess5Container.SetActor(sldTransformer)

	bakingProcess6 := bakeCake.NewNode("bakingProcess6")
	sldTransformer = transformers.NewSLDTransformer(BakingProcess6)
	bakingProcess6Container := bakingProcess6.GetContainer()
	bakingProcess6Container.SetActor(sldTransformer)

	bakingProcess7 := bakeCake.NewNode("bakingProcess7")
	sldTransformer = transformers.NewSLDTransformer(BakingProcess7)
	bakingProcess7Container := bakingProcess7.GetContainer()
	bakingProcess7Container.SetActor(sldTransformer)

	bakingProcess8 := bakeCake.NewNode("bakingProcess8")
	sldTransformer = transformers.NewSLDTransformer(BakingProcess8)
	bakingProcess8Container := bakingProcess8.GetContainer()
	bakingProcess8Container.SetActor(sldTransformer)

	bakingProcess1.AddDependant(bakingProcess4)
	bakingProcess2.AddDependant(bakingProcess8)
	bakingProcess3.AddDependant(bakingProcess4, bakingProcess5)
	bakingProcess4.AddDependant(bakingProcess6)
	bakingProcess5.AddDependant(bakingProcess7)
	bakingProcess6.AddDependant(bakingProcess8)
	bakingProcess7.AddDependant(bakingProcess8)

	var replyErr error
	replyChan := make(chan error)
	start := time.Now()
	// b.StartTimer()
	controlChnl <- &controllers.ControllerJobPacket{
		ReplyChannel: replyChan,
		Graph:        bakeCake,
	}
	// b.StopTimer()
	replyErr = <-replyChan
	respTimeSlice[jobNo] = time.Since(start).Milliseconds()
	if replyErr != nil {
		fmt.Println(replyErr)
	}

	return
}
