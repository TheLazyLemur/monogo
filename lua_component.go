package monogo

import (
	lua "github.com/yuin/gopher-lua"
)

// LuaComponent runs a Lua script as a component
type LuaComponent struct {
	BaseComponent
	state    *lua.LState
	path     string
	hasStart bool
}

// NewLuaComponent creates a component from a Lua file
func NewLuaComponent(path string) *LuaComponent {
	return &LuaComponent{path: path}
}

func (l *LuaComponent) Start() {
	l.state = lua.NewState()

	// Register API
	l.registerAPI()

	// Load script
	if err := l.state.DoFile(l.path); err != nil {
		println("Lua error:", err.Error())
		return
	}

	// Check if start function exists
	if l.state.GetGlobal("start").Type() == lua.LTFunction {
		l.hasStart = true
		if err := l.state.CallByParam(lua.P{
			Fn:      l.state.GetGlobal("start"),
			NRet:    0,
			Protect: true,
		}); err != nil {
			println("Lua start error:", err.Error())
		}
	}
}

func (l *LuaComponent) Update(dt float32) {
	if l.state == nil {
		return
	}

	// Update globals
	l.state.SetGlobal("dt", lua.LNumber(dt))

	// Call update if exists
	fn := l.state.GetGlobal("update")
	if fn.Type() != lua.LTFunction {
		return
	}

	if err := l.state.CallByParam(lua.P{
		Fn:      fn,
		NRet:    0,
		Protect: true,
	}, lua.LNumber(dt)); err != nil {
		println("Lua update error:", err.Error())
	}
}

func (l *LuaComponent) registerAPI() {
	// Transform table
	l.state.SetGlobal("transform", l.makeTransformTable())

	// Input functions
	l.state.SetGlobal("key_pressed", l.state.NewFunction(l.luaKeyPressed))
	l.state.SetGlobal("key_held", l.state.NewFunction(l.luaKeyHeld))
	l.state.SetGlobal("horizontal", l.state.NewFunction(l.luaHorizontal))
	l.state.SetGlobal("vertical", l.state.NewFunction(l.luaVertical))

	// Scene functions
	l.state.SetGlobal("load_scene", l.state.NewFunction(l.luaLoadScene))
	l.state.SetGlobal("find", l.state.NewFunction(l.luaFind))

	// Key constants
	keys := l.state.NewTable()
	l.state.SetField(keys, "space", lua.LNumber(KeySpace))
	l.state.SetField(keys, "enter", lua.LNumber(KeyEnter))
	l.state.SetField(keys, "escape", lua.LNumber(KeyEscape))
	l.state.SetField(keys, "shift", lua.LNumber(KeyShift))
	l.state.SetField(keys, "w", lua.LNumber(KeyW))
	l.state.SetField(keys, "a", lua.LNumber(KeyA))
	l.state.SetField(keys, "s", lua.LNumber(KeyS))
	l.state.SetField(keys, "d", lua.LNumber(KeyD))
	l.state.SetGlobal("keys", keys)
}

func (l *LuaComponent) makeTransformTable() *lua.LTable {
	t := l.state.NewTable()

	// Position getters/setters
	l.state.SetField(t, "get_x", l.state.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LNumber(l.GameObject.Transform.Position.X))
		return 1
	}))
	l.state.SetField(t, "get_y", l.state.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LNumber(l.GameObject.Transform.Position.Y))
		return 1
	}))
	l.state.SetField(t, "get_z", l.state.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LNumber(l.GameObject.Transform.Position.Z))
		return 1
	}))
	l.state.SetField(t, "set_x", l.state.NewFunction(func(L *lua.LState) int {
		l.GameObject.Transform.Position.X = float32(L.CheckNumber(1))
		return 0
	}))
	l.state.SetField(t, "set_y", l.state.NewFunction(func(L *lua.LState) int {
		l.GameObject.Transform.Position.Y = float32(L.CheckNumber(1))
		return 0
	}))
	l.state.SetField(t, "set_z", l.state.NewFunction(func(L *lua.LState) int {
		l.GameObject.Transform.Position.Z = float32(L.CheckNumber(1))
		return 0
	}))

	// Rotation
	l.state.SetField(t, "get_rot_y", l.state.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LNumber(l.GameObject.Transform.Rotation.Y))
		return 1
	}))
	l.state.SetField(t, "set_rot_y", l.state.NewFunction(func(L *lua.LState) int {
		l.GameObject.Transform.Rotation.Y = float32(L.CheckNumber(1))
		return 0
	}))

	// Move helper
	l.state.SetField(t, "move", l.state.NewFunction(func(L *lua.LState) int {
		dx := float32(L.CheckNumber(1))
		dy := float32(L.CheckNumber(2))
		dz := float32(L.CheckNumber(3))
		l.GameObject.Transform.Position.X += dx
		l.GameObject.Transform.Position.Y += dy
		l.GameObject.Transform.Position.Z += dz
		return 0
	}))

	return t
}

// Input bindings
func (l *LuaComponent) luaKeyPressed(L *lua.LState) int {
	key := L.CheckInt(1)
	L.Push(lua.LBool(KeyPressed(int32(key))))
	return 1
}

func (l *LuaComponent) luaKeyHeld(L *lua.LState) int {
	key := L.CheckInt(1)
	L.Push(lua.LBool(KeyHeld(int32(key))))
	return 1
}

func (l *LuaComponent) luaHorizontal(L *lua.LState) int {
	L.Push(lua.LNumber(Horizontal()))
	return 1
}

func (l *LuaComponent) luaVertical(L *lua.LState) int {
	L.Push(lua.LNumber(Vertical()))
	return 1
}

// Scene bindings
func (l *LuaComponent) luaLoadScene(L *lua.LState) int {
	name := L.CheckString(1)
	LoadScene(name)
	return 0
}

func (l *LuaComponent) luaFind(L *lua.LState) int {
	name := L.CheckString(1)
	obj := l.GameObject.Scene().Find(name)
	if obj == nil {
		L.Push(lua.LNil)
		return 1
	}
	// Return a table with the found object's transform
	t := L.NewTable()
	L.SetField(t, "get_x", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LNumber(obj.Transform.Position.X))
		return 1
	}))
	L.SetField(t, "get_y", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LNumber(obj.Transform.Position.Y))
		return 1
	}))
	L.SetField(t, "get_z", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LNumber(obj.Transform.Position.Z))
		return 1
	}))
	L.Push(t)
	return 1
}

// Cleanup
func (l *LuaComponent) OnDestroy() {
	if l.state != nil {
		l.state.Close()
	}
}
