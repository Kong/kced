// This package provides some basic functions for working with yaml nodes.
// The assumption is to never directly encode/decode yaml. Instead, we'll
// convert to/from interface{}.
package yamlbasics

import (
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

//
//
//  parsing
//
//

// FromObject converts the given map[string]interface{} to an yaml node (map).
func FromObject(data map[string]interface{}) (*yaml.Node, error) {
	if data == nil {
		return nil, errors.New("not an object, but <nil>")
	}
	encData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var yNode yaml.Node
	err = yaml.Unmarshal(encData, &yNode)
	if err != nil {
		return nil, err
	}
	if yNode.Kind == yaml.DocumentNode {
		return yNode.Content[0], nil
	}
	return &yNode, nil
}

// ToObject converts the given yaml node to a map[string]interface{}.
func ToObject(data *yaml.Node) (map[string]interface{}, error) {
	if data == nil || data.Kind != yaml.MappingNode {
		return nil, errors.New("data is not a mapping node/object")
	}

	encData, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}
	var jsonData interface{}
	err = json.Unmarshal(encData, &jsonData)
	if err != nil {
		return nil, err
	}
	return jsonData.(map[string]interface{}), nil
}

// ToArray converts the given yaml node to a []interface{}.
func ToArray(data *yaml.Node) ([]interface{}, error) {
	if data == nil || data.Kind != yaml.SequenceNode {
		return nil, errors.New("data is not a sequence node/array")
	}

	encData, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}
	var jsonData interface{}
	err = json.Unmarshal(encData, &jsonData)
	if err != nil {
		return nil, err
	}
	return jsonData.([]interface{}), nil
}

//
//
//  Handling objects and fields
//
//

// NewObject creates a new object node.
func NewObject() *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.MappingNode,
		Tag:   "!!map",
		Style: yaml.FlowStyle,
	}
}

// NewString creates a new string node.
func NewString(value string) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: value,
		Style: yaml.DoubleQuotedStyle,
	}
}

// FindFieldKeyIndex returns the index of the Node that contains the object-Key in the
// targets Content array. If the key is not found, it returns -1.
func FindFieldKeyIndex(targetObject *yaml.Node, key string) int {
	if targetObject.Kind != yaml.MappingNode {
		panic("targetObject is not a mapping node/object")
	}

	for i := 0; i < len(targetObject.Content); i += 2 {
		if targetObject.Content[i].Value == key {
			return i
		}
	}

	return -1
}

// FindFieldValueIndex returns the index of the Node that contains the object-Value in the
// targets Content array. If the value is not found, it returns -1.
func FindFieldValueIndex(targetObject *yaml.Node, key string) int {
	i := FindFieldKeyIndex(targetObject, key)
	if i != -1 {
		i++
	}

	return i
}

// RemoveFieldByIdx removes the key (by its index) and its value from the targetObject.
func RemoveFieldByIdx(targetObject *yaml.Node, idx int) {
	if idx < 0 || idx >= len(targetObject.Content) {
		panic("idx out of bounds")
	}
	targetObject.Content = append(targetObject.Content[:idx], targetObject.Content[idx+2:]...)
}

// RemoveField removes the given key and its value from the targetObject if it exists.
func RemoveField(targetObject *yaml.Node, key string) {
	if i := FindFieldKeyIndex(targetObject, key); i != -1 {
		RemoveFieldByIdx(targetObject, i)
	}
}

// GetFieldValue returns the value of the given key in the targetObject.
// If the key is not found, then nil is returned.
func GetFieldValue(targetObject *yaml.Node, key string) *yaml.Node {
	i := FindFieldValueIndex(targetObject, key)
	if i == -1 {
		return nil
	}
	return targetObject.Content[i]
}

// SetFieldValue sets/overwrites the value of the given key in the targetObject to the
// given value. If value is nil, then the key is removed from the targetObject if it exists.
func SetFieldValue(targetObject *yaml.Node, key string, value *yaml.Node) {
	i := FindFieldKeyIndex(targetObject, key)
	if i == -1 {
		// key not found, so field doesn't exist yet
		if value == nil {
			// nothing to do
			return
		}
		// add the field
		targetObject.Content = append(targetObject.Content, NewString(key), value)
		return
	}

	// key found, so field exists
	if value == nil {
		// remove the field
		RemoveFieldByIdx(targetObject, i)
		return
	}
	targetObject.Content[i+1] = value
}

//
//
//  Handling objects and fields
//
//

// NewArray creates a new array node.
func NewArray() *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.SequenceNode,
		Tag:   "!!seq",
		Style: yaml.FlowStyle,
	}
}

// Append adds the given values to the end of the targetArray. If no values are given,
// then nothing is done.
// If targetArray is nil or not a sequence node, then an error is returned.
// If any of the values are nil, then an error is returned (and the array remains unchanged).
func Append(targetArray *yaml.Node, values ...*yaml.Node) error {
	if targetArray == nil || targetArray.Kind != yaml.SequenceNode {
		return errors.New("targetArray is not a sequence node/array")
	}
	for i, value := range values {
		if value == nil {
			return fmt.Errorf("value at index %d is nil", i)
		}
	}

	targetArray.Content = append(targetArray.Content, values...)
	return nil
}

// AppendSlice appends all entries in a slice to the end of the targetArray.
// If targetArray is nil or not a sequence node, then an error is returned.
// If the slice is nil, then nothing is done.
// If any of the values in the slice are nil, then an error is returned (and the array remains unchanged).
func AppendSlice(targetArray *yaml.Node, values []*yaml.Node) error {
	if targetArray == nil || targetArray.Kind != yaml.SequenceNode {
		return errors.New("targetArray is not a sequence node/array")
	}
	if values == nil {
		return nil
	}

	err := Append(targetArray, values...)
	if err != nil {
		return err
	}
	return nil
}
