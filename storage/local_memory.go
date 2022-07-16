package storage

import (
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/BorisIosifov/mend-home-assignment/object"
)

var lmMu sync.RWMutex

type LocalMemory struct {
	// cache[object_type][object_id] = Object
	Cache map[string]map[int]object.Object
}

func PrepareLocalMemory() (lm LocalMemory, err error) {
	cache := make(map[string]map[int]object.Object)
	cache["books"] = make(map[int]object.Object)
	cache["cars"] = make(map[int]object.Object)
	lm = LocalMemory{
		Cache: cache,
	}
	log.Print("LocalMemory storage is ready")

	return lm, nil
}

func (lm LocalMemory) GetList(objectType string) (objects []object.Object, err error) {
	lmMu.RLock()
	defer lmMu.RUnlock()

	objects = make([]object.Object, 0, len(lm.Cache[objectType]))
	for _, object := range lm.Cache[objectType] {
		objects = append(objects, object)
	}
	sort.Slice(objects, func(i, j int) bool { return objects[i].GetID() < objects[j].GetID() })
	return objects, err
}

func (lm LocalMemory) Get(objectType string, ID int) (result object.Object, isNotFound bool, err error) {
	lmMu.RLock()
	defer lmMu.RUnlock()

	result, ok := lm.Cache[objectType][ID]
	if !ok {
		return nil, true, fmt.Errorf("Object %s with id %d not found", objectType, ID)
	}
	return result, false, err
}

func (lm LocalMemory) Post(objectType string, obj object.Object) (result object.Object, err error) {
	lmMu.Lock()
	defer lmMu.Unlock()

	maxID := 0
	for ID, _ := range lm.Cache[objectType] {
		if ID > maxID {
			maxID = ID
		}
	}

	// could be a race condition, mutex should be used
	newID := maxID + 1
	obj.SetID(newID)
	lm.Cache[objectType][newID] = obj
	return obj, err
}

func (lm LocalMemory) Put(objectType string, ID int, obj object.Object) (result object.Object, isNotFound bool, err error) {
	lmMu.Lock()
	defer lmMu.Unlock()

	_, ok := lm.Cache[objectType][ID]
	if !ok {
		return nil, true, fmt.Errorf("Object %s with id %d not found", objectType, ID)
	}

	obj.SetID(ID)
	lm.Cache[objectType][ID] = obj
	return obj, false, err
}

func (lm LocalMemory) Delete(objectType string, ID int) (isNotFound bool, err error) {
	lmMu.Lock()
	defer lmMu.Unlock()

	_, ok := lm.Cache[objectType][ID]
	if !ok {
		return true, fmt.Errorf("Object %s with id %d not found", objectType, ID)
	}

	delete(lm.Cache[objectType], ID)
	return false, err
}
