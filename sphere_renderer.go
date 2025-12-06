package monogo

import rl "github.com/gen2brain/raylib-go/raylib"

// SphereRenderer renders a colored sphere at the GameObject's transform
type SphereRenderer struct {
	BaseComponent
	Color         Color
	Radius        float32
	model         rl.Model
	hasModel      bool
	pendingShader *Shader
}

// NewSphereRenderer creates a sphere renderer with color and radius
func NewSphereRenderer(color Color, radius float32) *SphereRenderer {
	return &SphereRenderer{
		Color:  color,
		Radius: radius,
	}
}

// SetShader applies a custom shader to this renderer
func (s *SphereRenderer) SetShader(sh *Shader) *SphereRenderer {
	if s.hasModel {
		mats := s.model.GetMaterials()
		if len(mats) > 0 {
			mats[0].Shader = sh.GetRl()
		}
	} else {
		s.pendingShader = sh
	}
	return s
}

func (s *SphereRenderer) Start() {
	// Create sphere mesh and model (16 rings, 16 slices)
	mesh := rl.GenMeshSphere(s.Radius, 16, 16)
	s.model = rl.LoadModelFromMesh(mesh)
	s.hasModel = true

	// Apply pending shader if set before Start
	if s.pendingShader != nil {
		mats := s.model.GetMaterials()
		if len(mats) > 0 {
			mats[0].Shader = s.pendingShader.GetRl()
		}
		s.pendingShader = nil
	}
}

func (s *SphereRenderer) Draw() {
	if !s.hasModel {
		return
	}

	t := s.GameObject.Transform
	pos := t.WorldPosition()
	worldScale := t.WorldScale()

	// Apply full 3D rotation via transform matrix
	s.model.Transform = t.WorldRotationMatrix()

	// Draw with position and scale (rotation already in transform)
	rl.DrawModelEx(s.model, pos, Vec3(0, 1, 0), 0, worldScale, s.Color)
	rl.DrawModelWiresEx(s.model, pos, Vec3(0, 1, 0), 0, worldScale, Black)
}
