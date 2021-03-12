package test

import (
	goflow "INGO-goflow"
	"INGO-goflow/commons"
	"INGO-goflow/cores/controllers"
	"INGO-goflow/cores/graphs"
	"INGO-goflow/templates/actors/transformers"
	"context"
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

func TestBasic(t *testing.T) {
	ctx := context.WithValue(context.Background(), commons.UUIDGoFlowContextKey(), "UUID-TestBasic")
	bakeCake := graphs.NewTopologyGraph(ctx, "BakeCake")

	goflw := goflow.Init(10, 20, 10)
	goFlowLogger := goflw.EnableLogging()
	goFlowLogger.SetLogEnvironment(commons.Development)
	goFlowLogger.AddFileDescriptor(commons.StandardOutput)
	goFlowLogger.AddAdditionalFields("APP", "MyAppName")

	err := goFlowLogger.Build()
	if err != nil {
		t.Error(err)
	}

	goflw.Start()
	inputNode := bakeCake.GetInputHandlerNode()
	outputNode := bakeCake.GetOutputHandlerNode()

	sldTransformer := transformers.NewSLDTransformer(BakeCakeInputHandlerActorExecutionLogic)
	inpContainer := inputNode.GetContainer()
	inpContainer.SetActor(sldTransformer)

	sldTransformer = transformers.NewSLDTransformer(BakeCakeOutputHandlerActorExecutionLogic)
	outContainer := outputNode.GetContainer()
	outContainer.SetActor(sldTransformer)

	inputNode.AddDependant(outputNode)

	var replyErr error
	replyChan := make(chan error)
	controlChnl := goflw.GetControlChannel()
	controlChnl <- &controllers.ControllerJobPacket{
		ReplyChannel: replyChan,
		Graph:        bakeCake,
	}

	replyErr = <-replyChan
	if replyErr != nil {
		t.Error(replyErr)
	}

	return
}

func TestBakeCakeBasic(t *testing.T) {
	// setup
	goflw := goflow.Init(100, 200, 100)
	goFlowLogger := goflw.EnableLogging()
	goFlowLogger.SetLogEnvironment(commons.Development)
	goFlowLogger.AddFileDescriptor(commons.StandardOutput)
	goFlowLogger.AddAdditionalFields("APP", "MyAppName")

	err := goFlowLogger.Build()
	if err != nil {
		t.Error(err)
	}

	goflw.Start()
	ctx := context.WithValue(context.Background(),
		commons.UUIDGoFlowContextKey(), "UUID-BakeCakeBasic")
	bakeCake := graphs.NewTopologyGraph(ctx, "BakeCake")

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
	controlChnl := goflw.GetControlChannel()
	controlChnl <- &controllers.ControllerJobPacket{
		ReplyChannel: replyChan,
		Graph:        bakeCake,
	}

	replyErr = <-replyChan
	if replyErr != nil {
		t.Error(replyErr)
	}

	return
}
