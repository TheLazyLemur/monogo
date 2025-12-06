package monogo

import rl "github.com/gen2brain/raylib-go/raylib"

// Shader wraps raylib shader
type Shader struct {
	rl rl.Shader
}

// LoadShader loads a shader from vertex and fragment files
// Pass empty string for either to use default
func LoadShader(vsPath, fsPath string) *Shader {
	return &Shader{rl: rl.LoadShader(vsPath, fsPath)}
}

// LoadShaderFromMemory loads shader from code strings
func LoadShaderFromMemory(vsCode, fsCode string) *Shader {
	return &Shader{rl: rl.LoadShaderFromMemory(vsCode, fsCode)}
}

// SetFloat sets a float uniform
func (s *Shader) SetFloat(name string, value float32) {
	loc := rl.GetShaderLocation(s.rl, name)
	rl.SetShaderValue(s.rl, loc, []float32{value}, rl.ShaderUniformFloat)
}

// SetVec2 sets a vec2 uniform
func (s *Shader) SetVec2(name string, x, y float32) {
	loc := rl.GetShaderLocation(s.rl, name)
	rl.SetShaderValue(s.rl, loc, []float32{x, y}, rl.ShaderUniformVec2)
}

// SetVec3 sets a vec3 uniform
func (s *Shader) SetVec3(name string, x, y, z float32) {
	loc := rl.GetShaderLocation(s.rl, name)
	rl.SetShaderValue(s.rl, loc, []float32{x, y, z}, rl.ShaderUniformVec3)
}

// SetVec4 sets a vec4 uniform
func (s *Shader) SetVec4(name string, x, y, z, w float32) {
	loc := rl.GetShaderLocation(s.rl, name)
	rl.SetShaderValue(s.rl, loc, []float32{x, y, z, w}, rl.ShaderUniformVec4)
}

// SetInt sets an int uniform
func (s *Shader) SetInt(name string, value int32) {
	loc := rl.GetShaderLocation(s.rl, name)
	rl.SetShaderValue(s.rl, loc, []float32{float32(value)}, rl.ShaderUniformInt)
}

// SetTexture sets a texture uniform
func (s *Shader) SetTexture(name string, tex rl.Texture2D) {
	loc := rl.GetShaderLocation(s.rl, name)
	rl.SetShaderValueTexture(s.rl, loc, tex)
}

// Begin starts using this shader
func (s *Shader) Begin() {
	rl.BeginShaderMode(s.rl)
}

// End stops using this shader
func (s *Shader) End() {
	rl.EndShaderMode()
}

// Unload frees shader resources
func (s *Shader) Unload() {
	rl.UnloadShader(s.rl)
}

// GetRl returns the underlying raylib shader (for advanced use)
func (s *Shader) GetRl() rl.Shader {
	return s.rl
}

// PostProcess handles full-screen post-processing effects
type PostProcess struct {
	shader  *Shader
	target  rl.RenderTexture2D
	enabled bool
}

// NewPostProcess creates a post-processing pass
func NewPostProcess(shader *Shader) *PostProcess {
	w, h := int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight())
	return &PostProcess{
		shader:  shader,
		target:  rl.LoadRenderTexture(w, h),
		enabled: true,
	}
}

// SetEnabled enables/disables the post-process effect
func (p *PostProcess) SetEnabled(enabled bool) {
	p.enabled = enabled
}

// IsEnabled returns whether the effect is enabled
func (p *PostProcess) IsEnabled() bool {
	return p.enabled
}

// GetShader returns the shader for setting uniforms
func (p *PostProcess) GetShader() *Shader {
	return p.shader
}

// Begin starts rendering to the post-process buffer
func (p *PostProcess) Begin() {
	rl.BeginTextureMode(p.target)
}

// End finishes rendering and applies the effect
func (p *PostProcess) End() {
	rl.EndTextureMode()

	if p.enabled && p.shader != nil {
		rl.BeginShaderMode(p.shader.rl)
	}

	// Draw the render texture flipped
	rl.DrawTextureRec(
		p.target.Texture,
		rl.NewRectangle(0, 0, float32(p.target.Texture.Width), float32(-p.target.Texture.Height)),
		rl.NewVector2(0, 0),
		rl.White,
	)

	if p.enabled && p.shader != nil {
		rl.EndShaderMode()
	}
}

// Resize updates the render texture size (call on window resize)
func (p *PostProcess) Resize(width, height int) {
	rl.UnloadRenderTexture(p.target)
	p.target = rl.LoadRenderTexture(int32(width), int32(height))
}

// Unload frees resources
func (p *PostProcess) Unload() {
	rl.UnloadRenderTexture(p.target)
	if p.shader != nil {
		p.shader.Unload()
	}
}

