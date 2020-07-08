package form

import (
	"strconv"
)

// FormItem is the base of the form inputs. Every form item has a label
// that has to be string.
type FormItem struct {
	label string
}

// BoolValue return the string representation of the bool item.
func (fi *FormItem) BoolValue(v bool) string {
	return strconv.FormatBool(v)
}

// IntegerValue returns the string representation of the bool item.
func (fi *FormItem) IntegerValue(v int) string {
	return strconv.Itoa(v)
}

// FloatValue returns the string representation of the float number.
func (fi *FormItem) FloatValue(v float32) string {
	return strconv.FormatFloat(float64(v), 'G', -1, 32)
}

func newFormItem(l string) *FormItem {
	return &FormItem{
		label: l,
	}
}
