package test

import (
	goflow "INGO-goflow"
	"INGO-goflow/commons"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestSyncLogger(t *testing.T) {
	goflw := goflow.Init(10, 20, 10)
	goFlowLogger := goflw.EnableLogging()
	goFlowLogger.SetLogEnvironment(commons.Development)
	goFlowLogger.AddFileDescriptor(commons.StandardOutput)

	err := goFlowLogger.Build()
	if err != nil {
		t.Error(err)
	}

	goflw.Start()
	goFlowLogger.LogInfo(
		nil,
		"some message",
		goFlowLogger.String("key1", "key1 value"))

	goFlowLogger.LogDebug(
		nil,
		"another message",
		goFlowLogger.String("key2", "key2 value"))

	goFlowLogger.LogError(
		nil,
		errors.New("some error occured"),
		goFlowLogger.String("Errorkey", "some error value"))
}

func TestAsyncLogger(t *testing.T) {
	var waitGRP sync.WaitGroup

	goflw := goflow.Init(10, 20, 10)
	goFlowLogger := goflw.EnableLogging()
	goFlowLogger.SetLogEnvironment(commons.Development)
	goFlowLogger.AddFileDescriptor(commons.StandardOutput)
	goFlowLogger.DoAsyncLogging(100000, &waitGRP)

	err := goFlowLogger.Build()
	if err != nil {
		t.Error(err)
	}

	goflw.Start()
	goFlowLogger.LogInfo(
		nil,
		"some message",
		goFlowLogger.String("key1", "key1 value"))

	goFlowLogger.LogDebug(
		nil,
		"another message",
		goFlowLogger.String("key2", "key2 value"))

	goFlowLogger.LogError(
		nil,
		errors.New("some error occured"),
		goFlowLogger.String("Errorkey", "some error value"))

	// waitGRP.Wait()
	time.Sleep(1 * time.Second)
}
