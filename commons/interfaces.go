package commons

import (
	"context"
)

// validates a node's validity
type NodeValidator func(context.Context, GraphState) (bool, error)

// holds the logic to construct an actor
type ActorConstructionLogic func(context.Context, Actor, OperableDataStore) error

type Describable interface {
	// type
	MyType() string
	// name
	MyName() string
	// gives description
	Describe() map[string]interface{}
}

// major component of a graph node
// this defines what an actor can do including
// its construction logic, act logic
type Actor interface {
	Describable
	SetConstructionLogic(ActorConstructionLogic)
	Construct(context.Context, OperableDataStore) error
	Act(context.Context, OperableDataStore, OperableDataStore) error
}

// datastore inside nodes
// this holds all the data generated while executing
// a node
type OperableDataStore interface {
	Describable
	// Length returns the number of data items
	// present in data store
	Length() int
	// Push pushes data into datastore returning the
	// number of data currently present in updated datastore
	Push(data interface{}) int
	// RetrieveAll returns reference to the buffer
	// avoid modifying the buffer through it, use PeekAll in
	// instead
	RetrieveAll() []interface{}
	// PeekAll returns copy of all the data inside the buffer
	// as a slice
	PeekAll() []interface{}
	// Clear sets the buffer as nil
	Clear()
}

// graph nodes implements this
type Executable interface {
	Describable
	// MyID returns id which uniquely identifies the executable
	MyID() int
	// prepare for execution
	Prepare() error
	// executable is picked for execution
	MarkPicked()
	// am i picked for execution?
	IsPicked() bool
	// executable is executing
	MarkExecuting()
	// tells if the execution is ongoing
	IsExecuting() bool
	// executable's execution finished
	MarkExecuted()
	// tells if the execution finished
	HasExecuted() bool
	// execution logic
	Execute() error
	// executables(if any) dependant on this executable
	DependantExecutables() []Executable
	// executables(if any) on which this executable depends
	DependsOnExecutables() []Executable
}

// depicts an enity which holds group of executable
// graph implements this
type ExecutableCollection interface {
	Describable
	SetStateStart()
	SetStateStop()
	NextEligibleExecutables() []*Executable
}

// takes an executable collection and breaks it down into
// individual executables and execute them in topological
// sort order using Executors
type Controller interface {
	Control(ExecutableCollection) error
}

// executes an executable in a separate go-routine
type Executors interface {
	Run() int
	ClearExecutable()
}
