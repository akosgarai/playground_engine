# version 410
out vec4 FragColor;
  
struct Material {
    vec3 ambient;
    vec3 specular;
    vec3 diffuse;
    vec3 shininess;
};
in vec2 TexCoord;

uniform sampler2D paper;
uniform Material material;

void main()
{
    vec4 textureComponent = vec4(vec3(texture(paper, TexCoord)), 1.0);
    vec4 ambientComponent = vec4(material.ambient, 1.0) * textureComponent;
    vec4 diffuseComponent = vec4(material.diffuse, 1.0) * textureComponent;
    vec4 specularComponent = vec4(material.specular, 1.0) * textureComponent;
    FragColor = ambientComponent + diffuseComponent + specularComponent;
}
