package monogo

import rl "github.com/gen2brain/raylib-go/raylib"

// Game interface - implement this to create a game
type Game interface {
	Setup(scene *Scene) // Called once with initial scene
}

// SceneLoader is optional - implement to handle scene transitions
type SceneLoader interface {
	OnSceneLoad(name string, scene *Scene)
}

// Updater is optional - implement for per-frame logic outside components
type Updater interface {
	Update(scene *Scene, dt float32)
}

// Renderer is optional - implement for custom rendering after scene draw (UI, etc)
type Renderer interface {
	Render(scene *Scene)
}

// Config holds window and app settings
type Config struct {
	Title           string
	Width           int
	Height          int
	FPS             int
	Fullscreen      bool
	BackgroundColor Color
	ShowFPS         bool
	ShowGrid        bool
	GridSize        int
	InitialScene    string // Name of initial scene (default: "main")
}

// DefaultConfig returns sensible defaults
func DefaultConfig() Config {
	return Config{
		Title:           "MonoGo Game",
		Width:           1280,
		Height:          720,
		FPS:             60,
		BackgroundColor: RayWhite,
		ShowFPS:         true,
		ShowGrid:        true,
		GridSize:        20,
		InitialScene:    "main",
	}
}

// Run starts the game - this is the main entry point
func Run(game Game, configs ...Config) {
	config := DefaultConfig()
	if len(configs) > 0 {
		c := configs[0]
		if c.Title != "" {
			config.Title = c.Title
		}
		if c.Width != 0 {
			config.Width = c.Width
		}
		if c.Height != 0 {
			config.Height = c.Height
		}
		if c.FPS != 0 {
			config.FPS = c.FPS
		}
		config.Fullscreen = c.Fullscreen
		if c.BackgroundColor != (Color{}) {
			config.BackgroundColor = c.BackgroundColor
		}
		config.ShowFPS = c.ShowFPS
		config.ShowGrid = c.ShowGrid
		if c.GridSize != 0 {
			config.GridSize = c.GridSize
		}
		if c.InitialScene != "" {
			config.InitialScene = c.InitialScene
		}
	}

	// Init scene manager
	initSceneManager()

	// Set up scene load callback if game implements SceneLoader
	if loader, ok := game.(SceneLoader); ok {
		setSceneLoadCallback(loader.OnSceneLoad)
	}

	// Create initial scene
	scene := getOrCreateScene(config.InitialScene)

	// Init window
	rl.InitWindow(int32(config.Width), int32(config.Height), config.Title)
	defer rl.CloseWindow()

	if config.Fullscreen {
		rl.ToggleFullscreen()
	}

	rl.SetTargetFPS(int32(config.FPS))
	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	// Setup initial scene
	game.Setup(scene)
	scene.Start()

	// Game loop
	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()

		// Get active scene (may have changed via LoadScene)
		scene = GetActiveScene()

		// Update camera
		scene.Camera.update(dt)

		// Update game logic
		if updater, ok := game.(Updater); ok {
			updater.Update(scene, dt)
		}
		scene.Update(dt)

		// Draw - with optional post-processing
		usePostProcess := scene.PostProcess != nil && scene.PostProcess.IsEnabled()

		if usePostProcess {
			scene.PostProcess.Begin()
		}

		rl.BeginDrawing()
		rl.ClearBackground(config.BackgroundColor)

		rl.BeginMode3D(scene.Camera.getRl())

		if config.ShowGrid {
			rl.DrawGrid(int32(config.GridSize), 1.0)
		}

		scene.Draw()

		rl.EndMode3D()

		// Component UI rendering
		scene.DrawUI()

		// Game-level UI rendering
		if renderer, ok := game.(Renderer); ok {
			renderer.Render(scene)
		}

		if config.ShowFPS {
			rl.DrawFPS(10, 10)
		}

		if usePostProcess {
			// End() handles the final drawing
			scene.PostProcess.End()
		}

		rl.EndDrawing()
	}
}

// ScreenWidth returns current screen width
func ScreenWidth() int {
	return int(rl.GetScreenWidth())
}

// ScreenHeight returns current screen height
func ScreenHeight() int {
	return int(rl.GetScreenHeight())
}

// DrawText draws text on screen (for UI in Render)
func DrawText(text string, x, y, size int, color Color) {
	rl.DrawText(text, int32(x), int32(y), int32(size), color)
}
