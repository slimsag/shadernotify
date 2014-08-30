#version 130

uniform float time;

varying vec2 tc0;
uniform sampler2D Texture0;

varying float tiles0;

void main() {
	vec4 final = vec4(0, 0, 0, 0);
	vec4 t0 = texture2D(Texture0, tc0 * tiles0);
	final = mix(final, t0, t0.a);

  gl_FragColor = gl_Color * final;
}
