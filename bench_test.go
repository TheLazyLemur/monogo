package monogo

import (
	"testing"
)

// Mock component for benchmarking
type benchComponent struct {
	BaseComponent
	counter int
}

func (b *benchComponent) Update(dt float32) {
	b.counter++
}

func (b *benchComponent) Draw() {}

// setupBenchSceneRaw creates scene without raylib dependencies
func setupBenchScene(numObjects, componentsPerObject int) *Scene {
	scene := &Scene{
		gameObjects:    make([]*GameObject, 0, numObjects),
		toDestroy:      make([]*GameObject, 0),
		pendingObjects: make([]*GameObject, 0),
	}
	for i := 0; i < numObjects; i++ {
		obj := &GameObject{
			Name:        "Object",
			Transform:   &Transform{},
			components:  make([]Component, 0, componentsPerObject),
			drawables:   make([]Drawable, 0, componentsPerObject),
			uiDrawables: make([]UIDrawable, 0),
			Active:      true,
		}
		obj.Transform.gameObject = obj
		obj.scene = scene
		for j := 0; j < componentsPerObject; j++ {
			c := &benchComponent{}
			c.GameObject = obj
			obj.components = append(obj.components, c)
			obj.drawables = append(obj.drawables, c) // cache drawable
		}
		scene.gameObjects = append(scene.gameObjects, obj)
	}
	return scene
}

func setupBenchSceneWithHierarchy(numObjects int) *Scene {
	scene := &Scene{
		gameObjects:    make([]*GameObject, 0, numObjects),
		toDestroy:      make([]*GameObject, 0),
		pendingObjects: make([]*GameObject, 0),
	}
	var parent *GameObject
	for i := 0; i < numObjects; i++ {
		obj := &GameObject{
			Name:       "Object",
			Transform:  &Transform{Scale: Vector3{X: 1, Y: 1, Z: 1}},
			components: make([]Component, 0, 1),
			Active:     true,
			children:   make([]*GameObject, 0),
		}
		obj.Transform.gameObject = obj
		obj.scene = scene
		c := &benchComponent{}
		c.GameObject = obj
		obj.components = append(obj.components, c)
		if parent != nil && i%3 == 0 {
			obj.SetParent(parent)
		}
		parent = obj
		scene.gameObjects = append(scene.gameObjects, obj)
	}
	return scene
}

func BenchmarkSceneUpdate(b *testing.B) {
	scene := setupBenchScene(100, 3)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		scene.Update(0.016)
	}
}

func BenchmarkSceneDraw(b *testing.B) {
	scene := setupBenchScene(100, 3)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		scene.Draw()
	}
}

func BenchmarkSceneDrawUI(b *testing.B) {
	scene := setupBenchScene(100, 3)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		scene.DrawUI()
	}
}

func BenchmarkWorldPosition(b *testing.B) {
	scene := setupBenchSceneWithHierarchy(10)
	// Get deepest child
	obj := scene.gameObjects[len(scene.gameObjects)-1]
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = obj.Transform.WorldPosition()
	}
}

func BenchmarkWorldRotation(b *testing.B) {
	scene := setupBenchSceneWithHierarchy(10)
	obj := scene.gameObjects[len(scene.gameObjects)-1]
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = obj.Transform.WorldRotation()
	}
}

func BenchmarkWorldScale(b *testing.B) {
	scene := setupBenchSceneWithHierarchy(10)
	obj := scene.gameObjects[len(scene.gameObjects)-1]
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = obj.Transform.WorldScale()
	}
}

func BenchmarkGetComponent(b *testing.B) {
	scene := setupBenchScene(1, 5)
	obj := scene.gameObjects[0]
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = GetComponent[*benchComponent](obj)
	}
}

func BenchmarkSceneFind(b *testing.B) {
	scene := setupBenchScene(100, 1)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = scene.Find("Object")
	}
}

func BenchmarkProcessDestroyQueue_Empty(b *testing.B) {
	scene := setupBenchScene(100, 1)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		scene.processDestroyQueue()
	}
}

func BenchmarkProcessDestroyQueue_WithDestroys(b *testing.B) {
	// Pre-allocate objects to reuse
	objects := make([]*GameObject, 100)
	for i := range objects {
		objects[i] = &GameObject{Name: "X", Active: true}
	}
	scene := &Scene{
		gameObjects:    make([]*GameObject, 0, 100),
		toDestroy:      make([]*GameObject, 0, 10),
		pendingObjects: make([]*GameObject, 0),
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Reset scene with pre-allocated objects
		scene.gameObjects = scene.gameObjects[:0]
		scene.gameObjects = append(scene.gameObjects, objects...)
		scene.toDestroy = scene.toDestroy[:0]
		scene.toDestroy = append(scene.toDestroy, objects[50])
		scene.processDestroyQueue()
	}
}
