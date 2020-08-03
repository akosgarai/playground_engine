package screen

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/akosgarai/playground_engine/pkg/config"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/store"
	"github.com/akosgarai/playground_engine/pkg/testhelper"

	"github.com/go-gl/mathgl/mgl32"
)

func TestNewFormScreenBuilder(t *testing.T) {
	builder := NewFormScreenBuilder()
	if builder.headerLabel != "Default label" {
		t.Error("Invalid default label")
	}
	if builder.wrapper != nil {
		t.Error("Wrapper supposed to be nil by default")
	}
	if builder.charset != nil {
		t.Error("Charset supposed to be nil by default")
	}
}
func TestFormScreenBuilderSetHeaderLabel(t *testing.T) {
	builder := NewFormScreenBuilder()
	label := "new label"
	builder.SetHeaderLabel(label)
	if builder.headerLabel != label {
		t.Errorf("Invalid header label. Instead of '%s', we have '%s'.", label, builder.headerLabel)
	}
}
func TestFormScreenBuilderSetWrapper(t *testing.T) {
	builder := NewFormScreenBuilder()
	builder.SetWrapper(wrapperMock)
	if builder.wrapper != wrapperMock {
		t.Error("Invalid wrapper")
	}
}
func TestFormScreenBuilderSetWindowSize(t *testing.T) {
	builder := NewFormScreenBuilder()
	wW := float32(800)
	wH := float32(800)
	builder.SetWindowSize(wW, wH)
	if builder.windowWidth != wW {
		t.Errorf("Invalid window width. Instead of '%f', we have '%f'.", wW, builder.windowWidth)
	}
	if builder.windowHeight != wH {
		t.Errorf("Invalid window height. Instead of '%f', we have '%f'.", wH, builder.windowHeight)
	}
}
func TestFormScreenBuilderSetHeaderLabelColor(t *testing.T) {
	builder := NewFormScreenBuilder()
	color := mgl32.Vec3{1, 0, 1}
	builder.SetHeaderLabelColor(color)
	if builder.headerLabelColor != color {
		t.Error("Invalid header label color.")
	}
}
func TestFormScreenBuilderSetFormItemLabelColor(t *testing.T) {
	builder := NewFormScreenBuilder()
	color := mgl32.Vec3{1, 0, 1}
	builder.SetFormItemLabelColor(color)
	if builder.formItemLabelColor != color {
		t.Error("Invalid form item label color.")
	}
}
func TestFormScreenBuilderSetFormItemInputColor(t *testing.T) {
	builder := NewFormScreenBuilder()
	color := mgl32.Vec3{1, 0, 1}
	builder.SetFormItemInputColor(color)
	if builder.formItemInputColor != color {
		t.Error("Invalid form item input color.")
	}
}
func TestFormScreenBuilderSetClearColor(t *testing.T) {
	builder := NewFormScreenBuilder()
	color := mgl32.Vec3{1, 0, 1}
	builder.SetClearColor(color)
	if builder.clearColor != color {
		t.Error("Invalid background color.")
	}
}
func TestFormScreenBuilderSetConfig(t *testing.T) {
	conf := config.New()
	builder := NewFormScreenBuilder()
	builder.SetConfig(conf)
	if !reflect.DeepEqual(builder.config, conf) {
		t.Error("Invalid configuration.")
	}
}
func TestFormScreenBuilderSetConfigOrder(t *testing.T) {
	builder := NewFormScreenBuilder()
	order := []string{"o1", "p1", "l1"}
	builder.SetConfigOrder(order)
	if !reflect.DeepEqual(builder.configOrder, order) {
		t.Error("Invalid order.")
	}
}
func TestFormScreenBuilderSetCharset(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	runtime.LockOSThread()
	testhelper.GlfwInit()
	defer testhelper.GlfwTerminate()
	wrapperReal.InitOpenGL()
	builder := NewFormScreenBuilder()
	charset, err := model.LoadCharset("./assets/fonts/Desyrel/desyrel.ttf", 32, 127, 40.0, 72, wrapperReal)
	if err != nil {
		t.Errorf("Error during load charset: %#v.", err)
	}
	builder.SetCharset(charset)
	if builder.charset != charset {
		t.Error("Invalid charset.")
	}
}
func TestFormScreenBuilderAddConfigBool(t *testing.T) {
	builder := NewFormScreenBuilder()
	if len(builder.config) != 0 {
		t.Error("Invalid initial config length.")
	}
	confLabel := "label"
	confDesc := "desc"
	value := true
	builder.AddConfigBool(confLabel, confDesc, value)
	if len(builder.config) != 1 {
		t.Error("Invalid config length.")
	}
}
func TestFormScreenBuilderAddConfigInt(t *testing.T) {
	builder := NewFormScreenBuilder()
	if len(builder.config) != 0 {
		t.Error("Invalid initial config length.")
	}
	confLabel := "label"
	confDesc := "desc"
	value := 3
	builder.AddConfigInt(confLabel, confDesc, value, nil)
	if len(builder.config) != 1 {
		t.Error("Invalid config length.")
	}
}
func TestFormScreenBuilderAddConfigInt64(t *testing.T) {
	builder := NewFormScreenBuilder()
	if len(builder.config) != 0 {
		t.Error("Invalid initial config length.")
	}
	confLabel := "label"
	confDesc := "desc"
	value := int64(3)
	builder.AddConfigInt64(confLabel, confDesc, value, nil)
	if len(builder.config) != 1 {
		t.Error("Invalid config length.")
	}
}
func TestFormScreenBuilderAddConfigFloat(t *testing.T) {
	builder := NewFormScreenBuilder()
	if len(builder.config) != 0 {
		t.Error("Invalid initial config length.")
	}
	confLabel := "label"
	confDesc := "desc"
	value := float32(3.3)
	builder.AddConfigFloat(confLabel, confDesc, value, nil)
	if len(builder.config) != 1 {
		t.Error("Invalid config length.")
	}
}
func TestFormScreenBuilderAddConfigText(t *testing.T) {
	builder := NewFormScreenBuilder()
	if len(builder.config) != 0 {
		t.Error("Invalid initial config length.")
	}
	confLabel := "label"
	confDesc := "desc"
	value := "test"
	builder.AddConfigText(confLabel, confDesc, value, nil)
	if len(builder.config) != 1 {
		t.Error("Invalid config length.")
	}
}
func TestFormScreenBuilderAddConfigVector(t *testing.T) {
	builder := NewFormScreenBuilder()
	if len(builder.config) != 0 {
		t.Error("Invalid initial config length.")
	}
	confLabel := "label"
	confDesc := "desc"
	value := mgl32.Vec3{0, 1, 0}
	builder.AddConfigVector(confLabel, confDesc, value, nil)
	if len(builder.config) != 1 {
		t.Error("Invalid config length.")
	}
}
func TestFormScreenBuilderSetFrameSize(t *testing.T) {
	builder := NewFormScreenBuilder()
	w := float32(2)
	l := float32(3)
	r := float32(0.1)
	builder.SetFrameSize(w, l, r)
	if builder.frameWidth != w {
		t.Errorf("Invalid frame width. Instead of '%f', it is '%f'.", w, builder.frameWidth)
	}
	if builder.frameLength != l {
		t.Errorf("Invalid frame length. Instead of '%f', it is '%f'.", l, builder.frameLength)
	}
	if builder.frameTopLeftWidth != r {
		t.Errorf("Invalid frame top left length. Instead of '%f', it is '%f'.", r, builder.frameTopLeftWidth)
	}
}
func TestFormScreenBuilderSetFrameMaterial(t *testing.T) {
	builder := NewFormScreenBuilder()
	mat := material.Ruby
	builder.SetFrameMaterial(mat)
	if builder.frameMaterial != mat {
		t.Error("Invalid frame material")
	}
}
func TestFormScreenBuilderSetFormItemMaterial(t *testing.T) {
	builder := NewFormScreenBuilder()
	mat := material.Ruby
	builder.SetFormItemMaterial(mat)
	if builder.formItemMaterial != mat {
		t.Error("Invalid formItemMaterial.")
	}
}
func TestFormScreenBuilderSetFormItemHighlightMaterial(t *testing.T) {
	builder := NewFormScreenBuilder()
	mat := material.Ruby
	builder.SetFormItemHighlightMaterial(mat)
	if builder.formItemHighlightMaterial != mat {
		t.Error("Invalid formItemHighlightMaterial.")
	}
}
func TestFormScreenBuilderSetDetailContentBoxHeight(t *testing.T) {
	builder := NewFormScreenBuilder()
	height := float32(0.2)
	builder.SetDetailContentBoxHeight(height)
	if builder.detailContentBoxHeight != height {
		t.Errorf("Invalid detailContentBoxHeight. Instead of '%f', it is '%f'.", height, builder.detailContentBoxHeight)
	}
}
func TestFormScreenBuilderBuild(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Errorf("Shouldn't have panic, %#v.", r)
			}
		}()
		runtime.LockOSThread()
		testhelper.GlfwInit()
		defer testhelper.GlfwTerminate()
		wrapperReal.InitOpenGL()

		builder := NewFormScreenBuilder()
		wW := float32(800)
		wH := float32(800)
		builder.SetWindowSize(wW, wH)
		w := float32(2)
		l := float32(3)
		r := float32(0.1)
		builder.SetFrameSize(w, l, r)
		builder.SetWrapper(wrapperReal)
		confLabel := "label"
		confDesc := "desc"

		boolVal := true
		intVal := 3
		int64Val := int64(3)
		floatVal := float32(3.1)
		stringVal := "nine"
		vecVal := mgl32.Vec3{1, 0, 0}

		configOrder := []string{
			// bool line
			builder.AddConfigBool(confLabel, confDesc, boolVal),
			builder.AddConfigBool(confLabel, confDesc, boolVal),
			builder.AddConfigBool(confLabel, confDesc, boolVal),
			// int line
			builder.AddConfigInt(confLabel, confDesc, intVal, nil),
			builder.AddConfigInt(confLabel, confDesc, intVal, nil),
			// int64 val
			builder.AddConfigInt64(confLabel, confDesc, int64Val, nil),
			builder.AddConfigBool(confLabel, confDesc, boolVal),
			builder.AddConfigInt64(confLabel, confDesc, int64Val, nil),
			// float val
			builder.AddConfigFloat(confLabel, confDesc, floatVal, nil),
			builder.AddConfigFloat(confLabel, confDesc, floatVal, nil),
			// text val
			builder.AddConfigText(confLabel, confDesc, stringVal, nil),
			// vector val
			builder.AddConfigVector(confLabel, confDesc, vecVal, nil),
			// LL -> F
			builder.AddConfigInt64(confLabel, confDesc, int64Val, nil),
			builder.AddConfigText(confLabel, confDesc, stringVal, nil),
			// LS -> F
			builder.AddConfigBool(confLabel, confDesc, boolVal),
			builder.AddConfigVector(confLabel, confDesc, vecVal, nil),
			// LS -> RL
			builder.AddConfigBool(confLabel, confDesc, boolVal),
			builder.AddConfigInt64(confLabel, confDesc, int64Val, nil),
			// MS -> F
			builder.AddConfigBool(confLabel, confDesc, boolVal),
			builder.AddConfigBool(confLabel, confDesc, boolVal),
			builder.AddConfigText(confLabel, confDesc, stringVal, nil),
			// MS -> LL
			builder.AddConfigBool(confLabel, confDesc, boolVal),
			builder.AddConfigBool(confLabel, confDesc, boolVal),
			builder.AddConfigInt64(confLabel, confDesc, int64Val, nil),
			// MS -> LH
			builder.AddConfigVector(confLabel, confDesc, vecVal, nil),
			builder.AddConfigBool(confLabel, confDesc, boolVal),
			builder.AddConfigBool(confLabel, confDesc, boolVal),
			builder.AddConfigInt(confLabel, confDesc, intVal, nil),
			// LL -> LL
			builder.AddConfigText(confLabel, confDesc, stringVal, nil),
			builder.AddConfigInt64(confLabel, confDesc, int64Val, nil),
			builder.AddConfigInt64(confLabel, confDesc, int64Val, nil),
			// On bool val
			builder.AddConfigBool(confLabel, confDesc, true),

			builder.AddConfigInt(confLabel, confDesc, intVal, nil),
			builder.AddConfigInt(confLabel, confDesc, intVal, nil),
			builder.AddConfigInt(confLabel, confDesc, intVal, nil),

			builder.AddConfigFloat(confLabel, confDesc, floatVal, nil),
			builder.AddConfigFloat(confLabel, confDesc, floatVal, nil),
			builder.AddConfigFloat(confLabel, confDesc, floatVal, nil),
			builder.AddConfigFloat(confLabel, confDesc, floatVal, nil),

			builder.AddConfigText(confLabel, confDesc, stringVal, nil),
			builder.AddConfigText(confLabel, confDesc, stringVal, nil),

			builder.AddConfigVector(confLabel, confDesc, vecVal, nil),
			builder.AddConfigVector(confLabel, confDesc, vecVal, nil),
			builder.AddConfigVector(confLabel, confDesc, vecVal, nil),
		}
		builder.SetConfigOrder(configOrder)
		_ = builder.Build()
	}()
}
func TestFormScreenUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Errorf("Shouldn't have panic, %#v.", r)
			}
		}()
		wW := float32(800)
		wH := float32(800)
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		defer testhelper.GlfwTerminate()
		builder := NewFormScreenBuilder()
		builder.SetWrapper(wrapperReal)
		builder.SetWindowSize(wW, wH)
		builder.SetDetailContentBoxHeight(0.3)
		longDescription := "This is a long text to check that the wrap text function works well or not. This is a long text to check that the wrap text function works well or not."
		builder.SetConfigOrder([]string{
			builder.AddConfigBool("label bool", DefaultFormItemDescription, true),
			builder.AddConfigInt("label int", DefaultFormItemDescription, 1, nil),
			builder.AddConfigInt64("label int64", DefaultFormItemDescription, 10, nil),
			builder.AddConfigFloat("label float", longDescription, 0.44, nil),
			builder.AddConfigText("label text", DefaultFormItemDescription, "sample", nil),
			builder.AddConfigVector("label vector", DefaultFormItemDescription, mgl32.Vec3{0.01, 0.02, 0.03}, nil),
			builder.AddConfigBool("label bool", DefaultFormItemDescription, true),
			builder.AddConfigInt("label int", DefaultFormItemDescription, 1, nil),
			builder.AddConfigInt64("label int64", DefaultFormItemDescription, 10, nil),
			builder.AddConfigFloat("label float", longDescription, 0.44, nil),
			builder.AddConfigText("label text", DefaultFormItemDescription, "sample", nil),
			builder.AddConfigVector("label vector", DefaultFormItemDescription, mgl32.Vec3{0.01, 0.02, 0.03}, nil),
			builder.AddConfigBool("label bool", DefaultFormItemDescription, true),
			builder.AddConfigInt("label int", DefaultFormItemDescription, 1, nil),
			builder.AddConfigInt64("label int64", DefaultFormItemDescription, 10, nil),
			builder.AddConfigFloat("label float", longDescription, 0.44, nil),
			builder.AddConfigText("label text", DefaultFormItemDescription, "sample", nil),
			builder.AddConfigVector("label vector", DefaultFormItemDescription, mgl32.Vec3{0.01, 0.02, 0.03}, nil),
			builder.AddConfigBool("label bool", DefaultFormItemDescription, true),
			builder.AddConfigInt("label int", DefaultFormItemDescription, 1, nil),
			builder.AddConfigInt64("label int64", DefaultFormItemDescription, 10, nil),
			builder.AddConfigFloat("label float", longDescription, 0.44, nil),
			builder.AddConfigText("label text", DefaultFormItemDescription, "sample", nil),
			builder.AddConfigVector("label vector", DefaultFormItemDescription, mgl32.Vec3{0.01, 0.02, 0.03}, nil),
			builder.AddConfigBool("label bool", DefaultFormItemDescription, true),
			builder.AddConfigInt("label int", DefaultFormItemDescription, 1, nil),
			builder.AddConfigInt64("label int64", DefaultFormItemDescription, 10, nil),
			builder.AddConfigFloat("label float", longDescription, 0.44, nil),
			builder.AddConfigText("label text", DefaultFormItemDescription, "sample", nil),
			builder.AddConfigVector("label vector", DefaultFormItemDescription, mgl32.Vec3{0.01, 0.02, 0.03}, nil),
			builder.AddConfigBool("label bool", DefaultFormItemDescription, true),
			builder.AddConfigInt("label int", DefaultFormItemDescription, 1, nil),
			builder.AddConfigInt64("label int64", DefaultFormItemDescription, 10, nil),
			builder.AddConfigFloat("label float", longDescription, 0.44, nil),
			builder.AddConfigText("label text", DefaultFormItemDescription, "sample", nil),
			builder.AddConfigVector("label vector", DefaultFormItemDescription, mgl32.Vec3{0.01, 0.02, 0.03}, nil),
		})
		form := builder.Build()
		ks := store.NewGlfwKeyStore()
		ms := store.NewGlfwMouseStore()
		form.Update(10, 0.5, 0.5, ks, ms)
		form.Update(10, -0.4, 0.79, ks, ms)
		ms.Set(LEFT_MOUSE_BUTTON, true)
		form.sinceLastClick = 201
		form.Update(10, -0.4, 0.79, ks, ms) // bool
		form.sinceLastClick = 201
		form.Update(10, 0.4, 0.79, ks, ms) // int
		form.sinceLastClick = 201
		form.sinceLastDelete = 201
		ks.Set(BACK_SPACE, true)
		form.Update(10, 0.4, 0.79, ks, ms) // int
		form.sinceLastClick = 201
		form.sinceLastDelete = 201
		form.Update(10, 0.4, 0.69, ks, ms) // int64
		form.sinceLastClick = 201
		form.sinceLastDelete = 201
		form.Update(10, -0.4, 0.69, ks, ms) // int64
		form.sinceLastClick = 201
		form.sinceLastDelete = 201
		form.Update(10, -0.4, 0.59, ks, ms) // float
		form.sinceLastClick = 201
		form.sinceLastDelete = 201
		form.Update(10, -0.4, 0.69, ks, ms) // int64
		form.sinceLastClick = 201
		form.sinceLastDelete = 201
		form.Update(10, -0.4, 0.59, ks, ms) // float
		form.sinceLastClick = 201
		form.sinceLastDelete = 201
		form.Update(10, -0.4, 0.49, ks, ms) // text
		form.sinceLastClick = 201
		form.sinceLastDelete = 201
		form.Update(10, -0.4, 0.39, ks, ms) // vector
		ms.Set(LEFT_MOUSE_BUTTON, false)
		ks.Set(KEY_UP, true)
		form.Update(0.4, -0.4, 0.79, ks, ms)
		ks.Set(KEY_UP, false)
		ks.Set(KEY_DOWN, true)
		form.Update(0.4, -0.4, 0.79, ks, ms)
		ks.Set(KEY_DOWN, false)
		form.Update(0.4, -0.4, 0.79, ks, ms)
	}()
}
func newFormScreen() *FormScreen {
	builder := NewFormScreenBuilder()
	builder.SetWrapper(wrapperReal)
	builder.SetWindowSize(800, 800)
	return builder.Build()
}
func TestFormScreenSetupFormScreen(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Errorf("Shouldn't have panic, %#v.", r)
			}
		}()
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		form := newFormScreen()
		form.setupFormScreen(wrapperReal)
	}()
}
func TestNewFormScreen(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Errorf("Shouldn't have panic, %#v.", r)
			}
		}()
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		form := newFormScreen()
		defer testhelper.GlfwTerminate()
		if len(form.formItemToConf) != 0 {
			t.Errorf("Invalid initial formItemToConf length. '%d'.", len(form.formItemToConf))
		}
	}()
}
func TestFormGetFormItemValidIndex(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Errorf("Shouldn't have panic, %#v.", r)
			}
		}()
		wW := float32(800)
		wH := float32(800)
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		defer testhelper.GlfwTerminate()
		builder := NewFormScreenBuilder()
		builder.SetWrapper(wrapperReal)
		builder.SetWindowSize(wW, wH)
		textKey := builder.AddConfigText("text", DefaultFormItemDescription, "", nil)
		intKey := builder.AddConfigInt("int", DefaultFormItemDescription, 0, nil)
		int64Key := builder.AddConfigInt64("int64", DefaultFormItemDescription, 3, nil)
		floatKey := builder.AddConfigFloat("float", DefaultFormItemDescription, 0.0, nil)
		boolKey := builder.AddConfigBool("bool", DefaultFormItemDescription, false)
		vectorKey := builder.AddConfigVector("vector", DefaultFormItemDescription, mgl32.Vec3{0, 1, 0}, nil)
		builder.SetConfigOrder([]string{textKey, intKey, int64Key, floatKey, boolKey, vectorKey})
		form := builder.Build()
		fi := form.GetFormItem(textKey)
		if fi.ValueToString() != "" {
			t.Error("Invalid form item initial value")
		}
		fi = form.GetFormItem(intKey)
		if fi.ValueToString() != "0" {
			t.Error("Invalid form item initial value")
		}
		fi = form.GetFormItem(int64Key)
		if fi.ValueToString() != "3" {
			t.Error("Invalid form item initial value")
		}
		fi = form.GetFormItem(floatKey)
		if fi.ValueToString() != "0" {
			t.Errorf("Invalid form item initial value. '%s'.", fi.ValueToString())
		}
		fi = form.GetFormItem(boolKey)
		if fi.ValueToString() != "false" {
			t.Error("Invalid form item initial value")
		}
	}()
}
func TestFormSetFormItemValue(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Errorf("Shouldn't have panic, %#v.", r)
			}
		}()
		wW := float32(800)
		wH := float32(800)
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		defer testhelper.GlfwTerminate()
		builder := NewFormScreenBuilder()
		builder.SetWrapper(wrapperReal)
		builder.SetWindowSize(wW, wH)
		textKey := builder.AddConfigText("text", DefaultFormItemDescription, "text value", nil)
		intKey := builder.AddConfigInt("int", DefaultFormItemDescription, 0, nil)
		int64Key := builder.AddConfigInt64("int64", DefaultFormItemDescription, 3, nil)
		floatKey := builder.AddConfigFloat("float", DefaultFormItemDescription, 0.0, nil)
		boolKey := builder.AddConfigBool("bool", DefaultFormItemDescription, false)
		vectorKey := builder.AddConfigVector("vector", DefaultFormItemDescription, mgl32.Vec3{0, 1, 0}, nil)
		builder.SetConfigOrder([]string{textKey, intKey, int64Key, floatKey, boolKey, vectorKey})
		form := builder.Build()
		newTextValue := "new text"
		fi := form.GetFormItem(textKey)
		form.SetFormItemValue(fi, newTextValue)
		if fi.ValueToString() != newTextValue {
			t.Errorf("Invalid form item value. Instead of '%s', it is '%s'.", newTextValue, fi.ValueToString())
		}
		fi = form.GetFormItem(intKey)
		newTextValue = "2"
		form.SetFormItemValue(fi, newTextValue)
		if fi.ValueToString() != newTextValue {
			t.Errorf("Invalid form item value. Instead of '%s', it is '%s'.", newTextValue, fi.ValueToString())
		}
		fi = form.GetFormItem(int64Key)
		newTextValue = "200"
		form.SetFormItemValue(fi, newTextValue)
		if fi.ValueToString() != newTextValue {
			t.Errorf("Invalid form item value. Instead of '%s', it is '%s'.", newTextValue, fi.ValueToString())
		}
		fi = form.GetFormItem(floatKey)
		newTextValue = "20.002"
		form.SetFormItemValue(fi, newTextValue)
		if fi.ValueToString() != newTextValue {
			t.Errorf("Invalid form item value. Instead of '%s', it is '%s'.", newTextValue, fi.ValueToString())
		}
		fi = form.GetFormItem(boolKey)
		form.SetFormItemValue(fi, true)
		if fi.ValueToString() != "true" {
			t.Errorf("Invalid form item value.")
		}
		fi = form.GetFormItem(vectorKey)
		form.SetFormItemValue(fi, [3]string{"0.1", "1.0", "0.5"})
	}()
}
func TestFormGetFormItemValidIndexValidators(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Errorf("Shouldn't have panic, %#v.", r)
			}
		}()
		wW := float32(800)
		wH := float32(800)
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		defer testhelper.GlfwTerminate()
		builder := NewFormScreenBuilder()
		builder.SetWrapper(wrapperReal)
		builder.SetWindowSize(wW, wH)
		textKey := builder.AddConfigText("text", DefaultFormItemDescription, "", func(t string) bool { return true })
		intKey := builder.AddConfigInt("int", DefaultFormItemDescription, 0, func(i int) bool { return true })
		int64Key := builder.AddConfigInt64("int64", DefaultFormItemDescription, 3, func(i int64) bool { return true })
		floatKey := builder.AddConfigFloat("float", DefaultFormItemDescription, 0.0, func(f float32) bool { return true })
		boolKey := builder.AddConfigBool("bool", DefaultFormItemDescription, false)
		vectorKey := builder.AddConfigVector("vector", DefaultFormItemDescription, mgl32.Vec3{0, 1, 0}, func(f float32) bool { return true })
		builder.SetConfigOrder([]string{textKey, intKey, int64Key, floatKey, boolKey, vectorKey})
		form := builder.Build()
		fi := form.GetFormItem(textKey)
		if fi.ValueToString() != "" {
			t.Error("Invalid form item initial value")
		}
		fi = form.GetFormItem(intKey)
		if fi.ValueToString() != "0" {
			t.Error("Invalid form item initial value")
		}
		fi = form.GetFormItem(int64Key)
		if fi.ValueToString() != "3" {
			t.Error("Invalid form item initial value")
		}
		fi = form.GetFormItem(floatKey)
		if fi.ValueToString() != "0" {
			t.Errorf("Invalid form item initial value. '%s'.", fi.ValueToString())
		}
		fi = form.GetFormItem(boolKey)
		if fi.ValueToString() != "false" {
			t.Error("Invalid form item initial value")
		}
		_ = form.GetFormItem(vectorKey)
	}()
}
func TestFormGetFormItemInvalidIndex(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r == nil {
				defer testhelper.GlfwTerminate()
				t.Error("Should have panic.")
			}
		}()
		wW := float32(800)
		wH := float32(800)
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		builder := NewFormScreenBuilder()
		builder.SetWrapper(wrapperReal)
		builder.SetWindowSize(wW, wH)
		form := builder.Build()
		defer testhelper.GlfwTerminate()
		_ = form.GetFormItem("invalidkey")
	}()
}
func TestFormScreenCharCallback(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Errorf("Shouldn't have panic, %#v.", r)
			}
		}()
		wW := float32(800)
		wH := float32(800)
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		builder := NewFormScreenBuilder()
		builder.SetWrapper(wrapperReal)
		builder.SetWindowSize(wW, wH)
		textKey := builder.AddConfigText("text", DefaultFormItemDescription, "", nil)
		intKey := builder.AddConfigInt("int", DefaultFormItemDescription, 1, nil)
		int64Key := builder.AddConfigInt64("int64", DefaultFormItemDescription, 2, nil)
		floatKey := builder.AddConfigFloat("float", DefaultFormItemDescription, 0.0, nil)
		builder.SetConfigOrder([]string{textKey, intKey, int64Key, floatKey})
		form := builder.Build()
		defer testhelper.GlfwTerminate()
		form.underEdit = form.GetFormItem(textKey).(*model.FormItemText)
		form.CharCallback('1', wrapperReal)
		if form.underEdit.ValueToString() != "1" {
			t.Errorf("Invalid value: '%s'.", form.underEdit.ValueToString())
		}
		form.underEdit = form.GetFormItem(intKey).(*model.FormItemInt)
		form.CharCallback('1', wrapperReal)
		if form.underEdit.ValueToString() != "11" {
			t.Errorf("Invalid value: '%s'.", form.underEdit.ValueToString())
		}
		form.underEdit = form.GetFormItem(int64Key).(*model.FormItemInt64)
		form.CharCallback('1', wrapperReal)
		if form.underEdit.ValueToString() != "21" {
			t.Errorf("Invalid value: '%s'.", form.underEdit.ValueToString())
		}
		form.underEdit = form.GetFormItem(floatKey).(*model.FormItemFloat)
		form.CharCallback('1', wrapperReal)
		if form.underEdit.ValueToString() != "0" {
			t.Errorf("Invalid value: '%s'.", form.underEdit.ValueToString())
		}
		form.CharCallback('.', wrapperReal)
		if form.underEdit.ValueToString() != "0." {
			t.Errorf("Invalid value: '%s'.", form.underEdit.ValueToString())
		}
		form.CharCallback('1', wrapperReal)
		if form.underEdit.ValueToString() != "0.1" {
			t.Errorf("Invalid value: '%s'.", form.underEdit.ValueToString())
		}
	}()
}
