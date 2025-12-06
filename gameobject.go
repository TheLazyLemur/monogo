package monogo

// GameObject is the base entity in the game world
type GameObject struct {
	Name       string
	Tag        string
	Transform  *Transform
	components []Component
	drawables  []Drawable   // cached for fast iteration
	uiDrawables []UIDrawable // cached for fast iteration
	Active     bool
	scene      *Scene
	parent     *GameObject
	children   []*GameObject
	persistent bool // survives scene loads (DontDestroyOnLoad)
}

// NewGameObject creates a new GameObject with a name
func NewGameObject(name string) *GameObject {
	go_ := &GameObject{
		Name:       name,
		Tag:        "Untagged",
		Transform:  NewTransform(),
		components: make([]Component, 0),
		Active:     true,
		children:   make([]*GameObject, 0),
	}
	go_.Transform.gameObject = go_
	return go_
}

// AddComponent attaches a component to the GameObject
func (g *GameObject) AddComponent(c Component) {
	c.SetGameObject(g)
	g.components = append(g.components, c)
	// Cache drawable interfaces for fast iteration
	if d, ok := c.(Drawable); ok {
		g.drawables = append(g.drawables, d)
	}
	if u, ok := c.(UIDrawable); ok {
		g.uiDrawables = append(g.uiDrawables, u)
	}
}

// Add attaches a component and returns the GameObject for chaining
func (g *GameObject) Add(c Component) *GameObject {
	g.AddComponent(c)
	return g
}

// GetAllComponents returns all components (use GetComponent[T] for typed access)
func (g *GameObject) GetAllComponents() []Component {
	return g.components
}

// Scene returns the scene this GameObject belongs to
func (g *GameObject) Scene() *Scene {
	return g.scene
}

// SetActive enables or disables the GameObject
func (g *GameObject) SetActive(active bool) {
	g.Active = active
}

// Parent returns the parent GameObject, or nil if none
func (g *GameObject) Parent() *GameObject {
	return g.parent
}

// Children returns all child GameObjects
func (g *GameObject) Children() []*GameObject {
	return g.children
}

// SetParent sets the parent of this GameObject
func (g *GameObject) SetParent(parent *GameObject) {
	// Remove from old parent
	if g.parent != nil {
		g.parent.removeChild(g)
	}

	g.parent = parent

	// Add to new parent
	if parent != nil {
		parent.children = append(parent.children, g)
	}
}

// AddChild adds a child GameObject
func (g *GameObject) AddChild(child *GameObject) {
	child.SetParent(g)
}

// RemoveChild removes a child GameObject
func (g *GameObject) RemoveChild(child *GameObject) {
	if child.parent == g {
		child.SetParent(nil)
	}
}

func (g *GameObject) removeChild(child *GameObject) {
	for i, c := range g.children {
		if c == child {
			g.children = append(g.children[:i], g.children[i+1:]...)
			return
		}
	}
}

// Start calls Start on all components
func (g *GameObject) Start() {
	for _, c := range g.components {
		c.Start()
	}
}

// Update calls Update on all components with delta time
func (g *GameObject) Update(dt float32) {
	if !g.Active {
		return
	}
	for _, c := range g.components {
		c.Update(dt)
	}
}
