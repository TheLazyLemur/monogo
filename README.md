# MonoGo

## Getting Started

```bash
go get github.com/TheLazyLemur/monogo
```

You'll also need raylib-go. On macOS:
```bash
brew install raylib
```

### Minimal Example

```go
package main

import "github.com/TheLazyLemur/monogo"

type MyGame struct{}

func (g *MyGame) Setup(scene *monogo.Scene) {
    // Create a spinning cube
    cube := scene.Create("Cube")
    cube.Transform.Position = monogo.Vec3(0, 1, 0)
    cube.Add(monogo.NewCubeRenderer(monogo.Blue, monogo.Vec3(1, 1, 1)))
    cube.Add(&Spinner{Speed: 45})
}

// Spinner rotates whatever it's attached to
type Spinner struct {
    monogo.BaseComponent
    Speed float32
}

func (s *Spinner) Update(dt float32) {
    s.GameObject.Transform.Rotation.Y += s.Speed * dt
}

func main() {
    monogo.Run(&MyGame{}, monogo.Config{
        Title: "My Game",
    })
}
```

That's it. You get a window with a blue cube spinning. The engine handles the game loop, rendering, input - you just write components.

## Core Concepts

### GameObjects and Components

Everything in your game is a GameObject. On its own, a GameObject is just a transform (position, rotation, scale). You make it do things by adding Components.

```go
player := scene.Create("Player")
player.Transform.Position = monogo.Vec3(0, 0, 0)
player.Tag = "Player"

// Add a visual
player.Add(monogo.NewCubeRenderer(monogo.Green, monogo.Vec3(1, 2, 1)))

// Add behavior
player.Add(&PlayerController{MoveSpeed: 5})

// Add audio
player.Add(monogo.NewSoundSource("footsteps.wav"))
```

### Writing Components

Embed `BaseComponent` and implement what you need:

```go
type PlayerController struct {
    monogo.BaseComponent
    MoveSpeed float32
    health    int
}

// Start runs once when the scene loads
func (p *PlayerController) Start() {
    p.health = 100
}

// Update runs every frame
func (p *PlayerController) Update(dt float32) {
    h := monogo.Horizontal() // A/D or arrow keys
    v := monogo.Vertical()   // W/S or arrow keys

    p.GameObject.Transform.Position.X += h * p.MoveSpeed * dt
    p.GameObject.Transform.Position.Z += v * p.MoveSpeed * dt
}
```

For rendering, implement `Drawable` (3D) or `UIDrawable` (2D overlay):

```go
// Draw renders in 3D space
func (p *PlayerController) Draw() {
    if p.health < 20 {
        // flash red or something
    }
}

// DrawUI renders on screen (after 3D, no depth)
func (p *PlayerController) DrawUI() {
    monogo.DrawText(fmt.Sprintf("HP: %d", p.health), 10, 50, 20, monogo.White)
}
```

### Finding Things

```go
// By name
player := scene.Find("Player")

// By tag
enemies := scene.FindGameObjectsWithTag("Enemy")

// Get a specific component
unit := monogo.GetComponent[*Unit](someObject)
if unit != nil {
    unit.TakeDamage(10)
}

// Find all components of a type in the scene
allUnits := monogo.FindComponents[*Unit](scene)

// Find all GameObjects that have a specific component
targetObjects := monogo.FindGameObjectsWithComponent[*Target](scene)
```

### Transform Hierarchy

GameObjects can have parents. Child transforms are relative to their parent.

```go
tank := scene.Create("Tank")
tank.Transform.Position = monogo.Vec3(0, 0, 0)

turret := scene.Create("Turret")
turret.Transform.Position = monogo.Vec3(0, 1, 0) // 1 unit above parent
turret.SetParent(tank)

// When tank moves, turret moves with it
tank.Transform.Position.X = 5 // turret is now at (5, 1, 0) world space

// Get world position (accounts for all parents)
worldPos := turret.Transform.WorldPosition()
```

## Built-in Components

### Renderers

