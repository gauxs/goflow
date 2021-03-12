package graphs

import (
	"INGO-goflow/commons"
	"reflect"
)

func newDataStore(capacity int) *dataStore {
	if capacity == 0 {
		capacity = commons.DefaultSliceCapacity
	}

	return &dataStore{
		buffer: make([]interface{}, 0, capacity),
	}
}

type dataStore struct {
	buffer []interface{}
}

func (ds *dataStore) MyType() string {
	return reflect.TypeOf(ds).String()
}

func (ds *dataStore) MyName() string {
	return ""
}

func (ds *dataStore) Describe() map[string]interface{} {
	descriptionMap := make(map[string]interface{})
	descriptionMap["my type"] = ds.MyType()
	dataStoreType := "nil"
	if ds.buffer != nil {
		dataStoreType = reflect.TypeOf(ds.buffer).String()
	}
	descriptionMap["data store type"] = dataStoreType
	descriptionMap["number of data elements"] = len(ds.buffer)
	descriptionMap["data elements"] = ds.buffer
	return descriptionMap
}

func (ds *dataStore) Length() int {
	return len(ds.buffer)
}

func (ds *dataStore) Push(data interface{}) int {
	ds.buffer = append(ds.buffer, data)
	return len(ds.buffer)
}

func (ds *dataStore) RetrieveAll() []interface{} {
	return ds.buffer
}

func (ds *dataStore) PeekAll() []interface{} {
	tempBuffer := make([]interface{}, len(ds.buffer))
	copy(tempBuffer, ds.buffer)
	return tempBuffer
}

func (ds *dataStore) Clear() {
	ds.buffer = nil
}
