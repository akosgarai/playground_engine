package screen

import (
	"runtime"
	"testing"

	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/testhelper"
)

func TestNewScreenWithFrameBuilder(t *testing.T) {
	fb := NewScreenWithFrameBuilder()
	dfw := DefaultFrameWidth
	if fb.frameWidth != dfw {
		t.Errorf("Invalid frameWidth. Instead of '%f', it is '%f'.", dfw, fb.frameWidth)
	}
	dfl := DefaultFrameLength
	if fb.frameLength != dfl {
		t.Errorf("Invalid frameLength. Instead of '%f', it is '%f'.", dfl, fb.frameLength)
	}
	dtlw := TopLeftFrameWidth
	if fb.frameTopLeftWidth != dtlw {
		t.Errorf("Invalid frameTopLeftWidth. Instead of '%f', it is '%f'.", dtlw, fb.frameTopLeftWidth)
	}
	fov := float32(45)
	if fb.fov != fov {
		t.Errorf("Invalid fov. Instead of '%f', it is '%f'.", fov, fb.fov)
	}
	labelWidth := float32(0)
	if fb.labelWidth != labelWidth {
		t.Errorf("Invalid labelWidth. Instead of '%f', it is '%f'.", labelWidth, fb.labelWidth)
	}
	detailContentBoxHeight := float32(0)
	if fb.detailContentBoxHeight != detailContentBoxHeight {
		t.Errorf("Invalid labelWidth. Instead of '%f', it is '%f'.", detailContentBoxHeight, fb.detailContentBoxHeight)
	}
	detailContentBoxMaterial := DefaultFormItemMaterial
	if fb.detailContentBoxMaterial != detailContentBoxMaterial {
		t.Error("Invalid detailContentBoxMaterial.")
	}
}
func TestScreenWithFrameBuilderSetWindowSize(t *testing.T) {
	fb := NewScreenWithFrameBuilder()
	wW := float32(100)
	wH := float32(100)
	fb.SetWindowSize(wW, wH)
	if fb.windowWidth != wW {
		t.Errorf("Invalid window width. Instead of '%f', it is '%f'.", wW, fb.windowWidth)
	}
	if fb.windowHeight != wH {
		t.Errorf("Invalid window height. Instead of '%f', it is '%f'.", wH, fb.windowHeight)
	}
}
func TestScreenWithFrameBuilderSetFrameSize(t *testing.T) {
	fb := NewScreenWithFrameBuilder()
	w := float32(2)
	l := float32(3)
	r := float32(0.1)
	fb.SetFrameSize(w, l, r)
	if fb.frameWidth != w {
		t.Errorf("Invalid frame width. Instead of '%f', it is '%f'.", w, fb.frameWidth)
	}
	if fb.frameLength != l {
		t.Errorf("Invalid frame length. Instead of '%f', it is '%f'.", l, fb.frameLength)
	}
	if fb.frameTopLeftWidth != r {
		t.Errorf("Invalid frame top left length. Instead of '%f', it is '%f'.", r, fb.frameTopLeftWidth)
	}
}
func TestScreenWithFrameBuilderSetWrapper(t *testing.T) {
	fb := NewScreenWithFrameBuilder()
	fb.SetWrapper(wrapperMock)
	if fb.wrapper != wrapperMock {
		t.Error("Invalid wrapper")
	}
}
func TestScreenWithFrameBuilderSetFrameMaterial(t *testing.T) {
	fb := NewScreenWithFrameBuilder()
	mat := material.Ruby
	fb.SetFrameMaterial(mat)
	if fb.frameMaterial != mat {
		t.Error("Invalid frame material")
	}
}
func TestScreenWithFrameBuilderSetDetailContentBoxMaterial(t *testing.T) {
	fb := NewScreenWithFrameBuilder()
	mat := material.Ruby
	fb.SetDetailContentBoxMaterial(mat)
	if fb.detailContentBoxMaterial != mat {
		t.Error("Invalid detailContentBox material")
	}
}
func TestScreenWithFrameBuilderSetLabelWidth(t *testing.T) {
	fb := NewScreenWithFrameBuilder()
	width := float32(0.5)
	fb.SetLabelWidth(width)
	if fb.labelWidth != width {
		t.Errorf("invalid label width. Instead of '%f', it is '%f'.", width, fb.labelWidth)
	}
}
func TestScreenWithFrameBuilderSetDetailContentBoxHeight(t *testing.T) {
	fb := NewScreenWithFrameBuilder()
	height := float32(0.2)
	fb.SetDetailContentBoxHeight(height)
	if fb.detailContentBoxHeight != height {
		t.Errorf("Invalid detailContentBoxHeight. Instead of '%f', it is '%f'.", height, fb.detailContentBoxHeight)
	}
}
func TestScreenWithFrameBuilderBuild(t *testing.T) {
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
		testhelper.GlfwInit(glwrapper.GL_MAJOR_VERSION, glwrapper.GL_MINOR_VERSION)
		defer testhelper.GlfwTerminate()
		wrapperReal.InitOpenGL()
		fb := NewScreenWithFrameBuilder()
		fb.SetWrapper(wrapperReal)
		wW := float32(100)
		wH := float32(100)
		fb.SetWindowSize(wW, wH)
		_ = fb.Build()
	}()
}
func TestScreenWithFrameBuilderBuildWithoutWrapper(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r == nil {
				defer testhelper.GlfwTerminate()
				t.Errorf("Should have panic.")
			}
		}()
		runtime.LockOSThread()
		testhelper.GlfwInit(glwrapper.GL_MAJOR_VERSION, glwrapper.GL_MINOR_VERSION)
		defer testhelper.GlfwTerminate()
		wrapperReal.InitOpenGL()
		fb := NewScreenWithFrameBuilder()
		wW := float32(100)
		wH := float32(100)
		fb.SetWindowSize(wW, wH)
		_ = fb.Build()
	}()
}
func TestScreenWithFrame(t *testing.T) {
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
		testhelper.GlfwInit(glwrapper.GL_MAJOR_VERSION, glwrapper.GL_MINOR_VERSION)
		defer testhelper.GlfwTerminate()
		wrapperReal.InitOpenGL()
		fb := NewScreenWithFrameBuilder()
		fb.SetWrapper(wrapperReal)
		wW := float32(100)
		wH := float32(100)
		fb.SetWindowSize(wW, wH)
		w := float32(2)
		l := float32(0.5)
		r := float32(0.1)
		fb.SetFrameSize(w, l, r)
		s := fb.Build()
		if s.frameWidth != w {
			t.Errorf("Invalid frame width. Instead of '%f', it is '%f'.", w, s.frameWidth)
		}
		if s.frameLength != l {
			t.Errorf("Invalid frame length. Instead of '%f', it is '%f'.", l, s.frameLength)
		}
		if s.GetFullWidth() != float32(1.0) {
			t.Errorf("Invalid full width. Instead of '1.0', it is '%f'.", s.GetFullWidth())
		}
	}()
}
