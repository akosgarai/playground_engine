#version 410
layout(location = 0) in vec3 vVertex;
layout(location = 1) in vec3 vNormal;

out vec3 FragPos;
out vec3 Normal;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main()
{
    FragPos = vec3(model * vec4(vVertex, 1.0));
    Normal = mat3(transpose(inverse(model))) * vNormal;
    gl_Position = projection * view * vec4(FragPos,1.0);
}
