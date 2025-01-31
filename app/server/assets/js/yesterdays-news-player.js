import { CaptionRenderer } from './caption-renderer.js'
import { VideoRenderer } from './video-renderer.js'
import { StaticRenderer } from './static-renderer.js'

export class YesterdaysNewsPlayer extends HTMLElement {
  static observedAttributes = ['rotated', 'fullscreen', 'width', 'height']

  constructor() {
    super()
  }

  connectedCallback() {
    // initialize renderers
    this.captionRenderer = new CaptionRenderer(this.captionsResourceURL)
    this.videoRenderer = new VideoRenderer(
      this.clipsResourceURL,
      this.initialClipURL
    )
    this.staticRenderer = new StaticRenderer()

    this.canvas = document.createElement('canvas')
    this.ctx = this.canvas.getContext('2d')
    this.appendChild(this.canvas)

    this.setup()
  }

  disconnectedCallback() {}

  play() {
    this.videoRenderer.play()
  }

  // start playback
  setup() {
    this.captionRenderer.setup()
    this.videoRenderer.setup()
    this.staticRenderer.setup()
    const tick = () => {
      this.draw()
      window.requestAnimationFrame(tick)
    }
    window.requestAnimationFrame(tick)
  }

  update() {
    this.canvas.setAttribute('width', this.width + 'px')
    this.canvas.setAttribute('height', this.height + 'px')

    const rotated = this.rotated
    const fullscreen = this.fullscreen
    this.captionRenderer.update(rotated, fullscreen)
    this.videoRenderer.update(rotated, fullscreen)
    this.staticRenderer.update(rotated, fullscreen)
  }

  draw() {
    const ctx = this.canvas.getContext('2d')
    const { width, height } = this.#getVideoDimensions()

    // clear the frame
    ctx.clearRect(0, 0, this.canvas.width, this.canvas.height)

    // set transform
    ctx.resetTransform()
    this.setTranslation(ctx, width, height)
    this.setRotation(ctx)

    // draw parts
    // this.staticRenderer.draw(ctx, width, height)
    this.videoRenderer.draw(ctx, 0, 0, width, height)
    this.captionRenderer.draw(ctx, 0, 0, width, height)
  }

  setTranslation(ctx, width, height) {
    let x, y

    if (this.rotated) {
      x = ctx.canvas.width / 2 + height / 2
      y = ctx.canvas.height / 2 - width / 2
    } else {
      x = (ctx.canvas.width - width) / 2
      y = (ctx.canvas.height - height) / 2
    }

    ctx.translate(x, y)
  }

  setRotation(ctx) {
    if (this.rotated) {
      ctx.rotate((90 * Math.PI) / 180)
    }
  }

  attributeChangedCallback(name, oldValue, newValue) {
    switch (name) {
      case 'rotated':
        this.update()
        break
      case 'fullscreen':
        this.update()
        break
      // CAUTION: intentional switch fallthrough
      case 'width':
      case 'height':
        // handle change to width or height
        this.update()
        break
    }
  }

  get width() {
    return this.getAttribute('width')
  }

  get height() {
    return this.getAttribute('height')
  }

  get rotated() {
    return (
      this.hasAttribute('rotated') && this.getAttribute('rotated') !== 'false'
    )
  }

  get fullscreen() {
    return (
      this.hasAttribute('rotated') && this.getAttribute('rotated') !== 'false'
    )
  }

  get clipsResourceURL() {
    return this.getAttribute('clips-resource-url')
  }

  get captionsResourceURL() {
    return this.getAttribute('captions-resource-url')
  }

  get initialClipURL() {
    return this.getAttribute('initial-clip-url')
  }

  get videoAspectRatio() {
    return parseFloat(this.getAttribute('video-aspect-ratio'))
  }

  #getVideoDimensions() {
    const aspectRatio = this.videoAspectRatio
    let width, height

    if (this.rotated) {
      if (this.canvas.height * aspectRatio > this.canvas.width) {
        height = this.canvas.width
        width = height / aspectRatio
      } else {
        width = this.canvas.height
        height = width * aspectRatio
      }
    } else {
      if (this.canvas.width * aspectRatio > this.canvas.height) {
        height = this.canvas.height
        width = height / aspectRatio
      } else {
        width = this.canvas.width
        height = width * aspectRatio
      }
    }

    return { width, height }
  }
}

customElements.define('yesterdays-news-player', YesterdaysNewsPlayer)
