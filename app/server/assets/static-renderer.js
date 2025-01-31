export class StaticRenderer {
  #rotated = false
  #fullscreen = false

  constructor() {}
  setup() {}
  update(rotated = false, fullscreen = false) {
    this.#rotated = rotated
    this.#fullscreen = fullscreen
  }

  draw(ctx, width, height) {
    let imageData
    if (this.#rotated) {
      imageData = ctx.createImageData(height, width)
    } else {
      imageData = ctx.createImageData(width, height)
    }
    const data = imageData.data

    for (let i = 0; i < data.length; i += 4) {
      const value = Math.random() * 64
      data[i] = value // red
      data[i + 1] = value // green
      data[i + 2] = value // blue
      data[i + 3] = 255 // alpha
    }

    let x, y
    if (this.#rotated) {
      x = ctx.canvas.width / 2 - height / 2
      y = ctx.canvas.height / 2 - width / 2
    } else {
      x = ctx.canvas.width / 2 - width / 2
      y = ctx.canvas.height / 2 - height / 2
    }

    // note that because we are using `putImageData`, transforms will not be applied
    ctx.putImageData(imageData, x, y)
  }
}
