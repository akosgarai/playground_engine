# version 410
layout(location = 0) in vec3 vVertex;
layout(location = 1) in vec3 vNormal;
layout(location = 2) in vec2 vTexCoord;

out vec3 FragPos;
out vec3 Normal;
out vec2 TexCoords;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
uniform float time;
uniform vec3 viewPosition;
uniform float amplitude;
uniform float frequency;
uniform float waterLevel;

const float PI = 3.14159;

void main()
{
    vec3 Surface = vec3(vVertex.x, waterLevel, vVertex.z);
    FragPos = vec3(model * vec4(Surface, 1.0));
    float distance = length(FragPos);
    float h = amplitude*sin(-PI*distance*frequency+time);

    Normal = mat3(transpose(inverse(model))) * vNormal;
    TexCoords = vTexCoord;
    FragPos = vec3(FragPos.x, FragPos.y+h, FragPos.z);
    gl_Position = projection * view * vec4(FragPos, 1.0);
}
