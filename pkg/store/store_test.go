package store

import (
	"testing"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func TestNewGlfwKeyStore(t *testing.T) {
	store := NewGlfwKeyStore()
	if len(store) != 0 {
		t.Errorf("Invalid store length. It is %d.", len(store))
	}
}
func TestNewGlfwMouseStore(t *testing.T) {
	store := NewGlfwMouseStore()
	if len(store) != 0 {
		t.Errorf("Invalid store length. It is %d.", len(store))
	}
}
func TestKeyStoreSet(t *testing.T) {
	store := NewGlfwKeyStore()
	btn := glfw.KeyA
	value := true
	store.Set(btn, value)
	if store[btn] != value {
		t.Errorf("Invalid value set. It is %v.", store[btn])
	}
}
func TestKeyStoreGet(t *testing.T) {
	store := NewGlfwKeyStore()
	btn := glfw.KeyA
	value := true
	store.Set(btn, value)
	if store.Get(btn) != value {
		t.Errorf("Invalid value get. It is %v.", store.Get(btn))
	}
}
func TestMouseStoreSet(t *testing.T) {
	store := NewGlfwMouseStore()
	btn := glfw.MouseButtonLeft
	value := true
	store.Set(btn, value)
	if store[btn] != value {
		t.Errorf("Invalid value set. It is %v.", store[btn])
	}
}
func TestMouseStoreGet(t *testing.T) {
	store := NewGlfwMouseStore()
	btn := glfw.MouseButtonLeft
	value := true
	store.Set(btn, value)
	if store.Get(btn) != value {
		t.Errorf("Invalid value get. It is %v.", store.Get(btn))
	}
}
