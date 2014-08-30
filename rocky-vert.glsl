#version 130

attribute vec4 Vertex;
uniform mat4 MVP;

attribute vec2 TexCoord0;
varying vec2 tc0;

uniform sampler2D Texture0;

uniform float time;

out float tiles0;

void main() {
	tiles0 = 1;
	tc0 = TexCoord0;

	// sample texture to deform mesh
	vec4 t0 = texture2D(Texture0, tc0 * tiles0);
	float offset = normalize(t0).x;
	// use offset as cheap ambient occlusion
	float ao = offset;
	// motion sickness
	// offset += cos(time) * sin(time);
	
	vec3 center = vec3(0, 0, 0);
	vec3 fpos = vec3(Vertex);
	vec3 verts = fpos + (offset * (fpos - center));

	// normal
	// gl_Position = MVP * Vertex;

	// deformed
	gl_Position = MVP * vec4(verts, 1.0);

	// lighting
	vec3 vertexPos = vec3(gl_Position);
	vec3 normalDirection = normalize(gl_NormalMatrix * gl_Normal);
	vec3 lightPos = vec3(0, 10, 2);
	vec3 lightDirection = normalize(lightPos - vertexPos);

	//
	vec4 sceneLight = vec4(0.8, 0.7, 0.6, 1.0);
	
	//
	vec4 amb = vec4(1, 1, 1, 1) * ao;

	//
	vec4 diffuse = vec4(1, 1, 1, 1);
	diffuse = diffuse * max(dot(lightDirection, normalDirection), 0.0);
  diffuse = clamp(diffuse, 0.0, 1.0);     

	//
	vec3 viewDirection = normalize(-vertexPos);
	vec3 reflection = normalize(reflect(-lightDirection, normalDirection));
	float shininess = 2;

	//
	vec4 lightSpec = vec4(1, 1, 1, 1);
	vec4 spec = lightSpec * pow(max(0.0, dot(reflection, viewDirection)), shininess);
	
	gl_FrontColor = (sceneLight + spec + diffuse + amb) * ao;
}
