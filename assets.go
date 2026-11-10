package main

import (
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const vsLit = `#version 330
in vec3 vertexPosition;
in vec2 vertexTexCoord;
in vec3 vertexNormal;
in vec4 vertexColor;
uniform mat4 mvp;
uniform mat4 matModel;
uniform mat4 matNormal;
out vec3 fragPos;
out vec2 fragTexCoord;
out vec3 fragNormal;
out vec4 fragColor;
void main() {
	fragPos = vec3(matModel * vec4(vertexPosition, 1.0));
	fragTexCoord = vertexTexCoord;
	fragColor = vertexColor;
	fragNormal = normalize(vec3(matNormal * vec4(vertexNormal, 0.0)));
	gl_Position = mvp * vec4(vertexPosition, 1.0);
}`

const fsLit = `#version 330
in vec3 fragPos;
in vec2 fragTexCoord;
in vec3 fragNormal;
in vec4 fragColor;
uniform sampler2D texture0;
uniform vec4 colDiffuse;
uniform vec3 lightDir;
uniform vec3 viewPos;
uniform vec3 fogColor;
uniform float fogDensity;
out vec4 finalColor;
void main() {
	vec4 tex = texture(texture0, fragTexCoord) * colDiffuse * fragColor;
	float d = max(dot(normalize(fragNormal), -lightDir), 0.0);
	float sky = 0.5 + 0.5 * normalize(fragNormal).y;
	vec3 c = tex.rgb * (0.35 + 0.25 * sky + 0.65 * d);
	float dist = length(viewPos - fragPos);
	float f = 1.0 - exp(-pow(dist * fogDensity, 2.0));
	c = mix(c, fogColor, clamp(f, 0.0, 1.0));
	finalColor = vec4(c, tex.a);
}`

type Assets struct {
	Dir     string
	Models  map[string]rl.Model
	Shader  rl.Shader
	locView int32
	Font    rl.Font
	Sounds  map[string]rl.Sound
	Engine  rl.Music
	Audio   bool
}

var fogColor = rl.NewColor(178, 222, 245, 255)

func LoadAssets(dir string, audio bool) *Assets {
	a := &Assets{Dir: dir, Models: map[string]rl.Model{}, Sounds: map[string]rl.Sound{}, Audio: audio}
	a.Shader = rl.LoadShaderFromMemory(vsLit, fsLit)
	a.Shader.UpdateLocation(rl.ShaderLocMatrixModel, rl.GetShaderLocation(a.Shader, "matModel"))
	a.Shader.UpdateLocation(rl.ShaderLocMatrixNormal, rl.GetShaderLocation(a.Shader, "matNormal"))
	a.locView = rl.GetShaderLocation(a.Shader, "viewPos")
	l := rl.Vector3Normalize(rl.NewVector3(-0.4, -1, -0.6))
	rl.SetShaderValue(a.Shader, rl.GetShaderLocation(a.Shader, "lightDir"), []float32{l.X, l.Y, l.Z}, rl.ShaderUniformVec3)
	rl.SetShaderValue(a.Shader, rl.GetShaderLocation(a.Shader, "fogColor"),
		[]float32{float32(fogColor.R) / 255, float32(fogColor.G) / 255, float32(fogColor.B) / 255}, rl.ShaderUniformVec3)
	rl.SetShaderValue(a.Shader, rl.GetShaderLocation(a.Shader, "fogDensity"), []float32{0.012}, rl.ShaderUniformFloat)

	models := map[string][]string{
		"city": {"building-garage", "building-small-a", "building-small-b", "building-small-c", "building-small-d",
			"grass-trees-tall", "grass-trees", "grass", "pavement-fountain", "pavement",
			"road-corner", "road-intersection", "road-split", "road-straight-lightposts", "road-straight"},
		"racing":     {"vehicle-truck-red", "vehicle-truck-yellow", "vehicle-truck-green", "vehicle-truck-purple"},
		"platformer": {"coin", "flag", "cloud"},
	}
	for kit, names := range models {
		for _, n := range names {
			m := rl.LoadModel(filepath.Join(dir, kit, n+".glb"))
			mats := m.GetMaterials()
			for i := range mats {
				mats[i].Shader = a.Shader
				tex := mats[i].GetMap(rl.MapDiffuse).Texture
				rl.SetTextureFilter(tex, rl.FilterBilinear)
			}
			a.Models[n] = m
		}
	}
	ball := rl.LoadModelFromMesh(rl.GenMeshSphere(1, 12, 16))
	ball.GetMaterials()[0].Shader = a.Shader
	a.Models["ball"] = ball

	a.Font = rl.LoadFontEx(filepath.Join(dir, "fonts", "lilita_one_regular.ttf"), 96, FontRunes())
	rl.SetTextureFilter(a.Font.Texture, rl.FilterBilinear)

	if audio {
		for _, n := range []string{"coin", "impact", "placement-a", "toggle", "rotate", "skid"} {
			a.Sounds[n] = rl.LoadSound(filepath.Join(dir, "sounds", n+".ogg"))
		}
		a.Engine = rl.LoadMusicStream(filepath.Join(dir, "sounds", "engine.ogg"))
		a.Engine.Looping = true
	}
	return a
}

func (a *Assets) SetView(cam rl.Camera3D) {
	rl.SetShaderValue(a.Shader, a.locView, []float32{cam.Position.X, cam.Position.Y, cam.Position.Z}, rl.ShaderUniformVec3)
}

func (a *Assets) Play(name string) {
	if !a.Audio {
		return
	}
	if s, ok := a.Sounds[name]; ok {
		rl.PlaySound(s)
	}
}

// Draw draws a model at pos rotated by rotDeg around Y, uniformly scaled.
func (a *Assets) Draw(name string, pos rl.Vector3, rotDeg, scale float32) {
	rl.DrawModelEx(a.Models[name], pos, rl.NewVector3(0, 1, 0), rotDeg, rl.NewVector3(scale, scale, scale), rl.White)
}

// Ball draws a unit sphere squashed to size, moved by off, then placed by the matrix place.
func (a *Assets) Ball(off, size rl.Vector3, col rl.Color, place rl.Matrix) {
	m := a.Models["ball"]
	m.Transform = rl.MatrixMultiply(rl.MatrixMultiply(rl.MatrixScale(size.X, size.Y, size.Z),
		rl.MatrixTranslate(off.X, off.Y, off.Z)), place)
	rl.DrawModel(m, rl.Vector3{}, 1, col)
}

// FontRunes is the glyph set loaded from the font: printable ASCII plus the
// German and French letters used in the UI and in player names.
func FontRunes() []rune {
	var rs []rune
	for r := rune(32); r < 127; r++ {
		rs = append(rs, r)
	}
	return append(rs, []rune("äöüÄÖÜéèàçÉÈÀÇ–«»")...)
}
