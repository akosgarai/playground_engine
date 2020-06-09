# version 410
out vec4 FragColor;

struct Tex {
    sampler2D diffuse;
    sampler2D specular;
};

uniform Tex tex;

in vec3 AmbientColor;
in vec2 FragPos;

void main()
{
    vec3 ambient = texture(tex.diffuse, FragPos).rgb * AmbientColor;
    vec3 diffuse = texture(tex.diffuse, FragPos).rgb;
    vec3 specular = texture(tex.specular, FragPos).rgb;
    FragColor = vec4(ambient + diffuse + specular, 1.0);
}
