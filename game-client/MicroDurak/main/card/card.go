components {
  id: "card"
  component: "/main/card/card.script"
}
embedded_components {
  id: "sprite"
  type: "sprite"
  data: "default_animation: \"CardBack\"\n"
  "material: \"/builtins/materials/sprite.material\"\n"
  "textures {\n"
  "  sampler: \"texture_sampler\"\n"
  "  texture: \"/assets/images/cards.atlas\"\n"
  "}\n"
  ""
  scale {
    x: 0.1
    y: 0.1
  }
}
