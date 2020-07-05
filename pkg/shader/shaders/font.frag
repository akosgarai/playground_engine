#version 410
in vec2 TexCoords;
in vec3 Color;

layout(location=0) out vec4 FragColor;

uniform sampler2D tex;

void main()
{
    vec4 sampled = vec4(1.0, 1.0, 1.0, texture(tex, TexCoords).r);
    FragColor = vec4(Color,1.0) * sampled;
}