```go
// Colored cube (supports full 3D rotation)
cube := monogo.NewCubeRenderer(monogo.Red, monogo.Vec3(1, 1, 1))

// Colored sphere
sphere := monogo.NewSphereRenderer(monogo.Blue, 0.5) // radius 0.5

// 3D model from file
model := monogo.NewModelRendererFromFile("assets/robot.obj")
model.SetColor(monogo.White).SetScale(0.5)

// Apply custom shader
model.SetShader(myShader)
```

### Audio

```go
// Short sound effects (loads entire file)
sfx := monogo.NewSoundSource("shoot.wav")
sfx.Play()

// Music (streams from disk)
music := monogo.NewMusicSource("background.ogg")
music.Loop = true
music.PlayOnStart = true
music.SetVolume(0.5)
```

### Lua Scripting

For when you want to tweak behavior without recompiling:

```go
obj.Add(monogo.NewLuaComponent("scripts/enemy_ai.lua"))
```

```lua
-- scripts/enemy_ai.lua
function start()
    print("Enemy spawned at", transform.get_x(), transform.get_z())
end

function update(dt)
    local h = horizontal()
    transform.move(h * 5 * dt, 0, 0)

    if key_pressed(keys.space) then
        -- do something
    end
end
```

## Scene Management

### Multiple Scenes

```go
type MyGame struct{}

func (g *MyGame) Setup(scene *monogo.Scene) {
    // This is the initial scene - set up your main menu or whatever
    scene.Create("MainMenu").Add(&MainMenuUI{})
}

// Implement SceneLoader to handle scene transitions
func (g *MyGame) OnSceneLoad(name string, scene *monogo.Scene) {
    switch name {
    case "gameplay":
        scene.Create("Player").Add(&Player{})
        // spawn enemies, etc.
    case "gameover":
        scene.Create("GameOverUI").Add(&GameOverUI{})
    }
}
```

Then from anywhere:
```go
monogo.LoadScene("gameplay")
```

### Persistent Objects

Some things should survive scene loads (audio manager, score tracker, fade effects):

```go
func (g *MyGame) Setup(scene *monogo.Scene) {
    audioMgr := scene.Create("AudioManager")
    audioMgr.Add(&AudioManager{})
    monogo.DontDestroyOnLoad(audioMgr) // survives LoadScene calls
}
```

## Camera

The scene has a camera. You can control it directly or make it follow things.

```go
cam := scene.Camera

// Fixed position
cam.SetPosition(0, 10, 10)
cam.SetTarget(0, 0, 0)

// Follow a player
cam.Follow(playerObj)
cam.SetOffset(0, 5, 10) // behind and above

// Smooth follow
cam.FollowSmooth(playerObj, 3.0) // lerp speed
```

For RTS or custom cameras, write a component:

```go
type RTSCamera struct {
    monogo.BaseComponent
    x, z, zoom float32
}

func (c *RTSCamera) Update(dt float32) {
    c.x += monogo.Horizontal() * 20 * dt
    c.z -= monogo.Vertical() * 20 * dt
    c.zoom -= monogo.MouseWheel() * 5

    cam := c.GameObject.Scene().Camera
    cam.SetPosition(c.x, c.zoom, c.z + c.zoom*0.5)
    cam.SetTarget(c.x, 0, c.z)
}
```

### FPS Camera Utilities

For first-person cameras, use yaw/pitch controls:

```go
cam := scene.Camera

// Set orientation (pitch clamped to -89 to 89)
cam.SetYaw(yaw)
cam.SetPitch(pitch)

// Get direction vectors for movement
forward := cam.GetForward()  // Where camera looks
right := cam.GetRight()      // Strafe direction

// Use with mouse delta for look
dx, dy := monogo.MouseDelta()
cam.SetYaw(cam.GetYaw() + dx * sensitivity)
cam.SetPitch(cam.GetPitch() - dy * sensitivity)
```

### Raycasting

```go
cam := scene.Camera

// Where does the mouse point on the ground?
if groundPos, ok := cam.GetGroundPoint(); ok {
    // groundPos is where the mouse ray hits Y=0
}

// Get raw ray for custom intersection
ray := cam.GetMouseRay()

// Convert world position to screen coords
screenX, screenY := cam.WorldToScreen(someObject.Transform.Position)
```

## Shaders

