# version 410
layout(location = 0) in vec3 vVertex;
layout(location = 1) in vec3 vNormal;
layout(location = 2) in vec2 vTexCoord;

out vec3 SurfacePos;
out vec3 BottomPos;
out vec3 Normal;
out vec2 TexCoords;
out float Depth;

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
    SurfacePos = vec3(model * vec4(Surface, 1.0));
    BottomPos =  vec3(model * vec4(vVertex, 1.0));
    float h = 0.0;
    Depth = waterLevel - vVertex.y;
    // waves are not allowed for underground water
    if (Depth > 0.0) {
        float distance = length(SurfacePos);
        h = amplitude*sin(-PI*distance*frequency+time);
    }

    // vNormal connects to the bottom normal vector.
    // For the surface on the water level, we can expect that
    // it points to the up direction.
    Normal = mat3(transpose(inverse(model))) * vec3(0.0,-1.0,0.0);
    TexCoords = vTexCoord;
    SurfacePos = vec3(SurfacePos.x, SurfacePos.y+h, SurfacePos.z);
    gl_Position = projection * view * vec4(SurfacePos, 1.0);
}
