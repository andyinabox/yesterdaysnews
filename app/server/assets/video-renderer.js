export class VideoRenderer {
  #resourceURL
  #initialURL

  constructor(resourceURL, initialURL) {
    this.#resourceURL = resourceURL
    this.#initialURL = initialURL
  }

  setup() {}

  update() {}

  draw(ctx, x, y, width, height) {}
}