### Custom Material Shaders

```go
// Load from strings
shader := monogo.LoadShaderFromMemory(vertexCode, fragmentCode)

// Apply to renderer
cube.SetShader(shader)

// Set uniforms
shader.SetFloat("time", float32(monogo.Time()))
shader.SetVec3("lightPos", 0, 10, 0)
```

### Post-Processing

```go
func (g *MyGame) Setup(scene *monogo.Scene) {
    scene.PostProcess = monogo.NewPostProcessChain()

    // Add vignette
    vignette := monogo.NewPostProcess(
        monogo.LoadShaderFromMemory("", monogo.VignetteFS),
    )
    vignette.GetShader().SetFloat("intensity", 0.5)
    scene.PostProcess.Add(vignette)

    // Add chromatic aberration
    chromatic := monogo.NewPostProcess(
        monogo.LoadShaderFromMemory("", monogo.ChromaticAberrationFS),
    )
    chromatic.GetShader().SetFloat("offset", 0.003)
    scene.PostProcess.Add(chromatic)
}
```

Built-in post-process shaders:
- `monogo.GrayscaleFS` - black and white
- `monogo.VignetteFS` - darkened edges (set `intensity` 0-1)
- `monogo.ChromaticAberrationFS` - RGB split (set `offset` ~0.001-0.01)

## Input

```go
// Keyboard
if monogo.KeyPressed(monogo.KeySpace) { /* just pressed */ }
if monogo.KeyHeld(monogo.KeyShift) { /* held down */ }
if monogo.KeyReleased(monogo.KeyEscape) { /* just released */ }

// Mouse
if monogo.MousePressed(monogo.MouseLeft) { }
if monogo.MouseHeld(monogo.MouseRight) { }
x, y := monogo.MousePosition()
dx, dy := monogo.MouseDelta()
wheel := monogo.MouseWheel()

// Mouse lock (for FPS controls)
monogo.LockMouse()              // Capture and hide cursor
monogo.UnlockMouse()            // Release cursor
if monogo.IsMouseLocked() { }   // Check state

// Axes (returns -1, 0, or 1)
h := monogo.Horizontal() // A/D or Left/Right
v := monogo.Vertical()   // W/S or Up/Down
custom := monogo.Axis(monogo.KeyQ, monogo.KeyE) // Q=-1, E=+1
```

## Helpers

```go
// Vectors
pos := monogo.Vec3(1, 2, 3)
uv := monogo.Vec2(0.5, 0.5)

// Colors
color := monogo.RGB(255, 128, 0)      // orange
color := monogo.RGBA(255, 0, 0, 128)  // semi-transparent red

// Predefined colors
monogo.White, monogo.Black, monogo.Red, monogo.Green, monogo.Blue
monogo.Yellow, monogo.Orange, monogo.Purple, monogo.Pink
monogo.Gray, monogo.LightGray, monogo.DarkGray
monogo.SkyBlue, monogo.Lime, monogo.Gold, monogo.Transparent

// Math
monogo.Lerp(0, 100, 0.5)           // 50
monogo.LerpVec3(posA, posB, 0.5)   // midpoint
monogo.Clamp(value, 0, 1)
monogo.Distance(posA, posB)
monogo.Normalize(direction)

// Random
monogo.RandomFloat(0, 1)
monogo.RandomInt(1, 6)  // dice roll

// Time
dt := monogo.Delta()         // frame time in seconds
elapsed := monogo.Time()     // total time since start

// Drawing (for UI/debug)
monogo.DrawText("Hello", 10, 10, 20, monogo.White)
monogo.DrawRect(10, 50, 100, 20, monogo.Red)
monogo.DrawRectLines(10, 50, 100, 20, monogo.White)
monogo.DrawCircle3D(x, y, z, radius, monogo.Green)
monogo.DrawLine3D(x1, y1, z1, x2, y2, z2, monogo.White)
```

## Game Interface

Your game struct can implement these interfaces:

```go
type Game interface {
    Setup(scene *Scene)  // Required - set up your initial scene
}

type SceneLoader interface {
    OnSceneLoad(name string, scene *Scene)  // Handle scene transitions
}

type Updater interface {
    Update(scene *Scene, dt float32)  // Per-frame logic outside components
}

type Renderer interface {
    Render(scene *Scene)  // Custom rendering after scene draws
}
```

