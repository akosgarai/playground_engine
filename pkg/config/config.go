package config

import (
	"errors"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	ValueTypeInt = iota
	ValueTypeInt64
	ValueTypeFloat
	ValueTypeText
	ValueTypeVector
	ValueTypeBool
)

var (
	InvalidType = errors.New("Invalid type")
)

type ConfigItem struct {
	key          string
	label        string
	description  string
	valueType    int
	defaultValue interface{}
	currentValue interface{}
}

// NewConfigItem returns a ConfigItem. It gets the valueType from the type of the defaultValue.
// In case of unhandled type, it panics.
func NewConfigItem(key, label, description string, defaultValue interface{}) *ConfigItem {
	var valueType int
	switch defaultValue.(type) {
	case string:
		valueType = ValueTypeText
		break
	case int:
		valueType = ValueTypeInt
		break
	case int64:
		valueType = ValueTypeInt64
		break
	case float32:
		valueType = ValueTypeFloat
		break
	case bool:
		valueType = ValueTypeBool
		break
	case mgl32.Vec3:
		valueType = ValueTypeVector
		break
	default:
		panic("Unhandled value type.")
	}
	return &ConfigItem{
		key:          key,
		label:        label,
		description:  description,
		valueType:    valueType,
		defaultValue: defaultValue,
		currentValue: defaultValue,
	}
}

// GetLabel returns the label of the ConfigItem
func (ci *ConfigItem) GetLabel() string {
	return ci.label
}

// GetDescription returns the description of the ConfigItem
func (ci *ConfigItem) GetDescription() string {
	return ci.description
}

// GetDefaultValue returns the default value of the ConfigItem
func (ci *ConfigItem) GetDefaultValue() interface{} {
	return ci.defaultValue
}

// GetCurrentValue returns the current value of the ConfigItem
func (ci *ConfigItem) GetCurrentValue() interface{} {
	return ci.currentValue
}

// GetValueType returns the value type integer of the ConfigItem
func (ci *ConfigItem) GetValueType() int {
	return ci.valueType
}

// SetCurrentValue gets a value interface input. It checks that the current value
// type and the default value type is same or not. In case of difference, the value
// will not be updated.
func (ci *ConfigItem) SetCurrentValue(v interface{}) error {
	switch v.(type) {
	case string:
		if ci.valueType != ValueTypeText {
			return InvalidType
		}
		break
	case int:
		if ci.valueType != ValueTypeInt {
			return InvalidType
		}
		break
	case int64:
		if ci.valueType != ValueTypeInt64 {
			return InvalidType
		}
		break
	case float32:
		if ci.valueType != ValueTypeFloat {
			return InvalidType
		}
		break
	case bool:
		if ci.valueType != ValueTypeBool {
			return InvalidType
		}
		break
	case mgl32.Vec3:
		if ci.valueType != ValueTypeVector {
			return InvalidType
		}
		break
	default:
		return InvalidType
	}
	ci.currentValue = v
	return nil
}

type Config map[string]*ConfigItem

func New() Config {
	return make(map[string]*ConfigItem)
}

// AddConfig inserts a new config item to the config. In case of invalid default value type
// a panic will occure.
func (c Config) AddConfig(key, label, description string, defaultValue interface{}) {
	c[key] = NewConfigItem(key, label, description, defaultValue)
}

// SetCurrentValue sets the current value of the config item for the given key.
func (c Config) SetCurrentValue(key string, currentValue interface{}) error {
	return c[key].SetCurrentValue(currentValue)
}
