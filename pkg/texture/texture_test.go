package texture

import (
	"testing"

	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/testhelper"
)

var testGlWrapper testhelper.GLWrapperMock

func TestLoadImageFromFile(t *testing.T) {
	_, err := loadImageFromFile("this-image-does-not-exist.jpg")
	if err == nil {
		t.Error("Image load should be failed.")
	}
	_, err = loadImageFromFile("assets/testing.jpg")
	if err != nil {
		t.Error("Issue during load.")
	}
}
func TestAddTexture(t *testing.T) {
	var textures Textures
	textures.AddTexture("assets/testing.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", testGlWrapper)
	if len(textures) != 1 {
		t.Error("AddTexture should be successful")
	}
}
func TestAddTextureInvalidName(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("AddTexture should have panicked due to the missing file!")
			}
		}()
		var textures Textures
		textures.AddTexture("this-image-does-not-exist.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", testGlWrapper)
	}()
}
func TestGenTextures(t *testing.T) {
	genTextures(testGlWrapper)
}
func TestBindTexture(t *testing.T) {
	var textures Textures
	textures.AddTexture("assets/testing.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", testGlWrapper)
	textures[0].Bind()
}
func TestUnBindTexture(t *testing.T) {
	var textures Textures
	textures.AddTexture("assets/testing.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", testGlWrapper)
	textures[0].UnBind()
}
func TestUnBindTextures(t *testing.T) {
	var textures Textures
	textures.AddTexture("assets/testing.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", testGlWrapper)
	textures.UnBind()
}

func TestAddCubeMapTexture(t *testing.T) {
	var textures Textures
	textures.AddCubeMapTexture("assets", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", testGlWrapper)
	if len(textures) != 1 {
		t.Error("AddTexture should be successful")
	}
}
func TestAddCubeMapTextureWrongDir(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("It should be panicked due to the wrong directory.")
			}
		}()
		var textures Textures
		textures.AddCubeMapTexture("wrongDir", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", testGlWrapper)
	}()
}