## Config

```go
monogo.Run(&MyGame{}, monogo.Config{
    Title:           "My Game",
    Width:           1920,
    Height:          1080,
    FPS:             60,
    Fullscreen:      false,
    BackgroundColor: monogo.Black,
    ShowFPS:         true,
    ShowGrid:        true,
    GridSize:        20,
    InitialScene:    "main",
})
```

All fields are optional - sensible defaults are provided.

## Performance Notes

The Update and Draw loops are designed to be allocation-free. This means no GC pauses during gameplay. We cache drawable interfaces on GameObjects and use in-place filtering for destroyed objects.

If you're doing something performance-critical, avoid allocating in Update:
```go
// Bad - allocates every frame
func (e *Enemy) Update(dt float32) {
    targets := scene.FindGameObjectsWithTag("Player") // allocates slice
}

// Better - cache it
func (e *Enemy) Start() {
    e.player = e.GameObject.Scene().Find("Player")
}
```

## Public API Reference

### Entry Point

```go
func Run(game Game, configs ...Config)
```
Starts the game loop. Pass your game struct and optional config.

---

### Types

#### Config
```go
type Config struct {
    Title           string  // Window title (default: "MonoGo Game")
    Width           int     // Window width (default: 1280)
    Height          int     // Window height (default: 720)
    FPS             int     // Target FPS (default: 60)
    Fullscreen      bool    // Fullscreen mode (default: false)
    BackgroundColor Color   // Clear color (default: RayWhite)
    ShowFPS         bool    // Show FPS counter (default: true)
    ShowGrid        bool    // Show ground grid (default: true)
    GridSize        int     // Grid size in units (default: 20)
    InitialScene    string  // Starting scene name (default: "main")
}
```

#### Vector3
```go
type Vector3 = rl.Vector3  // {X, Y, Z float32}
```

#### Vector2
```go
type Vector2 = rl.Vector2  // {X, Y float32}
```

#### Color
```go
type Color = rl.Color  // {R, G, B, A uint8}
```

#### Model
```go
type Model = rl.Model
```

#### Key
```go
type Key = int32
```

#### MouseButton
```go
type MouseButton = rl.MouseButton
```

#### CameraMode
```go
type CameraMode int

const (
    CameraModeFixed CameraMode = iota
    CameraModeFollow
    CameraModeFollowSmooth
)
```

---

### Game Interfaces

```go
// Required - called once with initial scene
type Game interface {
    Setup(scene *Scene)
}

// Optional - called when LoadScene() is used
type SceneLoader interface {
    OnSceneLoad(name string, scene *Scene)
}

// Optional - called every frame before component updates
type Updater interface {
    Update(scene *Scene, dt float32)
}

// Optional - called after scene draws for custom rendering
type Renderer interface {
    Render(scene *Scene)
}
```

---

### Component Interfaces

```go
// Base interface all components implement
type Component interface {
    Start()
    Update(dt float32)
    GetGameObject() *GameObject
    SetGameObject(go_ *GameObject)
}

// Implement for 3D rendering
type Drawable interface {
    Draw()
}

// Implement for 2D UI rendering
type UIDrawable interface {
    DrawUI()
}
```

---

### Scene

```go
type Scene struct {
    Camera      *Camera
    PostProcess *PostProcessChain
}
```

| Method | Description |
|--------|-------------|
| `Create(name string) *GameObject` | Create and add a new GameObject |
| `AddGameObject(go_ *GameObject)` | Add existing GameObject to scene |
| `Instantiate(original *GameObject) *GameObject` | Clone a GameObject |
| `Destroy(go_ *GameObject)` | Mark GameObject for removal at frame end |
| `Find(name string) *GameObject` | Find first GameObject by name |
| `FindWithTag(tag string) *GameObject` | Find first GameObject by tag |
| `FindGameObjectsWithTag(tag string) []*GameObject` | Find all GameObjects with tag |
| `GetAllGameObjects() []*GameObject` | Get all GameObjects |
| `GetActiveGameObjects() []*GameObject` | Get all active GameObjects |
| `Count() int` | Number of GameObjects |

