package material

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/akosgarai/playground_engine/pkg/transformations"
)

type Material struct {
	ambient   mgl32.Vec3
	diffuse   mgl32.Vec3
	specular  mgl32.Vec3
	shininess float32
}

func New(ambient, diffuse, specular mgl32.Vec3, shininess float32) *Material {
	return &Material{
		ambient:   ambient,
		diffuse:   diffuse,
		specular:  specular,
		shininess: shininess,
	}
}

func (m *Material) Log() string {
	logString := "Material\n"
	logString += " - Ambient: Vector{" + transformations.Vec3ToString(m.ambient) + "}\n"
	logString += " - Diffuse: Vector{" + transformations.Vec3ToString(m.diffuse) + "}\n"
	logString += " - Specualar: Vector{" + transformations.Vec3ToString(m.specular) + "}\n"
	logString += " - Shininess: " + transformations.Float32ToString(m.shininess) + "\n"
	return logString
}

// GetDiffuse returns the diffuse color of the material
func (m *Material) GetAmbient() mgl32.Vec3 {
	return m.ambient
}

// GetDiffuse returns the diffuse color of the material
func (m *Material) GetDiffuse() mgl32.Vec3 {
	return m.diffuse
}

// GetSpecular returns the specular color of the material
func (m *Material) GetSpecular() mgl32.Vec3 {
	return m.specular
}

// GetShininess returns the shininess of the material
func (m *Material) GetShininess() float32 {
	return m.shininess
}

