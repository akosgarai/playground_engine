# Materials

This package tries to implement the base needs of the following tutorial: https://learnopengl.com/Lighting/Materials

## Defined materials.

This [page](http://devernay.free.fr/cours/opengl/materials.html) contains the setup for materials. I made the named materials based on this document.

## The math behind it

First we have to know the variables. The goal is to calculate the color of a given object made by a given material when we have a given light source.

- Ambient

The ambient color component.

- Diffuse

The diffuse color component.

- Spectral

The spectral color component.

- Shininess

The shininess component.

### Ligh Source

The lamp is described with its position and color.

### The object

We have an object described somehow. This object is build (or made or grown etc.) by a given material.

### The calculation.

Based on the [learnopengl materials tutorial](https://learnopengl.com/Lighting/Materials), Setting material paragraph.

- First let's calculate the ambient component. 

```
ambientColorComponent = lightColor * material.ambient
```

And that's it. It's a 3D vector (`lightColor` \* `material.ambient`).

- The diffuse component.

It is a little bit complicated than the ambient, because we need more calculation. We need to know the position of the fragment (`vPosition`), the normal vector of the fragment in the given position (`vNormal`), the position of the light source (`lightPosition`), the light direction (`lightDirection`, we can calculate it from the light position and fragment position), and the color of the light source (`lightColor`).

```
norm = normalize(vNormal) // just to make sure, that it's normalized.
lightDirection = normalize(lightPosition - vPosition) // it's a direction vector, it has to be normalized. 
diffMult = max(dot(norm, lightDirection), 0.0) // it is the dot product of the lightDirection and norm vectors, but it can't be negative number.
diffuseColorComponent = lightColor * (diffMult * material.diffuse) // calculate the diffuse color component
```

It's also a 3D vector.

- The specular component

Calculating this is also fun. We need to know some further stuff for the calculation. We have to know the view position (`viewPosition`, the position of the camera) the view direction (`viewDirection`, we can calculate it from the lightPosition and the viewPosition), and the direction of the reflection (`reflectionDirection`, we can calculate it based on the light direction and the normal vector of the object)

```
viewDirection = normalize(viewPosition - lightPosition) // The light is coming from the lamp to the camera (eye)
reflectionDirection = reflect(-lightDirection, norm) // The calculation of the reflection direction
specMult = pow(max(dot(viewDirection, reflectionDirection), 0.0), material.shininess); // the dot product of the view direction and reflection direction, but it can't be negative, and it is powered to the shininess.
specularColorComponent = lightColor * (specMult * material.specular)
```

And a 3D vector again. The last step is putting the stuff together.

```
resultColor = ambientColorComponent + diffuseColorComponent + specularColorComponent
```

## New

This function returns a material. The inputs are the color component vectors (`ambient`, `diffuse`, `specular`) and the shininess.

## GetAmbient

Returns the ambient color component of the material.

## GetDiffuse

Returns the diffuse color component of the material.

## GetSpecular

Returns the specular color component of the material.

## GetShininess

Returns the shiniess of the material.
