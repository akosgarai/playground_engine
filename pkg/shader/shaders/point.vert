#version 410
layout(location = 0) in vec3 vVertex;
layout(location = 1) in vec3 vColor;
layout(location = 2) in float vSize;

smooth out vec3 vSmoothColor;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main()
{
    vSmoothColor = vColor;
    gl_Position = projection * view * model * vec4(vVertex,1);
    gl_PointSize = vSize;
}
