package transformers

import (
	"INGO-goflow/commons"
	"context"
	"errors"
	"reflect"
	"runtime/debug"
)

type Transformer interface {
	Transform(context.Context, commons.OperableDataStore, commons.OperableDataStore) error
}

type SLDLogic func(context.Context, commons.OperableDataStore, commons.OperableDataStore) error

func NewSLDTransformer(logic SLDLogic) *SLDTransformer {
	sldTransformer := &SLDTransformer{
		nil,
		nil,
	}
	sldTransformer.SetLogic(logic)

	return sldTransformer
}

type SLDTransformer struct {
	logic             SLDLogic
	constructionLogic commons.ActorConstructionLogic
}

func (sldTransformer *SLDTransformer) MyType() string {
	return reflect.TypeOf(sldTransformer).String()
}

func (sldTransformer *SLDTransformer) MyName() string {
	return ""
}

func (sldTransformer *SLDTransformer) Describe() map[string]interface{} {
	descriptionMap := make(map[string]interface{})
	descriptionMap["type"] = sldTransformer.MyType()
	descriptionMap["name"] = sldTransformer.MyName()
	descriptionMap["logic type"] = commons.NameOfFunction(sldTransformer.logic)
	descriptionMap["does actor knows its construction?"] = sldTransformer.constructionLogic != nil
	return descriptionMap
}

func (sldTransformer *SLDTransformer) SetLogic(logic SLDLogic) {
	sldTransformer.logic = logic
}

func (sldTransformer *SLDTransformer) SetConstructionLogic(logic commons.ActorConstructionLogic) {
	sldTransformer.constructionLogic = logic
}

func (sldTransformer *SLDTransformer) Tranform(
	ctx context.Context,
	sourceDataStore commons.OperableDataStore,
	destinationDataStore commons.OperableDataStore) error {

	return sldTransformer.logic(ctx, sourceDataStore, destinationDataStore)
}

func (sldTransformer *SLDTransformer) Act(
	ctx context.Context,
	sourceDataStore commons.OperableDataStore,
	destinationDataStore commons.OperableDataStore) (err error) {

	defer func() {
		if r := recover(); r != nil {
			// log runtime errors here
			err = errors.New(string(debug.Stack()))
		}
	}()

	err = sldTransformer.Tranform(ctx, sourceDataStore, destinationDataStore)
	if err != nil {
		// log transformation errors here
		return err
	}

	return nil
}

func (sldTransformer *SLDTransformer) Construct(
	ctx context.Context,
	sourceDataStore commons.OperableDataStore) (err error) {

	if sldTransformer.constructionLogic == nil {
		return nil
	}

	defer func() {
		if r := recover(); r != nil {
			// log runtime errors here
			err = errors.New(string(debug.Stack()))
		}
	}()

	err = sldTransformer.constructionLogic(ctx, sldTransformer, sourceDataStore)
	if err != nil {
		// log contruction errors here
		return err
	}

	return nil
}
