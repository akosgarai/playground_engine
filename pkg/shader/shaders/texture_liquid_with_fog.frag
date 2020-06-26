# version 410
out vec4 FragColor;

struct Material {
    sampler2D diffuse;
    sampler2D specular;
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
struct Fog {
    float minDistance;
    float maxDistance;
    vec3 color;
};


in vec3 SurfacePos;
in vec3 BottomPos;
in vec3 Normal;
in vec2 TexCoords;
in float Depth;

#define MAX_DIRECTION_LIGHTS 16
#define MAX_POINT_LIGHTS 16
#define MAX_SPOT_LIGHTS 16

uniform DirectionalLight dirLight[MAX_DIRECTION_LIGHTS];
uniform PointLight pointLight[MAX_POINT_LIGHTS];
uniform SpotLight spotLight[MAX_SPOT_LIGHTS];
uniform Material material;
uniform int NumberOfDirectionalLightSources;
uniform int NumberOfPointLightSources;
uniform int NumberOfSpotLightSources;
uniform float Eta;
uniform Fog fog;

uniform vec3 viewPosition;

// function prototypes
vec4 CalculateDirectionalLight(DirectionalLight light, vec3 normal, vec3 viewDir);
vec4 CalculatePointLight(PointLight light, vec3 normal, vec3 fragPos, vec3 viewDir);
vec4 CalculateSpotLight(SpotLight light, vec3 normal, vec3 fragPos, vec3 viewDir);

const float FresnelPower = 5.0;

void main()
{
    vec3 normal = normalize(Normal);
    vec3 viewDirection = normalize(viewPosition - SurfacePos);

    vec4 resultReflect = vec4(0);
    vec4 resultRefract = vec4(0);
    vec4 resultView = vec4(0);
    float eta = Eta;
    vec3 norm = normal;
    if (viewPosition.y < SurfacePos.y) {
        eta = 1.0 / Eta;
    }


    vec3 reflectDir = reflect(viewDirection, norm);
    vec3 refractDir = refract(viewDirection, norm, eta);

    // calculate Directional lighting
    int nrDirLight = min(NumberOfDirectionalLightSources, MAX_DIRECTION_LIGHTS);
    for (int i = 0; i < nrDirLight; i++) {
        resultView += CalculateDirectionalLight(dirLight[i], norm, viewDirection);
        resultReflect += CalculateDirectionalLight(dirLight[i], norm, reflectDir);
        resultRefract += CalculateDirectionalLight(dirLight[i], norm, refractDir);
    }
    // calculate Point lighting
    int nrPointLight = min(NumberOfPointLightSources, MAX_POINT_LIGHTS);
    for (int i = 0; i < nrPointLight; i++) {
        resultView += CalculatePointLight(pointLight[i], norm, SurfacePos, viewDirection);
        resultReflect += CalculatePointLight(pointLight[i], norm, BottomPos, reflectDir);
        resultRefract += CalculatePointLight(pointLight[i], norm, BottomPos, refractDir);
    }
    // calculate spot lighting
    int nrSpotLight = min(NumberOfSpotLightSources, MAX_SPOT_LIGHTS);
    for (int i = 0; i < nrSpotLight; i++) {
        resultView += CalculateSpotLight(spotLight[i], norm, SurfacePos, viewDirection);
        resultReflect += CalculateSpotLight(spotLight[i], norm, BottomPos, reflectDir);
        resultRefract += CalculateSpotLight(spotLight[i], norm, BottomPos, refractDir);
    }
    FragColor = resultView;
    if (Depth > 0.0) {
        float F = ((1.0-eta) * (1.0-eta)) / ((1.0+eta) * (1.0+eta));
        float Ratio = F + (1.0 - F) * pow((1.0 - dot(-viewDirection, norm)), FresnelPower);
        vec4 mixed = mix(resultRefract, resultReflect, Ratio);
        FragColor = vec4(vec3(mixed), mixed.w/(resultView.w+mixed.w));
    }
    if (fog.minDistance > 0 && fog.maxDistance > 0) {
        float distance = length(viewPosition - SurfacePos);
        float fogFactor = (fog.maxDistance - distance) / (fog.maxDistance - fog.minDistance);
        fogFactor = clamp(fogFactor, 0.0, 1.0);
        FragColor = vec4(mix(fog.color, FragColor.xyz, fogFactor), FragColor.w);
    }

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
    vec3 ambient = light.ambient * texture(material.diffuse, TexCoords).rgb;
    vec3 diffuse = light.diffuse * diff * texture(material.diffuse, TexCoords).rgb;
    vec3 specular = light.specular * spec * texture(material.specular, TexCoords).rgb;
    return (vec4(ambient,texture(material.diffuse, TexCoords).a) + vec4(diffuse,texture(material.diffuse, TexCoords).a) + vec4(specular, texture(material.specular, TexCoords).a));
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
    vec3 ambient = light.ambient * texture(material.diffuse, TexCoords).rbg;
    vec3 diffuse = light.diffuse * diff * texture(material.diffuse, TexCoords).rbg;
    vec3 specular = light.specular * spec * texture(material.specular, TexCoords).rbg;
    ambient *= attenuation;
    diffuse *= attenuation;
    specular *= attenuation;
    return (vec4(ambient,texture(material.diffuse, TexCoords).a) + vec4(diffuse,texture(material.diffuse, TexCoords).a) + vec4(specular, texture(material.specular, TexCoords).a));
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
    vec3 ambient = light.ambient * texture(material.diffuse, TexCoords).rgb;
    vec3 diffuse = light.diffuse * diff * texture(material.diffuse, TexCoords).rgb;
    vec3 specular = light.specular * spec * texture(material.specular, TexCoords).rgb;
    ambient *= attenuation * intensity;
    diffuse *= attenuation * intensity;
    specular *= attenuation * intensity;
    return (vec4(ambient,texture(material.diffuse, TexCoords).a) + vec4(diffuse,texture(material.diffuse, TexCoords).a) + vec4(specular, texture(material.specular, TexCoords).a));
}
