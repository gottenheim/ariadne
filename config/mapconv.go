package config

import (
	"fmt"
	"strings"
)

type MapStr map[string]interface{}
type MapItf map[interface{}]interface{}

func cleanUpInterfaceArray(in []interface{}) []interface{} {
	result := make([]interface{}, len(in))
	for i, v := range in {
		result[i] = cleanUpMapValue(v)
	}
	return result
}

func cleanUpInterfaceMap(in MapItf) MapStr {
	result := make(MapStr)
	for k, v := range in {
		result[fmt.Sprintf("%v", k)] = cleanUpMapValue(v)
	}
	return result
}

func cleanUpMapValue(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return cleanUpInterfaceArray(v)
	case MapItf:
		return cleanUpInterfaceMap(v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (m MapStr) deepClone() MapStr {
	clone := MapStr{}
	deepCloneMap(m, clone)
	return clone
}

func deepCloneMap(src MapStr, dest MapStr) {
	for key, value := range src {
		switch src[key].(type) {
		case MapStr:
			dest[key] = MapStr{}
			deepCloneMap(src[key].(MapStr), dest[key].(MapStr))
		case map[string]interface{}:
			dest[key] = MapStr{}
			deepCloneMap(src[key].(map[string]interface{}), dest[key].(MapStr))
		default:
			dest[key] = value
		}
	}
}

func (m MapStr) RemovePath(path string) {
	if len(path) == 0 {
		return
	}

	pathItems := strings.Split(path, "/")
	lastIndex := len(pathItems) - 1
	lastItem := pathItems[lastIndex]

	pathOwner := m.FindPath(strings.Join(pathItems[:lastIndex], "/"))
	if pathOwner == nil {
		return
	}

	pathOwnerMap, isMap := pathOwner.(MapStr)
	if !isMap {
		return
	}

	delete(pathOwnerMap, lastItem)
}

func (m MapStr) FindPath(path string) interface{} {
	if len(path) == 0 {
		return m
	}

	var head string
	var tail string
	index := strings.Index(path, "/")
	if index == -1 {
		head = path
	} else {
		head = path[0:index]
		tail = path[index+1:]
	}

	value, ok := m.lookupValueByKey(head)

	if !ok || value == nil {
		return nil
	}

	if len(tail) == 0 {
		return value
	}

	mapValue, ok := value.(MapStr)
	if !ok || mapValue == nil {
		mapValue, ok = value.(map[string]interface{})
		if !ok || mapValue == nil {
			return nil
		}
	}

	return mapValue.FindPath(tail)
}

func (m MapStr) lookupValueByKey(requiredKey string) (interface{}, bool) {
	// find simple key
	value, ok := m[requiredKey]

	if ok {
		return value, ok
	}

	// find composite key (a, b, c)
	for mapKey, mapValue := range m {
		if strings.Contains(mapKey, ",") {
			mapKeyParts := strings.Split(mapKey, ",")

			for _, mapKeyPart := range mapKeyParts {
				if strings.TrimSpace(mapKeyPart) == requiredKey {
					return mapValue, true
				}
			}
		}
	}

	return nil, false
}

func (m MapStr) SetValue(path string, value interface{}) interface{} {
	var head string
	var tail string
	index := strings.Index(path, "/")
	if index == -1 {
		head = path
	} else {
		head = path[0:index]
		tail = path[index+1:]
	}

	value, ok := m[head]
	if !ok || value == nil {
		return nil
	}

	if len(tail) == 0 {
		return value
	}

	mapValue, ok := value.(MapStr)
	if !ok || mapValue == nil {
		return nil
	}

	return mapValue.FindPath(tail)
}

func (m MapStr) MakePath(path string) MapStr {
	if len(path) == 0 {
		return m
	}

	var head string
	var tail string
	index := strings.Index(path, "/")
	if index == -1 {
		head = path
	} else {
		head = path[0:index]
		tail = path[index+1:]
	}

	value, ok := m[head]
	if !ok || value == nil {
		// Create new map if value doesn't exist
		value = MapStr{}
	}

	mapValue, ok := value.(MapStr)
	if !ok || mapValue == nil {
		// Replace existed value with map if it exists
		mapValue = MapStr{}
	}

	m[head] = mapValue

	if len(tail) == 0 {
		return mapValue
	}

	return mapValue.MakePath(tail)
}

func (m MapStr) GetKeys() []string {
	keys := make([]string, len(m))

	index := 0
	for key, _ := range m {
		keys[index] = key
		index++
	}

	return keys
}
