#version 130

in vec3 FragPos;
in float color;

varying vec2 tc0;
varying float tiles0;
uniform sampler2D Texture0;

void main()
{
	vec3 normal = vec3(1, 1, 1);
	vec3 pos = vec3(0, 2, 0);
	vec3 surfaceToLight = pos - FragPos;
	float brightness = dot(normal, surfaceToLight) / (length(surfaceToLight) * length(normal));
	brightness = clamp(brightness, 0, 1) * 4.5;
	vec3 light = brightness * vec3(0.9, 0.8, 0.7);

	vec4 final = vec4(0, 0, 0, 0);
	vec4 t0 = texture2D(Texture0, tc0 * tiles0);
	final = mix(final, t0, t0.a);
	
	gl_FragColor = vec4(light, 1.0) * color * final;
}
