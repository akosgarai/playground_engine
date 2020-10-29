package ui

import (
	"testing"
)

func TestNew(t *testing.T) {
	model := New()
	if _, err := model.GetMeshByIndex(0); err.Error() != "EMPTY_MESHES" {
		t.Error("Invalid empty meshes error.")
	}
}
