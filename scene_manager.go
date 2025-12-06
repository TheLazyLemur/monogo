package monogo

// SceneManager handles multiple scenes and scene transitions
type SceneManager struct {
	scenes      map[string]*Scene
	active      *Scene
	activeName  string
	persistent  []*GameObject // DontDestroyOnLoad objects
	onLoadFunc  func(name string, scene *Scene)
}

var sceneManager *SceneManager

func initSceneManager() {
	sceneManager = &SceneManager{
		scenes:     make(map[string]*Scene),
		persistent: make([]*GameObject, 0),
	}
}

// RegisterScene registers a scene with a name
func RegisterScene(name string, scene *Scene) {
	if sceneManager == nil {
		initSceneManager()
	}
	sceneManager.scenes[name] = scene
}

// LoadScene switches to a scene by name
func LoadScene(name string) {
	if sceneManager == nil {
		return
	}

	// Always create a fresh scene
	scene := NewScene()
	sceneManager.scenes[name] = scene

	// Move persistent objects to new scene
	for _, go_ := range sceneManager.persistent {
		go_.scene = scene
		scene.gameObjects = append(scene.gameObjects, go_)
	}

	sceneManager.active = scene
	sceneManager.activeName = name

	// Call scene load callback to populate scene
	if sceneManager.onLoadFunc != nil {
		sceneManager.onLoadFunc(name, scene)
	}

	// Only start non-persistent objects (persistent already started)
	for _, go_ := range scene.gameObjects {
		if !go_.persistent {
			go_.Start()
		}
	}
}

// GetActiveScene returns the currently active scene
func GetActiveScene() *Scene {
	if sceneManager == nil {
		return nil
	}
	return sceneManager.active
}

// GetActiveSceneName returns the name of the active scene
func GetActiveSceneName() string {
	if sceneManager == nil {
		return ""
	}
	return sceneManager.activeName
}

// DontDestroyOnLoad marks a GameObject to persist across scene loads
func DontDestroyOnLoad(go_ *GameObject) {
	if sceneManager == nil {
		initSceneManager()
	}

	go_.persistent = true

	// Check if already in persistent list
	for _, p := range sceneManager.persistent {
		if p == go_ {
			return
		}
	}
	sceneManager.persistent = append(sceneManager.persistent, go_)
}

// IsPersistent returns true if GameObject survives scene loads
func IsPersistent(go_ *GameObject) bool {
	return go_.persistent
}

// internal: set scene load callback
func setSceneLoadCallback(fn func(name string, scene *Scene)) {
	if sceneManager == nil {
		initSceneManager()
	}
	sceneManager.onLoadFunc = fn
}

// internal: get or create initial scene
func getOrCreateScene(name string) *Scene {
	if sceneManager == nil {
		initSceneManager()
	}

	if scene, exists := sceneManager.scenes[name]; exists {
		sceneManager.active = scene
		sceneManager.activeName = name
		return scene
	}

	scene := NewScene()
	sceneManager.scenes[name] = scene
	sceneManager.active = scene
	sceneManager.activeName = name
	return scene
}