var (
	TestMaterialGreen = &Material{
		diffuse:   mgl32.Vec3{0, 1, 0},
		specular:  mgl32.Vec3{0, 1, 0},
		shininess: 0,
	}
	TestMaterialRed = &Material{
		diffuse:   mgl32.Vec3{1, 0, 0},
		specular:  mgl32.Vec3{1, 0, 0},
		shininess: 0,
	}
	Emerald = &Material{
		ambient:   mgl32.Vec3{0.0215, 0.1745, 0.0215},
		diffuse:   mgl32.Vec3{0.07568, 0.61424, 0.07568},
		specular:  mgl32.Vec3{0.633, 0.727811, 0.633},
		shininess: 0.6 * 128.0,
	}
	Jade = &Material{
		ambient:   mgl32.Vec3{0.135, 0.2225, 0.1575},
		diffuse:   mgl32.Vec3{0.54, 0.89, 0.63},
		specular:  mgl32.Vec3{0.316228, 0.316228, 0.316228},
		shininess: 0.1 * 128.0,
	}
	Obsidian = &Material{
		ambient:   mgl32.Vec3{0.05375, 0.05, 0.06625},
		diffuse:   mgl32.Vec3{0.18275, 0.17, 0.22525},
		specular:  mgl32.Vec3{0.332741, 0.328634, 0.346435},
		shininess: 0.3 * 128.0,
	}
	Pearl = &Material{
		ambient:   mgl32.Vec3{0.25, 0.20725, 0.20725},
		diffuse:   mgl32.Vec3{1, 0.829, 0.829},
		specular:  mgl32.Vec3{0.296648, 0.296648, 0.296648},
		shininess: 0.088 * 128.0,
	}
	Ruby = &Material{
		ambient:   mgl32.Vec3{0.1745, 0.01175, 0.01175},
		diffuse:   mgl32.Vec3{0.61424, 0.04136, 0.04136},
		specular:  mgl32.Vec3{0.727811, 0.626959, 0.626959},
		shininess: 0.6 * 128.0,
	}
	Turquoise = &Material{
		ambient:   mgl32.Vec3{0.1, 0.18725, 0.1745},
		diffuse:   mgl32.Vec3{0.396, 0.74151, 0.69102},
		specular:  mgl32.Vec3{0.297254, 0.30829, 0.306678},
		shininess: 0.1 * 128.0,
	}
	Brass = &Material{
		ambient:   mgl32.Vec3{0.329412, 0.223529, 0.027451},
		diffuse:   mgl32.Vec3{0.780392, 0.568627, 0.113725},
		specular:  mgl32.Vec3{0.992157, 0.941176, 0.807843},
		shininess: 0.21794872 * 128.0,
	}
	Bronze = &Material{
		ambient:   mgl32.Vec3{0.2125, 0.1275, 0.054},
		diffuse:   mgl32.Vec3{0.714, 0.4284, 0.18144},
		specular:  mgl32.Vec3{0.393548, 0.271906, 0.166721},
		shininess: 0.2 * 128.0,
	}
	Chrome = &Material{
		ambient:   mgl32.Vec3{0.25, 0.25, 0.25},
		diffuse:   mgl32.Vec3{0.4, 0.4, 0.4},
		specular:  mgl32.Vec3{0.774597, 0.774597, 0.774597},
		shininess: 0.6 * 128.0,
	}
	Copper = &Material{
		ambient:   mgl32.Vec3{0.19125, 0.0735, 0.0225},
		diffuse:   mgl32.Vec3{0.7038, 0.27048, 0.0828},
		specular:  mgl32.Vec3{0.256777, 0.137622, 0.086014},
		shininess: 0.1 * 128.0,
	}
	Gold = &Material{
		ambient:   mgl32.Vec3{0.24725, 0.1995, 0.0745},
		diffuse:   mgl32.Vec3{0.75164, 0.60648, 0.22648},
		specular:  mgl32.Vec3{0.628281, 0.555802, 0.366065},
		shininess: 0.4 * 128.0,
	}
	Silver = &Material{
		ambient:   mgl32.Vec3{0.19225, 0.19225, 0.19225},
		diffuse:   mgl32.Vec3{0.50754, 0.50754, 0.50754},
		specular:  mgl32.Vec3{0.508273, 0.508273, 0.508273},
		shininess: 0.4 * 128.0,
	}
	Blackplastic = &Material{
		ambient:   mgl32.Vec3{0.0, 0.0, 0.0},
		diffuse:   mgl32.Vec3{0.01, 0.01, 0.01},
		specular:  mgl32.Vec3{0.50, 0.50, 0.50},
		shininess: 0.25 * 128.0,
	}
	Cyanplastic = &Material{
		ambient:   mgl32.Vec3{0.0, 0.1, 0.06},
		diffuse:   mgl32.Vec3{0.0, 0.50980392, 0.50980392},
		specular:  mgl32.Vec3{0.50196078, 0.50196078, 0.50196078},
		shininess: 0.25 * 128.0,
	}
	Greenplastic = &Material{
		ambient:   mgl32.Vec3{0.0, 0.0, 0.0},
		diffuse:   mgl32.Vec3{0.1, 0.35, 0.1},
		specular:  mgl32.Vec3{0.45, 0.55, 0.45},
		shininess: 0.25 * 128.0,
	}
	Redplastic = &Material{
		ambient:   mgl32.Vec3{0.0, 0.0, 0.0},
		diffuse:   mgl32.Vec3{0.5, 0.0, 0.0},
		specular:  mgl32.Vec3{0.7, 0.6, 0.6},
		shininess: 0.25 * 128.0,
	}
	Whiteplastic = &Material{
		ambient:   mgl32.Vec3{0.0, 0.0, 0.0},
		diffuse:   mgl32.Vec3{0.55, 0.55, 0.55},
		specular:  mgl32.Vec3{0.70, 0.70, 0.70},
		shininess: 0.25 * 128.0,
	}
	Yellowplastic = &Material{
		ambient:   mgl32.Vec3{0.0, 0.0, 0.0},
		diffuse:   mgl32.Vec3{0.5, 0.5, 0.0},
		specular:  mgl32.Vec3{0.60, 0.60, 0.50},
		shininess: 0.25 * 128.0,
	}
	Blackrubber = &Material{
		ambient:   mgl32.Vec3{0.02, 0.02, 0.02},
		diffuse:   mgl32.Vec3{0.01, 0.01, 0.01},
		specular:  mgl32.Vec3{0.4, 0.4, 0.4},
		shininess: 0.078125 * 128.0,
	}
	Cyanrubber = &Material{
		ambient:   mgl32.Vec3{0.0, 0.05, 0.05},
		diffuse:   mgl32.Vec3{0.4, 0.5, 0.5},
		specular:  mgl32.Vec3{0.04, 0.7, 0.7},
		shininess: 0.078125 * 128.0,
	}
	Greenrubber = &Material{
		ambient:   mgl32.Vec3{0.0, 0.05, 0.0},
		diffuse:   mgl32.Vec3{0.4, 0.5, 0.4},
		specular:  mgl32.Vec3{0.04, 0.7, 0.04},
		shininess: 0.078125 * 128.0,
	}
	Redrubber = &Material{
		ambient:   mgl32.Vec3{0.05, 0.0, 0.0},
		diffuse:   mgl32.Vec3{0.5, 0.4, 0.4},
		specular:  mgl32.Vec3{0.7, 0.04, 0.04},
		shininess: 0.078125 * 128.0,
	}
	Whiterubber = &Material{
		ambient:   mgl32.Vec3{0.05, 0.05, 0.05},
		diffuse:   mgl32.Vec3{0.5, 0.5, 0.5},
		specular:  mgl32.Vec3{0.7, 0.7, 0.7},
		shininess: 0.078125 * 128.0,
	}
	Yellowrubber = &Material{
		ambient:   mgl32.Vec3{0.05, 0.05, 0.0},
		diffuse:   mgl32.Vec3{0.5, 0.5, 0.4},
		specular:  mgl32.Vec3{0.7, 0.7, 0.04},
		shininess: 0.078125 * 128.0,
	}
)
