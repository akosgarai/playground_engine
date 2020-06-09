#version 410
smooth in vec3 vSmoothColor;

layout(location=0) out vec4 FragColor;

void main()
{
    FragColor = vec4(vSmoothColor, 1.0);
}
