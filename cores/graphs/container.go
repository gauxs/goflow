package graphs

import (
	"context"
	"errors"
	"goflow/commons"
	"reflect"
)

func newContainer(ctx context.Context) *container {
	return &container{
		ctx,
		nil,
		newDataStore(0),
		newDataStore(0),
	}
}

type container struct {
	ctx      context.Context
	actor    commons.Actor
	inbound  commons.OperableDataStore
	outbound commons.OperableDataStore
}

func (container *container) MyType() string {
	return reflect.TypeOf(container).String()
}

func (container *container) MyName() string {
	return ""
}

func (container *container) Describe() map[string]interface{} {
	descriptionMap := make(map[string]interface{})
	descriptionMap["my type"] = container.MyType()
	actoryType := "nil"
	if container.actor != nil {
		actoryType = container.actor.MyType()
	}
	descriptionMap["actor type"] = actoryType
	descriptionMap["non-nil context"] = container.ctx != nil
	return descriptionMap
}

func (container *container) pushIntoInboundDataStore(inputData interface{}) error {
	if container.inbound != nil {
		container.inbound.Push(inputData)
	} else {
		return errors.New("container has nil inbound data store")
	}

	return nil
}

func (container *container) getFromOutboundDataStore() []interface{} {
	return container.outbound.PeekAll()
}

func (container *container) SetActor(actor commons.Actor) {
	container.actor = actor
}

func (container *container) getActor() commons.Actor {
	return container.actor
}

func (container *container) constructActor() error {
	goFlowLogger := commons.GetLogger()

	if goFlowLogger != nil {
		goFlowLogger.LogDebug(
			container.ctx,
			"actor construction",
			goFlowLogger.Description(container),
			goFlowLogger.Int("number of data items in inbound datastore", container.inbound.Length()),
		)
	}

	// log error inside actor construction, pass the error back to inform
	// that actor did not constructed successfully
	return container.actor.Construct(container.ctx, container.inbound)
}

func (container *container) askActorToAct() error {
	goFlowLogger := commons.GetLogger()

	// log error inside actor acting, pass the error back to inform
	// that actor did not acted successfully
	err := container.actor.Act(container.ctx, container.inbound, container.outbound)
	if err != nil {
		return err
	}

	if goFlowLogger != nil {
		goFlowLogger.LogDebug(
			container.ctx,
			"actor's acted successfully",
			goFlowLogger.String("actor name", container.actor.MyName()),
			goFlowLogger.Int("number of data items in inbound datastore", container.inbound.Length()),
			goFlowLogger.Int("number of data items in outbound datastore", container.outbound.Length()),
		)
	}

	return nil
}
