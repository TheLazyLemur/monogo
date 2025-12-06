package monogo

// Scene holds all GameObjects
type Scene struct {
	Camera         *Camera
	PostProcess    *PostProcessChain // optional post-processing effects
	gameObjects    []*GameObject
	toDestroy      []*GameObject
	pendingObjects []*GameObject
}

// NewScene creates a new empty scene
func NewScene() *Scene {
	return &Scene{
		Camera:         NewCamera(),
		gameObjects:    make([]*GameObject, 0),
		toDestroy:      make([]*GameObject, 0),
		pendingObjects: make([]*GameObject, 0),
	}
}

// Create creates a new GameObject and adds it to the scene
func (s *Scene) Create(name string) *GameObject {
	go_ := NewGameObject(name)
	s.AddGameObject(go_)
	return go_
}

// AddGameObject adds a GameObject to the scene
func (s *Scene) AddGameObject(go_ *GameObject) {
	go_.scene = s
	s.gameObjects = append(s.gameObjects, go_)
}

// Instantiate creates a clone of a GameObject and adds it to the scene
func (s *Scene) Instantiate(original *GameObject) *GameObject {
	clone := NewGameObject(original.Name + "(Clone)")
	clone.Tag = original.Tag
	clone.Transform.Position = original.Transform.Position
	clone.Transform.Rotation = original.Transform.Rotation
	clone.Transform.Scale = original.Transform.Scale
	clone.Active = original.Active
	s.pendingObjects = append(s.pendingObjects, clone)
	clone.scene = s
	return clone
}

// Destroy marks a GameObject for removal at end of frame
func (s *Scene) Destroy(go_ *GameObject) {
	s.toDestroy = append(s.toDestroy, go_)
}

// Find returns the first GameObject with the given name, or nil
func (s *Scene) Find(name string) *GameObject {
	for _, go_ := range s.gameObjects {
		if go_.Name == name {
			return go_
		}
	}
	return nil
}

// FindWithTag returns the first GameObject with the given tag, or nil
func (s *Scene) FindWithTag(tag string) *GameObject {
	for _, go_ := range s.gameObjects {
		if go_.Tag == tag {
			return go_
		}
	}
	return nil
}

// FindGameObjectsWithTag returns all GameObjects with the given tag
func (s *Scene) FindGameObjectsWithTag(tag string) []*GameObject {
	result := make([]*GameObject, 0)
	for _, go_ := range s.gameObjects {
		if go_.Tag == tag {
			result = append(result, go_)
		}
	}
	return result
}

// FindComponents returns all components of type T in the scene
func FindComponents[T Component](s *Scene) []T {
	var result []T
	for _, go_ := range s.gameObjects {
		for _, c := range go_.components {
			if typed, ok := c.(T); ok {
				result = append(result, typed)
			}
		}
	}
	return result
}

// FindGameObjectsWithComponent returns all GameObjects that have a component of type T
func FindGameObjectsWithComponent[T Component](s *Scene) []*GameObject {
	var result []*GameObject
	for _, go_ := range s.gameObjects {
		for _, c := range go_.components {
			if _, ok := c.(T); ok {
				result = append(result, go_)
				break // Only add each GameObject once
			}
		}
	}
	return result
}

// GetAllGameObjects returns all GameObjects in the scene
func (s *Scene) GetAllGameObjects() []*GameObject {
	return s.gameObjects
}

// GetActiveGameObjects returns all active GameObjects
func (s *Scene) GetActiveGameObjects() []*GameObject {
	result := make([]*GameObject, 0)
	for _, go_ := range s.gameObjects {
		if go_.Active {
			result = append(result, go_)
		}
	}
	return result
}

// Count returns the number of GameObjects in the scene
func (s *Scene) Count() int {
	return len(s.gameObjects)
}

// Start calls Start on all GameObjects
func (s *Scene) Start() {
	for _, go_ := range s.gameObjects {
		go_.Start()
	}
}

// Update calls Update on all GameObjects with delta time
func (s *Scene) Update(dt float32) {
	// Add pending objects
	for _, go_ := range s.pendingObjects {
		s.gameObjects = append(s.gameObjects, go_)
		go_.Start()
	}
	s.pendingObjects = s.pendingObjects[:0]

	// Update all objects
	for _, go_ := range s.gameObjects {
		go_.Update(dt)
	}

	// Process destroy queue
	s.processDestroyQueue()
}

func (s *Scene) processDestroyQueue() {
	if len(s.toDestroy) == 0 {
		return
	}

	// Filter in-place without allocating
	n := 0
	for _, go_ := range s.gameObjects {
		if !s.isMarkedForDestroy(go_) {
			s.gameObjects[n] = go_
			n++
		} else {
			go_.scene = nil
		}
	}
	s.gameObjects = s.gameObjects[:n]
	s.toDestroy = s.toDestroy[:0]
}

func (s *Scene) isMarkedForDestroy(go_ *GameObject) bool {
	for _, d := range s.toDestroy {
		if d == go_ {
			return true
		}
	}
	return false
}

// Draw renders all drawable components (3D)
func (s *Scene) Draw() {
	for _, go_ := range s.gameObjects {
		if !go_.Active {
			continue
		}
		for _, d := range go_.drawables {
			d.Draw()
		}
	}
}

// DrawUI renders all UI drawable components (2D overlay)
func (s *Scene) DrawUI() {
	for _, go_ := range s.gameObjects {
		if !go_.Active {
			continue
		}
		for _, u := range go_.uiDrawables {
			u.DrawUI()
		}
	}
}
