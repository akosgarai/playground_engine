package testhelper

import ()

// This function returns true, if the given a, b is almost equal,
// the difference between them is less than epsilon.
func Float32ApproxEqual(a, b, epsilon float32) bool {
	return (a-b) < epsilon && (b-a) < epsilon
}
