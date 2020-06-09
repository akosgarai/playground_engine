# Light

This package aims to gather the lightsource related parameter into a package.

A basic light source is described with it's `position` vector, `ambient`, `diffuse`, `specular` color vector.

## Extend the material math.

For using the light source parameters, we have to extend the math mentioned in the material descriptions.

```
ambientColorComponent = light.ambient * material.ambient
diffuseColorComponent = light.diffuse * (diffMult * material.diffuse)
specularColorComponent = light.specular * (specMult * material.specular)
```

The paragraphs below are about the [light casters](https://learnopengl.com/Lighting/Light-casters) tutorial.

## Light casters

A light source that casts light upon object is called `light caster`. We can define different types of light casters.

- **Directional light**

This light source is far (ideally infinite far), then we can model it with directional light. A good example for this kind of light caster is the `Sun`. It's not infinite far, but very far. We have to update our light structure, if we want to distinguish the directional light sources from the others. Because for directional lights the position is irrelevant. Instead of this, we have to know the **direction**.

- **Point light**

This light source is given with it's position, and the light is going to every direction. Like a lamp in the ceiling. This kind of light reduces it's intensity over the distance (attenuation). To be able to use the calculation with this kind of light sources, we have to know some other variables. the `constant`, the `linear` and the `quadratic` terms.

- **Spot light**

This light source is given with it's position, and the light is going to a specific direction (like a flashlight). For the calculation we have to know the `cutoff` angle, that describes the size of the spot.

---

## NewPointLight

NewPointLight returns a Light with point light settings. The vectorComponent [4]mgl32.Vec3 input has to contain the `position`, `ambient`, `diffuse`, `specular` component vectors in this order. The terms [3]float32 input has to contain the `constant`, `linear`, `quadratic` term components in this order.

## NewDirectionalLight

NewDirectionalLight returns a Light with directional light settings. The vectorComponent [4]mgl32.Vec3 input has to contain the `direction`, `ambient`, `diffuse`, `specular` components in this order.

## NewSpotLight

NewSpotLight returns a Light with spot light settings. The vectorComponent [5]mgl32.Vec3 input has to contain the `position`, `direction`, `ambient`, `diffuse`, `specular` components in this order. The terms[5]float32 input has to contain the `constant`, `linear`, `quadratic` terms, `cutoff` and the `outerCutoff` components in this order.
