package model

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"os"

	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/texture"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Glyph struct {
	Debug    bool
	Width    int
	Height   int
	BearingX int
	BearingY int
	Advance  int
	tex      texture.Textures
}

// Build setups the Glyph parameters.
func (g *Glyph) Build(ch rune, ttf *truetype.Font, options *truetype.Options, filePath string, wrapper interfaces.GLWrapper) error {
	if g.Debug {
		fmt.Printf("Glyph.Build: %v\t'%c'\n\tOptions: %#v\n", ch, ch, options)
	}
	ttfFace := truetype.NewFace(ttf, options)
	gBnd, gAdv, ok := ttfFace.GlyphBounds(ch)
	if ok != true {
		msg := fmt.Sprintf("ttfFace.GlyphBounds(ch) not OK. %v\t'%c' ... Skipping\n", ch, ch)
		fmt.Printf("\t%s\n", msg)
		return errors.New(msg)
	}
	if g.Debug {
		fmt.Printf("\tgBnd: %#v\n\tgAdv: %#v\n", gBnd, gAdv)
	}
	g.Height = int((gBnd.Max.Y - gBnd.Min.Y) >> 6)
	g.Width = int((gBnd.Max.X - gBnd.Min.X) >> 6)
	g.Advance = int((gAdv) >> 6)
	if g.Width == 0 || g.Height == 0 {
		gBnd = ttf.Bounds(fixed.Int26_6(options.Size))
		fmt.Printf("\tNull handler (g.Height:%d, g.Width:%d)\n\tSet gBnd value to: %#v\n", g.Height, g.Width, gBnd)
		g.Width = int((gBnd.Max.X - gBnd.Min.X) >> 6)
		g.Height = int((gBnd.Max.Y - gBnd.Min.Y) >> 6)

		//above can sometimes yield 0 for font smaller than 48pt, 1 is minimum
		if g.Width == 0 || g.Height == 0 {
			fmt.Printf("\tFallback null handler (g.Height:%d, g.Width:%d)\n", g.Height, g.Width)
			g.Width = 1
			g.Height = 1
		}
	}
	g.BearingX = (int(gBnd.Min.X) >> 6)
	if g.Debug {
		fmt.Printf("\tg.Height:%d\n\tg.Width:%d\n", g.Height, g.Width)
	}
	gAscent := int(-gBnd.Min.Y) >> 6
	gDescent := int(gBnd.Max.Y) >> 6
	if g.Debug {
		fmt.Printf("\tgAscent: %d\n\tgdescent: %d\n", gAscent, gDescent)
	}
	g.BearingY = gDescent
	//create image to draw glyph
	background := g.rgba(image.Black)

	//create a freetype context for drawing
	c := g.context(ttf, background, image.White, options.DPI, options.Size)
	//set the glyph dot
	px := 0 - (int(gBnd.Min.X) >> 6)
	py := (gAscent)
	pt := freetype.Pt(px, py)
	fmt.Printf("\t(px, py): (%d, %d)\n\tpt: %#v\n", px, py, pt)
	// Draw the text from mask to image
	_, err := c.DrawString(string(ch), pt)
	if err != nil {
		return err
	}
	var tex texture.Textures
	// Generate texture
	tex.AddTextureRGBA(filePath, background, glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex", wrapper)
	g.tex = tex
	return nil
}
func (g *Glyph) context(ttf *truetype.Font, background *image.RGBA, foreground *image.Uniform, dpi, size float64) *freetype.Context {
	if g.Debug {
		fmt.Printf("\tCreating freetype context.\n\t\tDpi: %f\n\t\tsize: %f\n\t\tbounds: %#v\n", dpi, size, background.Bounds())
	}
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(ttf)
	c.SetFontSize(size)
	c.SetClip(background.Bounds())
	c.SetDst(background)
	c.SetSrc(foreground)
	c.SetHinting(font.HintingFull)
	return c
}

// rgba creates the background image and returns it.
func (g *Glyph) rgba(bg *image.Uniform) *image.RGBA {
	if g.Debug {
		fmt.Printf("\tCreating background image.\n")
	}
	rect := image.Rect(0, 0, int(g.Width), int(g.Height))
	rgba := image.NewRGBA(rect)
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	return rgba
}

type Charset struct {
	*BaseModel
	fonts   map[rune]*Glyph
	surface interfaces.Mesh
}

func LoadCharsetDebug(filePath string, low, high rune, scale float64, dpi float64, wrapper interfaces.GLWrapper) (*Charset, error) {
	fmt.Printf("Opening '%s'.\n", filePath)
	fd, err := os.Open(filePath)
	if err != nil {
		return &Charset{}, err
	}
	defer fd.Close()

	fmt.Printf("Reading '%s'.\n", filePath)
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return &Charset{}, err
	}

	fmt.Printf("Parsing '%s'.\n", filePath)
	ttf, err := truetype.Parse(data)
	if err != nil {
		return &Charset{}, err
	}
	fonts := make(map[rune]*Glyph)
	for ch := low; ch <= high; ch++ {
		g := &Glyph{Debug: true}
		options := &truetype.Options{
			Size:    scale,
			DPI:     dpi,
			Hinting: font.HintingFull,
		}
		err := g.Build(ch, ttf, options, filePath, wrapper)
		if err != nil {
			continue
		}
		fonts[ch] = g
	}
	m := New()
	return &Charset{
		BaseModel: m,
		fonts:     fonts,
	}, nil
}

// SetSurface sets the mesh that we are using for drawing.
func (c *Charset) SetSurface(mesh interfaces.Mesh) {
	c.surface = mesh
}
func (c *Charset) Print(text string, x, y, scale float32, wrapper interfaces.GLWrapper) {
	indices := []rune(text)
	fmt.Printf("The following text will be printed: '%s' as '%v'\n", text, indices)
	if len(indices) == 0 {
		return
	}
	// the low rune value from the LoadCharset function.
	lc := rune(32)
	var mshStore []interfaces.Mesh
	for i := range indices {
		runeIndex := indices[i]
		//skip runes that are not in font chacter range
		if int(runeIndex)-int(lc) > len(c.fonts) || runeIndex < lc {
			fmt.Printf("Skipping: %c %d\n", runeIndex, runeIndex)
			continue
		}
		ch := c.fonts[runeIndex]
		//calculate position and size for current rune
		xpos := x + float32(ch.BearingX)*scale
		ypos := y + float32(ch.Height-ch.BearingY)*scale
		w := float32(ch.Width) * scale
		h := float32(ch.Height) * scale
		rect := rectangle.NewExact(w, h)
		cols := []mgl32.Vec3{
			mgl32.Vec3{0.0, 1.0, 0.0},
		}
		v, i, _ := rect.TexturedColoredMeshInput(cols)
		rotTr := c.surface.RotationTransformation()
		position := mgl32.Vec3{x + float32(ch.BearingX+ch.Width/2)*scale, 0.01, y - float32(ch.BearingY-ch.Height/2)*scale}
		msh := mesh.NewTexturedColoredMesh(v, i, ch.tex, cols, wrapper)
		msh.SetPosition(mgl32.TransformCoordinate(position, rotTr))
		msh.SetParent(c.surface)
		mshStore = append(mshStore, msh)
		fmt.Printf("pos: %#v\nch: %#v\nw: %f, h: %f, xpos: %f, ypos: %f, adv: %f\n\n", position.Mul(scale), ch, w, h, xpos, ypos, float32(ch.Advance)*scale)
		x += float32(ch.Advance) * scale
	}
	for i := len(mshStore) - 1; i >= 0; i-- {
		c.Model.AddMesh(mshStore[i])
	}
}
