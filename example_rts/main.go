package main

import "github.com/TheLazyLemur/monogo"

// HologramVS - vertex shader that outputs normals and position for hologram effect
const HologramVS = `
#version 330
in vec3 vertexPosition;
in vec2 vertexTexCoord;
in vec3 vertexNormal;
in vec4 vertexColor;

uniform mat4 mvp;
uniform mat4 matModel;
uniform mat4 matNormal;

out vec2 fragTexCoord;
out vec4 fragColor;
out vec3 fragNormal;
out vec3 fragPosition;

void main() {
    fragTexCoord = vertexTexCoord;
    fragColor = vertexColor;
    fragNormal = normalize(vec3(matNormal * vec4(vertexNormal, 0.0)));
    fragPosition = vec3(matModel * vec4(vertexPosition, 1.0));
    gl_Position = mvp * vec4(vertexPosition, 1.0);
}
`

// HologramFS - sci-fi hologram effect
const HologramFS = `
#version 330
in vec2 fragTexCoord;
in vec4 fragColor;
in vec3 fragNormal;
in vec3 fragPosition;
uniform sampler2D texture0;
uniform vec4 colDiffuse;
uniform float time;
out vec4 finalColor;

void main() {
    vec3 normal = normalize(fragNormal);

    // Fresnel edge glow
    vec3 viewDir = normalize(-fragPosition);
    float fresnel = pow(1.0 - max(dot(normal, viewDir), 0.0), 2.0);

    // Scan lines
    float scanLine = sin(fragPosition.y * 50.0 + time * 5.0) * 0.5 + 0.5;
    scanLine = pow(scanLine, 0.5);

    // Flicker
    float flicker = sin(time * 30.0) * 0.05 + 0.95;

    // Glitch offset (occasional)
    float glitch = step(0.98, sin(time * 2.0)) * sin(time * 100.0) * 0.1;

    // Base hologram color (cyan)
    vec3 holoColor = vec3(0.0, 0.8, 1.0);

    // Mix with diffuse color
    vec3 baseColor = colDiffuse.rgb * fragColor.rgb;
    vec3 color = mix(baseColor, holoColor, 0.7);

    // Apply effects
    float alpha = (0.3 + fresnel * 0.7) * scanLine * flicker;
    alpha = clamp(alpha + glitch, 0.0, 1.0);

    // Edge highlight
    color += holoColor * fresnel * 0.5;

    finalColor = vec4(color, alpha * 0.8);
}
`

type RTSGame struct {
	holoShader *monogo.Shader
}

func (g *RTSGame) Setup(scene *monogo.Scene) {
	// Load hologram shader for units (needs custom VS for fragNormal/fragPosition)
	g.holoShader = monogo.LoadShaderFromMemory(HologramVS, HologramFS)

	// RTS Camera
	camObj := scene.Create("RTSCamera")
	camObj.Add(NewRTSCamera())

	// Selection manager
	selMgr := scene.Create("SelectionManager")
	selMgr.Add(NewSelectionManager())

	// Spawn some units
	positions := []monogo.Vector3{
		monogo.Vec3(-5, 0.5, -5),
		monogo.Vec3(-3, 0.5, -5),
		monogo.Vec3(-4, 0.5, -3),
		monogo.Vec3(5, 0.5, 5),
		monogo.Vec3(7, 0.5, 5),
		monogo.Vec3(6, 0.5, 7),
	}
	colors := []monogo.Color{
		monogo.Blue, monogo.Blue, monogo.Blue,
		monogo.Red, monogo.Red, monogo.Red,
	}

	for i, pos := range positions {
		unit := scene.Create("Unit")
		unit.Tag = "Unit"
		unit.Transform.Position = pos
		unit.Add(monogo.NewCubeRenderer(colors[i], monogo.Vec3(0.8, 1, 0.8)).SetShader(g.holoShader))
		unit.Add(NewUnit())
	}

	// Ground
	ground := scene.Create("Ground")
	ground.Transform.Position = monogo.Vec3(0, -0.1, 0)
	ground.Transform.Scale = monogo.Vec3(50, 0.2, 50)
	ground.Add(monogo.NewCubeRenderer(monogo.DarkGray, monogo.Vec3(1, 1, 1)))

	// UI
	ui := scene.Create("UI")
	ui.Add(NewRTSUI())

	// Post-processing chain: vignette + chromatic aberration
	scene.PostProcess = monogo.NewPostProcessChain()

	// Vignette effect
	vignette := monogo.NewPostProcess(
		monogo.LoadShaderFromMemory("", monogo.VignetteFS),
	)
	vignette.GetShader().SetFloat("intensity", 0.5)
	scene.PostProcess.Add(vignette)

	// Chromatic aberration (subtle)
	chromatic := monogo.NewPostProcess(
		monogo.LoadShaderFromMemory("", monogo.ChromaticAberrationFS),
	)
	chromatic.GetShader().SetFloat("offset", 0.002)
	scene.PostProcess.Add(chromatic)
}

func (g *RTSGame) Update(scene *monogo.Scene, dt float32) {
	// Update time uniform for hologram shader animation
	g.holoShader.SetFloat("time", float32(monogo.Time()))
}

func main() {
	monogo.Run(&RTSGame{}, monogo.Config{
		Title:    "MonoGo RTS Example",
		ShowFPS:  true,
		ShowGrid: true,
	})
}
