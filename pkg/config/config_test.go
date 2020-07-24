package config

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestNewConfigItemValidValue(t *testing.T) {
	testData := []struct {
		key       string
		label     string
		desc      string
		value     interface{}
		validator interface{}
		vtype     int
	}{
		{"key1", "label 1", "description 1", int(7), func(i int) bool { return i >= 0 && i <= 10 }, ValueTypeInt},
		{"key2", "label 2", "description 2", int64(7), nil, ValueTypeInt64},
		{"key3", "label 3", "description 3", float32(7), nil, ValueTypeFloat},
		{"key4", "label 4", "description 4", string("seven"), nil, ValueTypeText},
		{"key5", "label 5", "description 5", true, nil, ValueTypeBool},
		{"key6", "label 6", "description 6", mgl32.Vec3{0, 0, 0}, nil, ValueTypeVector},
	}
	for _, tt := range testData {
		fi := NewConfigItem(tt.key, tt.label, tt.desc, tt.value, tt.validator)
		if fi.key != tt.key {
			t.Errorf("Invalid key value. Instead of '%s', we have '%s'.", tt.key, fi.key)
		}
		if fi.label != tt.label {
			t.Errorf("Invalid label value. Instead of '%s', we have '%s'.", tt.label, fi.label)
		}
		if fi.GetLabel() != tt.label {
			t.Errorf("Invalid GetLabel value. Instead of '%s', we have '%s'.", tt.label, fi.GetLabel())
		}
		if fi.description != tt.desc {
			t.Errorf("Invalid description value. Instead of '%s', we have '%s'.", tt.desc, fi.description)
		}
		if fi.GetDescription() != tt.desc {
			t.Errorf("Invalid GetDescription value. Instead of '%s', we have '%s'.", tt.desc, fi.GetDescription())
		}
		if fi.defaultValue != tt.value {
			t.Errorf("Invalid default value. Instead of '%v', we have '%v'.", tt.value, fi.defaultValue)
		}
		if fi.GetDefaultValue() != tt.value {
			t.Errorf("Invalid GetDefultValue. Instead of '%v', we have '%v'.", tt.value, fi.GetDefaultValue())
		}
		if fi.GetValueType() != tt.vtype {
			t.Errorf("Invalid GetValueType. Instead of '%d', we have '%d'.", tt.vtype, fi.GetValueType())
		}
		if fi.GetValidatorFunction() != nil {
			_ = fi.GetValidatorFunction()
		}
		if fi.IsConfigOf(tt.key) == false {
			t.Error("It supposed to be the config of its own key")
		}
		if fi.IsConfigOf(tt.key+"wrong") == true {
			t.Error("It isn't supposed to be the config of the wrong key")
		}
	}
}
func TestNewConfigItemInvalidValue(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Invalid value should lead to panic.")
			}
		}()
		_ = NewConfigItem("key", "label", "desc", float64(3.4), nil)
	}()
}
func TestConfigItemSetCurrentValue(t *testing.T) {
	testData := []struct {
		key       string
		label     string
		desc      string
		value     interface{}
		validator interface{}
		newValue  interface{}
	}{
		{"key1", "label 1", "description 1", int(7), nil, int(9)},
		{"key2", "label 2", "description 2", int64(7), nil, int64(9)},
		{"key3", "label 3", "description 3", float32(7), nil, float32(9)},
		{"key4", "label 4", "description 4", string("seven"), nil, string("nine")},
		{"key5", "label 5", "description 5", true, nil, false},
		{"key6", "label 6", "description 6", mgl32.Vec3{0, 0, 0}, nil, mgl32.Vec3{1, 1, 1}},
	}
	for _, tt := range testData {
		fi := NewConfigItem(tt.key, tt.label, tt.desc, tt.value, tt.validator)
		err := fi.SetCurrentValue(tt.newValue)
		if err != nil {
			t.Errorf("Value set should be successful, but we have the following error: '%#v'.", err)
		}
	}
	// invalid values
	for _, tt := range testData {
		fi := NewConfigItem(tt.key, tt.label, tt.desc, tt.value, tt.validator)
		err := fi.SetCurrentValue(float64(9))
		if err == nil {
			t.Error("Value set shouldn't be successful")
		}
	}
	errorData := []struct {
		validKey string
		value    interface{}
	}{
		{"key1", int(11)},
		{"key2", int64(11)},
		{"key3", float32(11)},
		{"key4", string("eleven")},
		{"key5", false},
		{"key6", mgl32.Vec3{0, 1, 0}},
	}
	for _, ed := range errorData {
		for _, tt := range testData {
			fi := NewConfigItem(tt.key, tt.label, tt.desc, tt.value, tt.validator)
			err := fi.SetCurrentValue(ed.value)
			if tt.key == ed.validKey {
				if err != nil {
					t.Errorf("Setting value should be successful for '%s'. We got the following error: '%#v'.", ed.validKey, err)
				}
				if fi.GetCurrentValue() != ed.value {
					t.Error("Invalid value.")
				}
			}
			if tt.key != ed.validKey && err == nil {
				t.Errorf("Setting value should be failed for '%s'.", ed.validKey)
			}
		}
	}
}
func TestNewConfig(t *testing.T) {
	conf := New()
	if len(conf) > 0 {
		t.Errorf("Config should be empty, but it contains '%d' items.", len(conf))
	}
}
func TestConfigAddConfigValidValues(t *testing.T) {
	testData := []struct {
		key       string
		label     string
		desc      string
		value     interface{}
		validator interface{}
	}{
		{"key1", "label 1", "description 1", int(7), nil},
		{"key2", "label 2", "description 2", int64(7), nil},
		{"key3", "label 3", "description 3", float32(7), nil},
		{"key4", "label 4", "description 4", string("seven"), nil},
		{"key5", "label 5", "description 5", true, nil},
		{"key6", "label 6", "description 6", mgl32.Vec3{0, 0, 0}, nil},
	}
	conf := New()
	for i, tt := range testData {
		conf.AddConfig(tt.key, tt.label, tt.desc, tt.value, tt.validator)
		if len(conf) != i+1 {
			t.Errorf("Invalid config length. It supposed to be '%d', but it is '%d'.", i+1, len(conf))
		}
	}
}
func TestNewConfigInvalidValue(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Invalid value should lead to panic.")
			}
		}()
		conf := New()
		conf.AddConfig("key", "label", "desc", float64(3.4), nil)
	}()
}
func TestConfigSetCurrentValue(t *testing.T) {
	testData := []struct {
		key       string
		label     string
		desc      string
		value     interface{}
		newValue  interface{}
		validator interface{}
	}{
		{"key1", "label 1", "description 1", int(7), int(9), nil},
		{"key2", "label 2", "description 2", int64(7), int64(9), nil},
		{"key3", "label 3", "description 3", float32(7), float32(9), nil},
		{"key4", "label 4", "description 4", string("seven"), string("nine"), nil},
		{"key5", "label 5", "description 5", true, false, nil},
		{"key6", "label 6", "description 6", mgl32.Vec3{0, 0, 0}, mgl32.Vec3{1, 1, 1}, nil},
	}
	conf := New()
	for _, tt := range testData {
		conf.AddConfig(tt.key, tt.label, tt.desc, tt.value, tt.validator)
		err := conf.SetCurrentValue(tt.key, tt.newValue)
		if err != nil {
			t.Errorf("Value set should be successful, but we have the following error: '%#v'.", err)
		}
		// invalid values
		err = conf.SetCurrentValue(tt.key, float64(9))
		if err == nil {
			t.Error("Value set shouldn't be successful")
		}
	}
}
