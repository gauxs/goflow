package goflow

import (
	"INGO-goflow/commons"
	"INGO-goflow/cores/controllers"
	"errors"
	"strconv"
	"time"

	newrelic "github.com/newrelic/go-agent"
)

type manager struct {
	controllerWorkerLimit        int
	controllerWorkerLimitChannel chan struct{}
	executorWorkerLimit          int
	executorWorkerLimitChannel   chan struct{}
	controlJobChannelLimit       int
	controlJobChannel            chan *controllers.ControllerJobPacket
}

func (mgr *manager) ControlChannel() chan *controllers.ControllerJobPacket {
	return mgr.controlJobChannel
}

func (mgr *manager) Manage() {
	go mgr.monitor()
}

func (mgr *manager) monitor() {
	var jobNumber int
	var ticker *time.Ticker
	var goFlowLogger *commons.Logger
	var managerTransactionName string
	var managerTransactionLogTime int
	var newRelicApp newrelic.Application
	var managerRoutine newrelic.Transaction

	goFlowLogger = commons.GetLogger()
	newRMonitor := commons.GetNewRelicMonitoring()
	if newRelicApp != nil {
		newRelicApp = newRMonitor.GetNewRelicApplication()
		managerTransactionName = newRMonitor.GetManagerTransactionName()
		managerTransactionLogTime = newRMonitor.GetManagerTransactionLogTime()

		managerRoutine = newRelicApp.StartTransaction(managerTransactionName, nil, nil)
	} else {
		// give a positive value initially, ticker will stop automatically
		managerTransactionLogTime = 1
	}

	ticker = time.NewTicker(time.Duration(managerTransactionLogTime) * time.Minute)
	for {
		select {
		case job := <-mgr.controlJobChannel:
			if job.Graph == nil {
				job.ReplyChannel <- errors.New("nil graph in ControllerJobPacket")
				continue
			}

			jobNumber++
			if goFlowLogger != nil {
				goFlowLogger.LogInfo(
					job.Graph.GetContext(),
					"job "+strconv.Itoa(jobNumber),
					goFlowLogger.Int("job's in channel", len(mgr.controlJobChannel)),
				)
			}

			// construct the graph
			job.Graph.Construct()

			if goFlowLogger != nil {
				goFlowLogger.LogDebug(
					job.Graph.GetContext(),
					"topology graph details",
					goFlowLogger.Description(job.Graph),
				)
			}

			// send empty struct in controller limit channel
			mgr.controllerWorkerLimitChannel <- struct{}{}

			// controller will die out after tasks complete
			// and we have the upper cap using channels
			go controllers.Control(job, mgr.controllerWorkerLimitChannel, mgr.executorWorkerLimitChannel)

		case <-ticker.C:
			// TODO: optimize this
			if managerRoutine != nil {
				managerRoutine.End()
				managerRoutine = newRelicApp.StartTransaction(managerTransactionName, nil, nil)
			} else {
				if ticker != nil {
					ticker.Stop()
				}
			}
		}
	}
}
