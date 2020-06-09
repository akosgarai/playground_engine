package boundingobject

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	ExpectedName := "AABB"
	ExpectedParams := make(map[string]float32)
	ExpectedParams["width"] = 1
	ExpectedParams["height"] = 1
	ExpectedParams["length"] = 1

	bo := New(ExpectedName, ExpectedParams)

	if bo.typeName != ExpectedName {
		t.Errorf("Invalid bo type name. Instead of '%s', we have '%s'.\n", ExpectedName, bo.typeName)
	}
	if !reflect.DeepEqual(bo.params, ExpectedParams) {
		t.Errorf("Invalid bo params. Instead of '%v', we have '%v'.\n", ExpectedParams, bo.params)
	}
}
func TestType(t *testing.T) {
	ExpectedName := "AABB"
	ExpectedParams := make(map[string]float32)

	bo := New(ExpectedName, ExpectedParams)

	if bo.Type() != ExpectedName {
		t.Errorf("Invalid bo type name. Instead of '%s', we have '%s'.\n", ExpectedName, bo.Type())
	}
}
func TestParams(t *testing.T) {
	ExpectedName := "AABB"
	ExpectedParams := make(map[string]float32)
	ExpectedParams["width"] = 1
	ExpectedParams["height"] = 1
	ExpectedParams["length"] = 1

	bo := New(ExpectedName, ExpectedParams)

	if !reflect.DeepEqual(bo.Params(), ExpectedParams) {
		t.Errorf("Invalid bo params. Instead of '%v', we have '%v'.\n", ExpectedParams, bo.Params())
	}
}
