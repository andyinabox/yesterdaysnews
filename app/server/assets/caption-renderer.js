import { SingleObjectCache } from '/assets/single-object-cache.js'

export class CaptionRenderer {
  #resourceURL
  #eventSource

  #caption
  #captionLinesCache

  constructor(resourceURL) {
    this.#resourceURL = resourceURL
    this.#captionLinesCache = new SingleObjectCache()
  }

  setup() {
    this.#eventSource = new EventSource(this.#resourceURL)
    this.#eventSource.addEventListener('message', this.handleMessage.bind(this))
  }

  update() {}

  draw(ctx, x, y, width, height) {}

  destroy() {
    this.#eventSource.removeEventListener(
      'message',
      this.handleMessage.bind(this)
    )
  }

  #getCaptionLines(ctx, cachedAttributes = {}) {
    // if (this.#captionLinesCache.has(cachedAttributes))
  }

  handleMessage(event) {
    this.caption = event.data
  }

  set caption(str) {
    this.#caption = str
  }

  get caption() {
    return this.#caption
  }
}
