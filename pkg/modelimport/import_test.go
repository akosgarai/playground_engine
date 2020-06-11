package modelimport

import (
	"testing"

	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/testhelper"
)

const (
	ObjectFileName = "object.obj"
	Directory      = "testdata"
)

var (
	wrapperMock   testhelper.GLWrapperMock
	realGlWrapper glwrapper.Wrapper
)

func TestNew(t *testing.T) {
	importer := New(Directory, ObjectFileName, wrapperMock)

	if importer.objectFile != ObjectFileName {
		t.Errorf("Invalid object file name. Instead of '%s', we have '%s'.", ObjectFileName, importer.objectFile)
	}
	if importer.basePath != Directory {
		t.Errorf("Invalid directory name. Instead of '%s', we have '%s'.", Directory, importer.basePath)
	}
	if len(importer.meshes) != 0 {
		t.Errorf("Invalid initial mesh length. Instead of '0', we have '%d'.", len(importer.meshes))
	}
}
func TestGetMeshes(t *testing.T) {
	importer := New(Directory, ObjectFileName, wrapperMock)
	importer.meshes = append(importer.meshes, mesh.NewPointMesh(wrapperMock))
	importer.meshes = append(importer.meshes, mesh.NewPointMesh(wrapperMock))
	if len(importer.GetMeshes()) != 2 {
		t.Errorf("Invalid mesh length. Instead of '2', we have '%d'.", len(importer.GetMeshes()))
	}
}
func TestImport(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Should be fine")
				t.Log(r)
			}
		}()
		testhelper.GlfwInit()
		defer testhelper.GlfwTerminate()
		realGlWrapper.InitOpenGL()
		importer := New(Directory, ObjectFileName, wrapperMock)
		importer.Import()
		meshes := importer.GetMeshes()
		t.Log(meshes)
	}()
}