---

### GameObject

```go
type GameObject struct {
    Name      string
    Tag       string
    Transform *Transform
    Active    bool
}
```

| Method | Description |
|--------|-------------|
| `Add(c Component) *GameObject` | Add component (chainable) |
| `AddComponent(c Component)` | Add component |
| `GetAllComponents() []Component` | Get all components |
| `Scene() *Scene` | Get parent scene |
| `SetActive(active bool)` | Enable/disable GameObject |
| `Parent() *GameObject` | Get parent (nil if none) |
| `Children() []*GameObject` | Get children |
| `SetParent(parent *GameObject)` | Set parent |
| `AddChild(child *GameObject)` | Add child |
| `RemoveChild(child *GameObject)` | Remove child |

---

### Transform

```go
type Transform struct {
    Position Vector3  // Local position
    Rotation Vector3  // Local rotation (euler degrees)
    Scale    Vector3  // Local scale
}
```

| Method | Description |
|--------|-------------|
| `WorldPosition() Vector3` | Absolute world position |
| `WorldRotation() Vector3` | Absolute world rotation |
| `WorldScale() Vector3` | Absolute world scale |
| `WorldRotationMatrix() rl.Matrix` | Rotation as matrix (for rendering) |
| `SetWorldPosition(pos Vector3)` | Set position in world space |
| `Forward() Vector3` | Forward direction vector |
| `Right() Vector3` | Right direction vector |

---

### Camera

```go
type Camera struct {}
```

| Method | Description |
|--------|-------------|
| `SetPosition(x, y, z float32)` | Set camera position |
| `SetTarget(x, y, z float32)` | Set look-at target |
| `SetOffset(x, y, z float32)` | Set follow offset |
| `SetFOV(fov float32)` | Set field of view |
| `Position() Vector3` | Get current position |
| `Target() Vector3` | Get current target |
| `LookAt(x, y, z float32)` | Point camera at position |
| `Follow(go_ *GameObject)` | Snap-follow a GameObject |
| `FollowSmooth(go_ *GameObject, speed float32)` | Smooth-follow a GameObject |
| `SetYaw(deg float32)` | Set horizontal rotation (updates target) |
| `SetPitch(deg float32)` | Set vertical rotation (clamped -89 to 89) |
| `GetYaw() float32` | Get horizontal rotation |
| `GetPitch() float32` | Get vertical rotation |
| `GetForward() Vector3` | Get forward direction from yaw/pitch |
| `GetRight() Vector3` | Get right direction from yaw |
| `GetMouseRay() rl.Ray` | Get ray from camera through mouse |
| `GetGroundPoint() (Vector3, bool)` | Get mouse-ray intersection with Y=0 |
| `WorldToScreen(pos Vector3) (x, y int)` | Convert 3D to screen coords |

---

### BaseComponent

Embed this in your custom components:

```go
type BaseComponent struct {
    GameObject *GameObject
}
```

| Method | Description |
|--------|-------------|
| `GetGameObject() *GameObject` | Returns the attached GameObject |
| `SetGameObject(go_ *GameObject)` | Sets the attached GameObject |
| `Start()` | Empty default (override in your component) |
| `Update(dt float32)` | Empty default (override in your component) |

---

### CubeRenderer

```go
func NewCubeRenderer(color Color, size Vector3) *CubeRenderer
```

| Method | Description |
|--------|-------------|
| `SetShader(s *Shader) *CubeRenderer` | Apply custom shader |

| Field | Type | Description |
|-------|------|-------------|
| `Color` | `Color` | Cube color |
| `Size` | `Vector3` | Cube dimensions |

---

### SphereRenderer

```go
func NewSphereRenderer(color Color, radius float32) *SphereRenderer
```

| Method | Description |
|--------|-------------|
| `SetShader(s *Shader) *SphereRenderer` | Apply custom shader |

| Field | Type | Description |
|-------|------|-------------|
| `Color` | `Color` | Sphere color |
| `Radius` | `float32` | Sphere radius |

---

