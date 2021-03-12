package goflow

import (
	"INGO-goflow/commons"
	"INGO-goflow/cores/controllers"
	"errors"
	"fmt"
	"sync"

	newrelic "github.com/newrelic/go-agent"
)

var goFlw *GoFlow
var goFlowOnce sync.Once

// Initializes a goflow instance with bounds on
// controller, executor while also initializing job
// channel size
func Init(
	controllerWrkrLimit int,
	executorWrkrLimit int,
	controlJobChannelCapacity int) *GoFlow {

	goFlowOnce.Do(func() {
		goFlw = &GoFlow{
			workerMgr: &manager{
				controllerWorkerLimit:        controllerWrkrLimit,
				controllerWorkerLimitChannel: make(chan struct{}, controllerWrkrLimit),
				executorWorkerLimit:          executorWrkrLimit,
				executorWorkerLimitChannel:   make(chan struct{}, executorWrkrLimit),
				controlJobChannelLimit:       controlJobChannelCapacity,
				controlJobChannel:            make(chan *controllers.ControllerJobPacket, controlJobChannelCapacity),
			},
		}
	})

	return goFlw
}

func GetGoFLow() *GoFlow {
	return goFlw
}

type GoFlow struct {
	newRelicEnabled bool
	logEnabled      bool
	workerMgr       *manager
}

func (goFlw *GoFlow) EnableNewRelicMonitoring(
	newRelicApp newrelic.Application,
	managerNewRelicTransactionLogTime int,
	managerTransactionName string,
	loggerNewRelicTransactionLogTime int,
	loggerTransactionName string) error {

	if newRelicApp == nil {
		return errors.New("enabling newrelic with a nil newrelic application")
	}

	if managerNewRelicTransactionLogTime == 0 {
		managerNewRelicTransactionLogTime = commons.DefaultManagerNewRelicTransactionLogTime
	}

	if managerTransactionName == "" {
		managerTransactionName = commons.DefaultManagerTransactionName
	}

	if loggerNewRelicTransactionLogTime == 0 {
		loggerNewRelicTransactionLogTime = commons.DefaultLoggerNewRelicTransactionLogTime
	}

	if loggerTransactionName == "" {
		loggerTransactionName = commons.DefaultLogConsumerTransactionName
	}

	goFlw.newRelicEnabled = true
	commons.EnableNewRelicMonitoring(
		loggerNewRelicTransactionLogTime,
		loggerTransactionName,
		managerNewRelicTransactionLogTime,
		managerTransactionName,
		newRelicApp)

	return nil
}

func (goFlw *GoFlow) EnableLogging() *commons.Logger {
	goFlw.logEnabled = true
	return commons.EnableLogging()
}

func (goFlw *GoFlow) Start() error {
	if goFlw.logEnabled && (commons.GetLogger() == nil || !commons.GetLogger().IsCoreLoggerSet()) {
		return errors.New("logging enabled but logger not built")
	}

	if goFlw.newRelicEnabled && commons.GetNewRelicMonitoring() == nil {
		return errors.New("monitoring enabled but newrelicmonitoring nil")
	}

	var goFlowLogger *commons.Logger
	goFlowLogger = commons.GetLogger()

	if goFlowLogger != nil {
		goFlowLogger.LogInfo(
			nil,
			"starting goflow with following config",
			goFlowLogger.String("newrelic monitoring", fmt.Sprintf("%t", goFlw.newRelicEnabled)),
			goFlowLogger.String("logging", fmt.Sprintf("%t", goFlw.logEnabled)),
		)
	}

	goFlw.workerMgr.Manage()
	return nil
}

func (goFlw *GoFlow) GetControlChannel() chan *controllers.ControllerJobPacket {
	return goFlw.workerMgr.controlJobChannel
}
