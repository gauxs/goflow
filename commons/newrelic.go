package commons

import (
	"sync"

	newrelic "github.com/newrelic/go-agent"
)

var newRelicOnce sync.Once
var newRMonitor *NewRelicMonitoring

func EnableNewRelicMonitoring(
	loggerNewRelicTransLogTime int,
	loggerTransactionName string,
	managerNewRelicTransLogTime int,
	managerTransactionName string,
	newRelApp newrelic.Application) *NewRelicMonitoring {
	newRelicOnce.Do(func() {
		newRMonitor = &NewRelicMonitoring{
			newRelicApp:                       newRelApp,
			managerNewRelicTransactionLogTime: managerNewRelicTransLogTime,
			managerNewRelicTransactionName:    managerTransactionName,
			loggerNewRelicTransactionLogTime:  loggerNewRelicTransLogTime,
			loggerNewRelicTransactionName:     loggerTransactionName,
		}
	})

	return newRMonitor
}

func GetNewRelicMonitoring() *NewRelicMonitoring {
	return newRMonitor
}

type NewRelicMonitoring struct {
	newRelicApp                       newrelic.Application
	managerNewRelicTransactionLogTime int
	managerNewRelicTransactionName    string
	loggerNewRelicTransactionLogTime  int
	loggerNewRelicTransactionName     string
}

func (nrm *NewRelicMonitoring) GetNewRelicApplication() newrelic.Application {
	return nrm.newRelicApp
}

func (nrm *NewRelicMonitoring) GetManagerTransactionLogTime() int {
	return nrm.managerNewRelicTransactionLogTime
}

func (nrm *NewRelicMonitoring) GetManagerTransactionName() string {
	return nrm.managerNewRelicTransactionName
}

func (nrm *NewRelicMonitoring) GetLoggerTransactionLogTime() int {
	return nrm.loggerNewRelicTransactionLogTime
}

func (nrm *NewRelicMonitoring) GetLoggerTransactionName() string {
	return nrm.loggerNewRelicTransactionName
}
