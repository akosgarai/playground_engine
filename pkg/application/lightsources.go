package application

import (
	"github.com/akosgarai/playground_engine/pkg/interfaces"
)

type DirectionalLightSource struct {
	LightSource          interfaces.DirectionalLight
	DirectionUniformName string
	AmbientUniformName   string
	DiffuseUniformName   string
	SpecularUniformName  string
}
type PointLightSource struct {
	LightSource              interfaces.PointLight
	PositionUniformName      string
	AmbientUniformName       string
	DiffuseUniformName       string
	SpecularUniformName      string
	ConstantTermUniformName  string
	LinearTermUniformName    string
	QuadraticTermUniformName string
}
type SpotLightSource struct {
	LightSource              interfaces.SpotLight
	PositionUniformName      string
	DirectionUniformName     string
	AmbientUniformName       string
	DiffuseUniformName       string
	SpecularUniformName      string
	ConstantTermUniformName  string
	LinearTermUniformName    string
	QuadraticTermUniformName string
	CutoffUniformName        string
	OuterCutoffUniformName   string
}