// Common post-process shader code

// GrayscaleFS is a simple grayscale fragment shader
const GrayscaleFS = `
#version 330
in vec2 fragTexCoord;
uniform sampler2D texture0;
out vec4 finalColor;
void main() {
    vec4 color = texture(texture0, fragTexCoord);
    float gray = dot(color.rgb, vec3(0.299, 0.587, 0.114));
    finalColor = vec4(gray, gray, gray, color.a);
}
`

// VignetteFS adds a vignette effect
const VignetteFS = `
#version 330
in vec2 fragTexCoord;
uniform sampler2D texture0;
uniform float intensity; // 0.0 to 1.0
out vec4 finalColor;
void main() {
    vec4 color = texture(texture0, fragTexCoord);
    vec2 uv = fragTexCoord;
    uv *= 1.0 - uv.yx;
    float vig = uv.x * uv.y * 15.0;
    vig = pow(vig, intensity * 0.5);
    finalColor = vec4(color.rgb * vig, color.a);
}
`

// ChromaticAberrationFS adds chromatic aberration
const ChromaticAberrationFS = `
#version 330
in vec2 fragTexCoord;
uniform sampler2D texture0;
uniform float offset; // 0.001 to 0.01
out vec4 finalColor;
void main() {
    float r = texture(texture0, fragTexCoord + vec2(offset, 0.0)).r;
    float g = texture(texture0, fragTexCoord).g;
    float b = texture(texture0, fragTexCoord - vec2(offset, 0.0)).b;
    float a = texture(texture0, fragTexCoord).a;
    finalColor = vec4(r, g, b, a);
}
`

// PostProcessChain chains multiple post-process effects
type PostProcessChain struct {
	passes  []*PostProcess
	enabled bool
}

// NewPostProcessChain creates an empty post-process chain
func NewPostProcessChain() *PostProcessChain {
	return &PostProcessChain{
		passes:  make([]*PostProcess, 0),
		enabled: true,
	}
}

// Add appends a post-process pass to the chain
func (c *PostProcessChain) Add(pp *PostProcess) *PostProcessChain {
	c.passes = append(c.passes, pp)
	return c
}

// Remove removes a post-process pass from the chain
func (c *PostProcessChain) Remove(pp *PostProcess) {
	for i, p := range c.passes {
		if p == pp {
			c.passes = append(c.passes[:i], c.passes[i+1:]...)
			return
		}
	}
}

// SetEnabled enables/disables the entire chain
func (c *PostProcessChain) SetEnabled(enabled bool) {
	c.enabled = enabled
}

// IsEnabled returns whether the chain is enabled
func (c *PostProcessChain) IsEnabled() bool {
	return c.enabled && len(c.passes) > 0
}

// Begin starts rendering to the first pass buffer
func (c *PostProcessChain) Begin() {
	if len(c.passes) == 0 {
		return
	}
	c.passes[0].Begin()
}

// End chains through all passes and renders final result
func (c *PostProcessChain) End() {
	if len(c.passes) == 0 {
		return
	}

	// End the first pass's texture capture
	rl.EndTextureMode()

	// Chain through passes: each pass reads from previous and writes to its buffer
	for i := 0; i < len(c.passes); i++ {
		pp := c.passes[i]
		isLast := i == len(c.passes)-1

		if !isLast {
			// Render to next pass's buffer
			c.passes[i+1].Begin()
		} else {
			// Last pass renders to screen
			rl.BeginDrawing()
		}

		// Apply this pass's shader
		if pp.enabled && pp.shader != nil {
			rl.BeginShaderMode(pp.shader.rl)
		}

		// Draw previous buffer (or first pass's buffer for i==0)
		rl.DrawTextureRec(
			pp.target.Texture,
			rl.NewRectangle(0, 0, float32(pp.target.Texture.Width), float32(-pp.target.Texture.Height)),
			rl.NewVector2(0, 0),
			rl.White,
		)

		if pp.enabled && pp.shader != nil {
			rl.EndShaderMode()
		}

		if !isLast {
			rl.EndTextureMode()
		}
	}
}

// Resize updates all render texture sizes
func (c *PostProcessChain) Resize(width, height int) {
	for _, pp := range c.passes {
		pp.Resize(width, height)
	}
}

// Unload frees all resources
func (c *PostProcessChain) Unload() {
	for _, pp := range c.passes {
		pp.Unload()
	}
	c.passes = nil
}

// GetPass returns a specific pass by index for setting uniforms
func (c *PostProcessChain) GetPass(index int) *PostProcess {
	if index < 0 || index >= len(c.passes) {
		return nil
	}
	return c.passes[index]
}
