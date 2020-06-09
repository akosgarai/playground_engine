package boundingobject

type BoundingObject struct {
	typeName string
	params   map[string]float32
}

// New returns a BoundingObject. The name input supposed to be 'AABB' or 'Sphere'
// In case of 'AABB', the following keys needs to be set as params: 'width', 'length',
// 'height'. In case of 'Sphere' name, only the 'radius' key needs to be set.
func New(name string, params map[string]float32) *BoundingObject {
	return &BoundingObject{
		typeName: name,
		params:   params,
	}
}

// Type return the type name fo the BO
func (bo *BoundingObject) Type() string {
	return bo.typeName
}

// Type return the type name fo the BO
func (bo *BoundingObject) Params() map[string]float32 {
	return bo.params
}