### ModelRenderer

```go
func NewModelRenderer(model Model) *ModelRenderer
func NewModelRendererFromFile(path string) *ModelRenderer
```

| Method | Description |
|--------|-------------|
| `SetColor(color Color) *ModelRenderer` | Set tint color |
| `SetScale(scale float32) *ModelRenderer` | Set model scale |
| `SetShader(s *Shader) *ModelRenderer` | Apply custom shader |
| `Unload()` | Free model resources |

| Field | Type | Description |
|-------|------|-------------|
| `Model` | `Model` | The 3D model |
| `Color` | `Color` | Tint color |
| `Scale` | `float32` | Scale multiplier |

---

### SoundSource

```go
func NewSoundSource(path string) *SoundSource   // Short clips
func NewMusicSource(path string) *SoundSource   // Streaming audio
```

| Method | Description |
|--------|-------------|
| `Play()` | Start playback |
| `Stop()` | Stop playback |
| `Pause()` | Pause playback |
| `Resume()` | Resume playback |
| `IsPlaying() bool` | Check if playing |
| `SetVolume(volume float32)` | Set volume (0.0-1.0) |
| `Unload()` | Free audio resources |

| Field | Type | Description |
|-------|------|-------------|
| `Volume` | `float32` | Playback volume |
| `Loop` | `bool` | Loop playback |
| `PlayOnStart` | `bool` | Auto-play on Start() |

---

### LuaComponent

```go
func NewLuaComponent(path string) *LuaComponent
```

Loads and runs a Lua script. The script can define `start()` and `update(dt)` functions.

**Lua API:**

| Global | Description |
|--------|-------------|
| `transform` | Table with position/rotation access |
| `transform.get_x()`, `get_y()`, `get_z()` | Get position |
| `transform.set_x(v)`, `set_y(v)`, `set_z(v)` | Set position |
| `transform.get_rot_y()`, `set_rot_y(v)` | Y rotation |
| `transform.move(dx, dy, dz)` | Move by delta |
| `key_pressed(key)` | Check key just pressed |
| `key_held(key)` | Check key held |
| `horizontal()` | Horizontal axis (-1 to 1) |
| `vertical()` | Vertical axis (-1 to 1) |
| `load_scene(name)` | Load a scene |
| `find(name)` | Find GameObject by name |
| `keys.space`, `keys.enter`, `keys.escape`, etc. | Key constants |
| `dt` | Delta time (updated each frame) |

---

### Shader

```go
func LoadShader(vsPath, fsPath string) *Shader
func LoadShaderFromMemory(vsCode, fsCode string) *Shader
```

| Method | Description |
|--------|-------------|
| `SetFloat(name string, value float32)` | Set float uniform |
| `SetVec2(name string, x, y float32)` | Set vec2 uniform |
| `SetVec3(name string, x, y, z float32)` | Set vec3 uniform |
| `SetVec4(name string, x, y, z, w float32)` | Set vec4 uniform |
| `SetInt(name string, value int32)` | Set int uniform |
| `SetTexture(name string, tex rl.Texture2D)` | Set texture uniform |
| `Begin()` | Start using shader |
| `End()` | Stop using shader |
| `Unload()` | Free shader resources |
| `GetRl() rl.Shader` | Get underlying raylib shader |

**Built-in Fragment Shaders:**
```go
const GrayscaleFS string           // Black and white
const VignetteFS string            // Darkened edges (uniform: intensity)
const ChromaticAberrationFS string // RGB split (uniform: offset)
```

---

### PostProcess

```go
func NewPostProcess(shader *Shader) *PostProcess
```

| Method | Description |
|--------|-------------|
| `SetEnabled(enabled bool)` | Enable/disable effect |
| `IsEnabled() bool` | Check if enabled |
| `GetShader() *Shader` | Get shader for setting uniforms |
| `Resize(width, height int)` | Update render texture size |
| `Unload()` | Free resources |

---

### PostProcessChain

```go
func NewPostProcessChain() *PostProcessChain
```

