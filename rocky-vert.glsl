#version 130

attribute vec4 Vertex;
attribute vec4 Color;
attribute vec2 TexCoord0;

varying vec2 tc0;

out vec3 FragPos;
out float color;
out float tiles0;

uniform mat4 Projection;
uniform mat4 MVP;

uniform sampler2D Texture0;

void main()
{
	FragPos = vec3(Vertex);
	tc0 = TexCoord0;
	tiles0 = 1;
	gl_FrontColor = Color;
	
	vec4 t0 = texture2D(Texture0, tc0 * tiles0);
	float offset = normalize(t0).x;
	color = offset;
	
	vec3 center = vec3(0, 0, 0);
	vec3 pos = vec3(Vertex);
	vec3 verts = pos + (offset * (pos - center));

	// render deformed
	gl_Position = MVP * vec4(verts, 1.0);
}
