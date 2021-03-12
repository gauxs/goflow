package commons

import (
	"errors"
)

const (
	// default slice capacity
	DefaultSliceCapacity = 3
	// new relic transaction context key name
	NewRelicTransactionContextKey = "goflow_newrelic_transaction"
	// uuuid context key name
	UUIDContextKey = "uuid"
	// node name of input handler nodes
	InputHandlerNodeName = "InputHandler"
	// node name of output handler nodes
	OutputHandlerNodeName = "OutputHandler"
	// default logging level keyname
	DefaultLoggingLevelKeyName = "LVL"
	// default logging time keyname
	DefaultLoggingTimeKeyName = "TIME"
	// default logging message keyname
	DefaultLoggingMessageKeyName = "MSG"
	// default logging UUID keyname
	DefaultLoggingUUIDKeyName = "UUID"
	// default async logger channel capacity
	DefaultAsyncLoggerChannelCapacity = 100000
	// default transaction name for log consumer
	DefaultLogConsumerTransactionName = "GoFlowLogConsumer"
	// default transaction name for manager
	DefaultManagerTransactionName = "GoFlowManager"
	// default logger transaction log time(in minutes)
	DefaultLoggerNewRelicTransactionLogTime = 60
	// default manager transaction log time(in minutes)
	DefaultManagerNewRelicTransactionLogTime = 60
)

type goFlowContextKey struct {
	keyName string
}

func (key *goFlowContextKey) String() string {
	return key.keyName
}

func CustomGoFlowContextKey(keyStr string) (*goFlowContextKey, error) {
	if keyStr == "" {
		return nil, errors.New("key length should be non-zero")
	}

	if keyStr == NewRelicTransactionContextKey {
		return nil, errors.New("cannot used a reserved key name")
	}

	return &goFlowContextKey{keyName: keyStr}, nil
}

func NewRelicTransactionGoFlowContextKey() *goFlowContextKey {
	return &goFlowContextKey{keyName: NewRelicTransactionContextKey}
}

func UUIDGoFlowContextKey() *goFlowContextKey {
	return &goFlowContextKey{keyName: UUIDContextKey}
}

type NodeState uint8

// [CAUTION] dont change the order of the states
const (
	DefaultNodeState NodeState = iota
	Unvalidated
	ValidatedUnexecutable
	ValidatedUnexecutableViaNormalisation
	ValidatedExecutable
	PickedForExecution
	ExecutionStarted
	ExecutionFinished
)

func (nodeState NodeState) String() string {
	switch nodeState {
	case DefaultNodeState:
		return "default"
	case Unvalidated:
		return "unvalidated"
	case ValidatedUnexecutable:
		return "validated unexecutable"
	case ValidatedUnexecutableViaNormalisation:
		return "validated unexecutable via normalisation"
	case ValidatedExecutable:
		return "validated executable"
	case PickedForExecution:
		return "picked for execution"
	case ExecutionStarted:
		return "execution started"
	case ExecutionFinished:
		return "execution finished"
	default:
		return "undefined node state"
	}
}

type GraphState uint8

// [CAUTION] don't change the order of the states
const (
	DefaultGraphState GraphState = iota
	Prerun
	Running
	Postrun
)

func (graphState GraphState) String() string {
	switch graphState {
	case DefaultGraphState:
		return "default"
	case Prerun:
		return "pre run"
	case Running:
		return "running"
	case Postrun:
		return "post run"
	default:
		return "undefined graph state"
	}
}

type GoflowRuntimeEnvironment uint8

const (
	UndefinedRuntimeEnvironment GoflowRuntimeEnvironment = iota
	Development
	PreProduction
	ProductionPP
	Production
)

func (enigmaRuntimeEnvironment GoflowRuntimeEnvironment) String() string {
	switch enigmaRuntimeEnvironment {
	case UndefinedRuntimeEnvironment:
		return "undefined"
	case Development:
		return "dev"
	case PreProduction:
		return "pp"
	case ProductionPP:
		return "prodpp"
	case Production:
		return "prod"
	default:
		return ""
	}
}

type LoggingFieldType uint8

const (
	Unknown LoggingFieldType = iota
	Bool
	Int
	String
	Interface
	Error
)

type LogLevel uint8

const (
	// DEBUG typically voluminous, and are usually disabled in production
	DEBUG LogLevel = iota
	// INFO default logging priority
	INFO
	// WARN important than Info, but don't need individual human review
	WARN
	// ERROR high-priority, if an application is running smoothly,
	// it shouldn't generate any error-level logs
	ERROR
	// DPANIC important errors, in development the logger panics after writing the message
	DPANIC
	// PANIC logs a message, then panics
	PANIC
	// FATAL logs a message, then calls os.Exit(1)
	FATAL
)

func (logLevel LogLevel) String() string {
	switch logLevel {
	case DEBUG:
		return "debug"
	case INFO:
		return "info"
	case WARN:
		return "warn"
	case ERROR:
		return "error"
	case DPANIC:
		return "dpanic"
	case PANIC:
		return "panic"
	case FATAL:
		return "fatal"
	default:
		return ""
	}
}

type FileDescriptor uint8

const (
	UndefinedFileDescriptor FileDescriptor = iota
	StandardInput
	StandardOutput
	StandardError
)

func (fileDescriptor FileDescriptor) String() string {
	switch fileDescriptor {
	case UndefinedFileDescriptor:
		return "undefined file descriptor"
	case StandardInput:
		return "stdin"
	case StandardOutput:
		return "stdout"
	case StandardError:
		return "stderr"
	default:
		return ""
	}
}

type GoRoutineStatus uint8

const (
	UndefinedGoRoutineStatus GoRoutineStatus = iota
	Scheduled
	Executing
	Executed
)

func (ws GoRoutineStatus) String() string {
	switch ws {
	case UndefinedGoRoutineStatus:
		return "undefined"
	case Scheduled:
		return "scheduled"
	case Executing:
		return "executing"
	case Executed:
		return "executed"
	default:
		return ""
	}
}
