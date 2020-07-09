package texture

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
)

// LoadImageFromFile takes a filepath string argument.
// It loads the file, decodes it as PNG or jpg, and returns the image and error
func loadImageFromFile(path string) (image.Image, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()
	img, _, err := image.Decode(imgFile)
	return img, err

}

// GenTextures returns the generated uint32 name of the texture.
func genTextures(wrapper interfaces.GLWrapper) uint32 {
	var id uint32
	wrapper.GenTextures(1, &id)
	return id
}

type Texture struct {
	// The id of the texture. eg gl.TEXTURE0
	Id uint32
	// The generated name that was given by the GenTextures command
	TextureName uint32
	// The target that we use for BindTexture. (eg: TEXTURE_2D)
	TargetId uint32

	// The Uniform name of the texture
	UniformName string
	Wrapper     interfaces.GLWrapper
	// The path of the file - used for exporting
	FilePath string
}

// Bing calls Activtexture with it's `Id`, then BindTexture with `TargetId` and `TextureName`.
func (t *Texture) Bind() {
	t.Wrapper.ActiveTexture(t.Id)
	t.Wrapper.BindTexture(t.TargetId, t.TextureName)
}

// UnBind binds the default texture. It is the cleanup for the texture.
func (t *Texture) UnBind() {
	t.Wrapper.BindTexture(t.TargetId, glwrapper.TEXTURE0)
}

type Textures []*Texture

// TransparentTexture creates a transparent surface with the given dimensions and returns the transparent texture.
func (t *Textures) TransparentTexture(width, height int, uniformName string, wrapper interfaces.GLWrapper) {
	upLeft := image.Point{0, 0}
	bottomRight := image.Point{width, height}
	rgba := image.NewRGBA(image.Rectangle{upLeft, bottomRight})
	rgba.Set(0, 0, color.RGBA{255, 255, 255, 0})

	t.AddTextureRGBA("transparent-gen", rgba, glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, uniformName, wrapper)
}

// AddTextureRGBA gets an RGBA and further necessary parameters, sets up a Texture, and appends it.
func (t *Textures) AddTextureRGBA(filePath string, rgba *image.RGBA, wrapR, wrapS, minificationFilter, magnificationFilter int32, uniformName string, wrapper interfaces.GLWrapper) {
	tex := &Texture{
		TextureName: genTextures(wrapper),
		TargetId:    glwrapper.TEXTURE_2D,
		Id:          glwrapper.TEXTURE0 + uint32(len(*t)),
		UniformName: uniformName,
		Wrapper:     wrapper,
		FilePath:    filePath,
	}

	tex.Bind()
	defer tex.UnBind()

	tex.Wrapper.TexParameteri(glwrapper.TEXTURE_2D, glwrapper.TEXTURE_WRAP_R, wrapR)
	tex.Wrapper.TexParameteri(glwrapper.TEXTURE_2D, glwrapper.TEXTURE_WRAP_S, wrapS)
	tex.Wrapper.TexParameteri(glwrapper.TEXTURE_2D, glwrapper.TEXTURE_MIN_FILTER, minificationFilter)
	tex.Wrapper.TexParameteri(glwrapper.TEXTURE_2D, glwrapper.TEXTURE_MAG_FILTER, magnificationFilter)

	tex.Wrapper.TexImage2D(tex.TargetId, 0, glwrapper.RGBA, int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y), 0, glwrapper.RGBA, uint32(glwrapper.UNSIGNED_BYTE), tex.Wrapper.Ptr(rgba.Pix))

	tex.Wrapper.GenerateMipmap(tex.TextureName)

	*t = append(*t, tex)
}

// AddTexture gets a filepath and further necessary parameters, loads the image from the filepath,
// validates it and sets up the RGBA, and calls AddTextureRGBA function
func (t *Textures) AddTexture(filePath string, wrapR, wrapS, minificationFilter, magnificationFilter int32, uniformName string, wrapper interfaces.GLWrapper) {
	img, err := loadImageFromFile(filePath)
	if err != nil {
		panic(err)
	}
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("not 32 bit color")
	}
	t.AddTextureRGBA(filePath, rgba, wrapR, wrapS, minificationFilter, magnificationFilter, uniformName, wrapper)
}
func (t *Textures) AddCubeMapTexture(directoryPath string, wrapR, wrapS, wrapT, minificationFilter, magnificationFilter int32, uniformName string, wrapper interfaces.GLWrapper) {
	tex := &Texture{
		TextureName: genTextures(wrapper),
		TargetId:    glwrapper.TEXTURE_CUBE_MAP,
		Id:          glwrapper.TEXTURE0 + uint32(len(*t)),
		UniformName: uniformName,
		Wrapper:     wrapper,
	}
	tex.Bind()
	defer tex.UnBind()
	fileNames := []string{"skybox-right.png", "skybox-left.png", "skybox-top.png", "skybox-bottom.png", "skybox-front.png", "skybox-back.png"}

	for index, file := range fileNames {
		img, err := loadImageFromFile(directoryPath + "/" + file)
		if err != nil {
			panic(err)
		}
		rgba := image.NewRGBA(img.Bounds())
		draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
		if rgba.Stride != rgba.Rect.Size().X*4 {
			panic("not 32 bit color")
		}
		tex.Wrapper.TexImage2D(glwrapper.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(index), 0, glwrapper.RGBA, int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y), 0, glwrapper.RGBA, uint32(glwrapper.UNSIGNED_BYTE), tex.Wrapper.Ptr(rgba.Pix))
	}
	tex.Wrapper.TexParameteri(glwrapper.TEXTURE_CUBE_MAP, glwrapper.TEXTURE_WRAP_R, wrapR)
	tex.Wrapper.TexParameteri(glwrapper.TEXTURE_CUBE_MAP, glwrapper.TEXTURE_WRAP_S, wrapS)
	tex.Wrapper.TexParameteri(glwrapper.TEXTURE_CUBE_MAP, glwrapper.TEXTURE_WRAP_T, wrapT)
	tex.Wrapper.TexParameteri(glwrapper.TEXTURE_CUBE_MAP, glwrapper.TEXTURE_MIN_FILTER, minificationFilter)
	tex.Wrapper.TexParameteri(glwrapper.TEXTURE_CUBE_MAP, glwrapper.TEXTURE_MAG_FILTER, magnificationFilter)

	*t = append(*t, tex)
}

func (t Textures) UnBind() {
	for i, _ := range t {
		t[i].UnBind()
	}
}
