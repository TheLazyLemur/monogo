package monogo

// Component interface - implement this for custom behaviors
type Component interface {
	Start()
	Update(dt float32)
	GetGameObject() *GameObject
	SetGameObject(go_ *GameObject)
}

// BaseComponent provides common functionality - embed this in your components
type BaseComponent struct {
	GameObject *GameObject
}

func (b *BaseComponent) GetGameObject() *GameObject {
	return b.GameObject
}

func (b *BaseComponent) SetGameObject(go_ *GameObject) {
	b.GameObject = go_
}

func (b *BaseComponent) Start()            {}
func (b *BaseComponent) Update(dt float32) {}

// Drawable is an optional interface for components that render in 3D
type Drawable interface {
	Draw()
}

// UIDrawable is an optional interface for components that render UI (2D overlay)
type UIDrawable interface {
	DrawUI()
}

// GetComponent returns the first component of type T, or nil
func GetComponent[T Component](go_ *GameObject) T {
	var zero T
	for _, c := range go_.components {
		if typed, ok := c.(T); ok {
			return typed
		}
	}
	return zero
}

// HasComponent returns true if GameObject has a component of type T
func HasComponent[T Component](go_ *GameObject) bool {
	for _, c := range go_.components {
		if _, ok := c.(T); ok {
			return true
		}
	}
	return false
}

// GetComponents returns all components of type T on a GameObject
func GetComponents[T Component](go_ *GameObject) []T {
	var result []T
	for _, c := range go_.components {
		if typed, ok := c.(T); ok {
			result = append(result, typed)
		}
	}
	return result
}
