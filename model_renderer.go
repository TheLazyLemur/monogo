package monogo

import rl "github.com/gen2brain/raylib-go/raylib"

// ModelRenderer renders a 3D model at the GameObject's transform
type ModelRenderer struct {
	BaseComponent
	Model Model
	Color Color
	Scale float32
}

// NewModelRenderer creates a renderer with a model
func NewModelRenderer(model Model) *ModelRenderer {
	return &ModelRenderer{
		Model: model,
		Color: White,
		Scale: 1.0,
	}
}

// NewModelRendererFromFile loads a model from file
func NewModelRendererFromFile(path string) *ModelRenderer {
	return &ModelRenderer{
		Model: rl.LoadModel(path),
		Color: White,
		Scale: 1.0,
	}
}

// SetColor sets the tint color
func (m *ModelRenderer) SetColor(color Color) *ModelRenderer {
	m.Color = color
	return m
}

// SetScale sets the model scale
func (m *ModelRenderer) SetScale(scale float32) *ModelRenderer {
	m.Scale = scale
	return m
}

// SetShader applies a custom shader to this renderer
func (m *ModelRenderer) SetShader(s *Shader) *ModelRenderer {
	mats := m.Model.GetMaterials()
	if len(mats) > 0 {
		mats[0].Shader = s.GetRl()
	}
	return m
}

func (m *ModelRenderer) Draw() {
	t := m.GameObject.Transform

	// Apply full 3D rotation via transform matrix
	m.Model.Transform = t.WorldRotationMatrix()

	scale := m.Scale * t.WorldScale().X
	rl.DrawModel(m.Model, t.WorldPosition(), scale, m.Color)
}

// Unload frees the model resources - call when done
func (m *ModelRenderer) Unload() {
	rl.UnloadModel(m.Model)
}
