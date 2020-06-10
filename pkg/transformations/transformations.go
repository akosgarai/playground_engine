package transformations

import (
	"strconv"

	"github.com/go-gl/mathgl/mgl32"
)

// Returns the mouse coordinate in window coordinate.
func MouseCoordinates(currentX, currentY, windowWidth, windowHeight float64) (float64, float64) {
	halfWidth := windowWidth / 2.0
	halfHeight := windowHeight / 2.0
	x := (currentX - halfWidth) / (halfWidth)
	y := (halfHeight - currentY) / (halfHeight)
	return x, y
}

// Vec3ToString helper function for the string representation of a vector. It is for the log.
func Vec3ToString(v mgl32.Vec3) string {
	x := Float32ToString(v.X())
	y := Float32ToString(v.Y())
	z := Float32ToString(v.Z())
	return "X : " + x + ", Y : " + y + ", Z : " + z
}

// Float64ToString returns the given float number in string format.
func Float64ToString(num float64) string {
	return strconv.FormatFloat(num, 'f', 10, 32)
}

// Float32ToString returns the given float number in string format.
func Float32ToString(num float32) string {
	return Float64ToString(float64(num))
}

// IntegerToString returns the string representation of the given integer
func IntegerToString(num int) string {
	return strconv.Itoa(num)
}

// This function returns the abs. value of a float32 number.
// If the number is less than 0, it returns -1*number, otherwise
// it returns the number itself.
func Float32Abs(a float32) float32 {
	if a < 0 {
		return -a
	}
	return a
}
