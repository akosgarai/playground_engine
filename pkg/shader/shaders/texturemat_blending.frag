# version 410
out vec4 FragColor;

struct Tex {
    sampler2D diffuse;
    sampler2D specular;
};
struct Material {
    vec3 ambient;
    vec3 diffuse;
    vec3 specular;
    float shininess;
};

struct DirectionalLight {
    vec3 direction;

    vec3 ambient;
    vec3 diffuse;
    vec3 specular;
};

struct PointLight {
    vec3 position;

    vec3 ambient;
    vec3 diffuse;
    vec3 specular;

    float constant;
    float linear;
    float quadratic;
};

struct SpotLight {
    vec3 position;
    vec3 direction;
    float cutOff;
    float outerCutOff;

    vec3 ambient;
    vec3 diffuse;
    vec3 specular;

    float constant;
    float linear;
    float quadratic;
};

in vec3 FragPos;
in vec3 Normal;
in vec2 TexCoords;

#define MAX_DIRECTION_LIGHTS 16
#define MAX_POINT_LIGHTS 16
#define MAX_SPOT_LIGHTS 16

uniform DirectionalLight dirLight[MAX_DIRECTION_LIGHTS];
uniform PointLight pointLight[MAX_POINT_LIGHTS];
uniform SpotLight spotLight[MAX_SPOT_LIGHTS];
uniform Material material;
uniform Tex tex;
uniform int NumberOfDirectionalLightSources;
uniform int NumberOfPointLightSources;
uniform int NumberOfSpotLightSources;

uniform vec3 viewPosition;

// function prototypes
vec4 CalculateDirectionalLight(DirectionalLight light, vec3 normal, vec3 viewDir);
vec4 CalculatePointLight(PointLight light, vec3 normal, vec3 fragPos, vec3 viewDir);
vec4 CalculateSpotLight(SpotLight light, vec3 normal, vec3 fragPos, vec3 viewDir);

void main()
{
    vec3 norm = normalize(Normal);
    vec3 viewDirection = normalize(viewPosition - FragPos);

    vec4 result = vec4(0);
    // calculate Directional lighting
    int nrDirLight = min(NumberOfDirectionalLightSources, MAX_DIRECTION_LIGHTS);
    for (int i = 0; i < nrDirLight; i++) {
        result += CalculateDirectionalLight(dirLight[i], norm, viewDirection);
    }
    // calculate Point lighting
    int nrPointLight = min(NumberOfPointLightSources, MAX_POINT_LIGHTS);
    for (int i = 0; i < nrPointLight; i++) {
        result += CalculatePointLight(pointLight[i], norm, FragPos, viewDirection);
    }
    // calculate spot lighting
    int nrSpotLight = min(NumberOfSpotLightSources, MAX_SPOT_LIGHTS);
    for (int i = 0; i < nrSpotLight; i++) {
        result += CalculateSpotLight(spotLight[i], norm, FragPos, viewDirection);
    }
    FragColor = result;
}

// calculates the color when using a directional light.
vec4 CalculateDirectionalLight(DirectionalLight light, vec3 normal, vec3 viewDir)
{
    vec3 lightDir = normalize(-light.direction);
    // diffuse shading
    float diff = max(dot(normal, lightDir), 0.0);
    // specular shading
    vec3 reflectDir = reflect(-lightDir, normal);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shininess);
    // combine results
    vec3 ambient = light.ambient * material.ambient * texture(tex.diffuse, TexCoords).rbg;
    vec3 diffuse = light.diffuse * diff * material.diffuse * texture(tex.diffuse, TexCoords).rbg;
    vec3 specular = light.specular * spec * material.specular * texture(tex.specular, TexCoords).rbg;
    return (vec4(ambient,texture(tex.diffuse, TexCoords).a) + vec4(diffuse,texture(tex.diffuse, TexCoords).a) + vec4(specular, texture(tex.specular, TexCoords).a));
}

// calculates the color when using a point light.
vec4 CalculatePointLight(PointLight light, vec3 normal, vec3 fragPos, vec3 viewDir)
{
    vec3 lightDir = normalize(light.position - fragPos);
    // diffuse shading
    float diff = max(dot(normal, lightDir), 0.0);
    // specular shading
    vec3 reflectDir = reflect(-lightDir, normal);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shininess);
    // attenuation
    float distance = length(light.position - fragPos);
    float attenuation = 1.0 / (light.constant + light.linear * distance + light.quadratic * (distance * distance));
    // combine results
    vec3 ambient = light.ambient * material.ambient * texture(tex.diffuse, TexCoords).rbg;
    vec3 diffuse = light.diffuse * material.diffuse * diff * texture(tex.diffuse, TexCoords).rbg;
    vec3 specular = light.specular * material.specular * spec * texture(tex.specular, TexCoords).rbg;
    ambient *= attenuation;
    diffuse *= attenuation;
    specular *= attenuation;
    return (vec4(ambient,texture(tex.diffuse, TexCoords).a) + vec4(diffuse,texture(tex.diffuse, TexCoords).a) + vec4(specular, texture(tex.specular, TexCoords).a));
}

// calculates the color when using a spot light.
vec4 CalculateSpotLight(SpotLight light, vec3 normal, vec3 fragPos, vec3 viewDir)
{
    vec3 lightDir = normalize(light.position - fragPos);
    // diffuse shading
    float diff = max(dot(normal, lightDir), 0.0);
    // specular shading
    vec3 reflectDir = reflect(-lightDir, normal);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shininess);
    // attenuation
    float distance = length(light.position - fragPos);
    float attenuation = 1.0 / (light.constant + light.linear * distance + light.quadratic * (distance * distance));
    // spotlight intensity
    float theta = dot(lightDir, normalize(-light.direction));
    float epsilon = light.cutOff - light.outerCutOff;
    float intensity = clamp((theta - light.outerCutOff) / epsilon, 0.0, 1.0);
    // combine results
    vec3 ambient = light.ambient * material.ambient * texture(tex.diffuse, TexCoords).rbg;
    vec3 diffuse = light.diffuse * diff * material.diffuse * texture(tex.diffuse, TexCoords).rbg;
    vec3 specular = light.specular * spec * material.specular * texture(tex.specular, TexCoords).rbg;
    ambient *= attenuation * intensity;
    diffuse *= attenuation * intensity;
    specular *= attenuation * intensity;
    return (vec4(ambient,texture(tex.diffuse, TexCoords).a) + vec4(diffuse,texture(tex.diffuse, TexCoords).a) + vec4(specular, texture(tex.specular, TexCoords).a));
}
