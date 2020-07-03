package store

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type GlfwKeyStore map[glfw.Key]bool

type GlfwMouseStore map[glfw.MouseButton]bool

// NewGlfwKeyStore returns an empty keyboard button state store.
func NewGlfwKeyStore() GlfwKeyStore {
	return make(map[glfw.Key]bool)
}

// Set sets the given key to the given value.
func (store GlfwKeyStore) Set(key glfw.Key, value bool) {
	store[key] = value
}

// Get return the value for the given key
func (store GlfwKeyStore) Get(key glfw.Key) bool {
	return store[key]
}

// NewGlfwMouseStore returns an empty mouse button state store.
func NewGlfwMouseStore() GlfwMouseStore {
	return make(map[glfw.MouseButton]bool)
}

// Set sets the given key to the given value.
func (store GlfwMouseStore) Set(key glfw.MouseButton, value bool) {
	store[key] = value
}

// Get return the value for the given key
func (store GlfwMouseStore) Get(key glfw.MouseButton) bool {
	return store[key]
}
