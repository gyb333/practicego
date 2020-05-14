package dataType_test

import (
	"../dataType"
	"testing"
)


func TestBaseStruct(t *testing.T) {
	dataType.BaseStruct()
}

func TestReflectStruct(t *testing.T) {
	dataType.ReflectStruct()
}

func TestStructOffset(t *testing.T) {
	dataType.StructOffset()
}