| Method | Description |
|--------|-------------|
| `Add(pp *PostProcess) *PostProcessChain` | Add effect to chain |
| `Remove(pp *PostProcess)` | Remove effect |
| `SetEnabled(enabled bool)` | Enable/disable chain |
| `IsEnabled() bool` | Check if enabled |
| `GetPass(index int) *PostProcess` | Get effect by index |
| `Resize(width, height int)` | Resize all passes |
| `Unload()` | Free all resources |

---

### Scene Management Functions

```go
func LoadScene(name string)                    // Switch to scene
func GetActiveScene() *Scene                   // Get current scene
func GetActiveSceneName() string               // Get current scene name
func RegisterScene(name string, scene *Scene)  // Pre-register a scene
func DontDestroyOnLoad(go_ *GameObject)        // Persist across loads
func IsPersistent(go_ *GameObject) bool        // Check persistence
```

---

### Component Functions

```go
func GetComponent[T Component](go_ *GameObject) T              // Get component by type
func HasComponent[T Component](go_ *GameObject) bool           // Check for component
func GetComponents[T Component](go_ *GameObject) []T           // Get all of type
func FindComponents[T Component](s *Scene) []T                 // Find all in scene
func FindGameObjectsWithComponent[T Component](s *Scene) []*GameObject  // Find objects by component
```

---

### Input Functions

**Keyboard:**
```go
func KeyPressed(key Key) bool   // Just pressed this frame
func KeyReleased(key Key) bool  // Just released this frame
func KeyHeld(key Key) bool      // Currently held down
```

**Mouse:**
```go
func MousePressed(button MouseButton) bool   // Just pressed
func MouseReleased(button MouseButton) bool  // Just released
func MouseHeld(button MouseButton) bool      // Currently held
func MousePosition() (x, y int)              // Screen position
func MouseDelta() (x, y float32)             // Movement since last frame
func MouseWheel() float32                    // Wheel movement
func LockMouse()                             // Capture and hide cursor
func UnlockMouse()                           // Release cursor
func IsMouseLocked() bool                    // Check if cursor captured
```

**Axes:**
```go
func Axis(negative, positive Key) float32  // Returns -1, 0, or 1
func Horizontal() float32                  // A/D or Left/Right
func Vertical() float32                    // W/S or Up/Down
```

**Key Constants:**
```go
KeyW, KeyA, KeyS, KeyD, KeyQ, KeyE, KeyR, KeyF
KeySpace, KeyShift, KeyCtrl, KeyAlt, KeyTab, KeyEscape, KeyEnter
KeyUp, KeyDown, KeyLeft, KeyRight
Key1, Key2, Key3, Key4, Key5
```

**Mouse Constants:**
```go
MouseLeft, MouseRight, MouseMiddle
```

---

### Math Functions

```go
func Vec3(x, y, z float32) Vector3
func Vec2(x, y float32) Vector2
func RGB(r, g, b uint8) Color
func RGBA(r, g, b, a uint8) Color

func Lerp(a, b, t float32) float32
func LerpVec3(a, b Vector3, t float32) Vector3
func Clamp(value, min, max float32) float32
func Distance(a, b Vector3) float32
func Normalize(v Vector3) Vector3
func Sin(x float32) float32  // radians
func Cos(x float32) float32  // radians

func RandomFloat(min, max float32) float32
func RandomInt(min, max int) int

func Delta() float32   // Frame time in seconds
func Time() float64    // Total elapsed time
```

---

### Drawing Functions

```go
func DrawText(text string, x, y, size int, color Color)
func DrawRect(x, y, width, height int, color Color)
func DrawRectLines(x, y, width, height int, color Color)
func DrawCircle3D(x, y, z, radius float32, color Color)
func DrawLine3D(x1, y1, z1, x2, y2, z2 float32, color Color)
```

---

### Utility Functions

```go
func ScreenWidth() int
func ScreenHeight() int
func LoadModel(path string) Model
func DefaultConfig() Config
```

---

### Color Constants

```go
White, Black, Gray, LightGray, DarkGray
Red, Green, Blue, Yellow, Orange, Purple, Pink
Brown, Beige, SkyBlue, Lime, Gold, Maroon, Magenta, Violet
RayWhite, Transparent
```

---

## License

MIT - do whatever you want with it.
