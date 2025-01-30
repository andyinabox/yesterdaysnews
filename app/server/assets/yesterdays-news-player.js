import { ClipLoader } from '/assets/clip-loader.js'
import { CaptionLoader } from '/assets/caption-loader.js'

const ASPECT_RATIO = 720 / 1280
const STYLESHEET_URL = '/assets/yesterdays-news-player.css'

export class YesterdaysNewsPlayer extends HTMLElement {
  #currentCaption = ''
  #rotated = false

  constructor() {
    super()
    // initialize video clip loader
    this.clipLoader = new ClipLoader(this.clipsResourceURL, this.initialClipURL)
    // initialize caption loader
    this.captionLoader = new CaptionLoader(
      this.captionsResourceURL,
      (caption) => {
        this.#currentCaption = caption
      }
    )
  }

  connectedCallback() {
    // start video player
    this.clipLoader.video.play()
    // connect to caption loader
    this.captionLoader.connect()

    // set up shadow dom
    this.shadow = this.attachShadow({ mode: 'open' })
    this.loadStyles()

    // set styles
    const sheet = new CSSStyleSheet()
    sheet.replaceSync(`@import url("${STYLESHEET_URL}");`)
    this.shadow.adoptedStyleSheets = [sheet]

    // setup canvas
    this.canvas = document.createElement('canvas')
    this.ctx = this.canvas.getContext('2d')
    this.setCanvasSize()
    this.shadow.appendChild(this.canvas)

    // add window resize listener
    // TODO: add throttling
    window.addEventListener('resize', this.setCanvasSize.bind(this))

    // start canvas animation
    window.requestAnimationFrame(this.drawCanvas.bind(this))
  }

  loadStyles() {
    const link = document.createElement('link')
    link.setAttribute('href', STYLESHEET_URL)
    link.setAttribute('rel', 'stylesheet')
    link.setAttribute('type', 'text/css')
    this.shadow.appendChild(link)
  }

  drawCanvas() {
    this.clearCanvas()
    this.setRotation()
    const video = this.clipLoader.video
    if (video.paused || video.ended) {
      this.drawStatic()
    } else {
      this.drawVideo()
    }
    this.drawCaption()
    window.requestAnimationFrame(this.drawCanvas.bind(this))
  }

  clearCanvas() {
    this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height)
  }

  setRotation() {
    this.ctx.resetTransform()
    if (this.#rotated) {
      this.ctx.translate(this.canvas.width, 0)
      this.ctx.rotate((90 * Math.PI) / 180)
    }
  }

  drawCaption() {
    // get video bounding box
    const { x, y, width, height } = this.getVideoBoundingBox()

    // set canvas based on video width
    const fontSize = 0.04 * width
    const verticalPadding = 0.1 * fontSize
    const horizontalPadding = 0.2 * fontSize

    // get current caption
    const caption = this.#currentCaption

    // set font size and then use that to get text metricss
    this.ctx.font = `${fontSize}px monospace`
    let metrics = this.ctx.measureText(caption)

    const lines = [caption]

    // break into lines if necessart
    if (metrics.width > width - 4 * horizontalPadding) {
      const tokens = caption.split(' ')
      const divider = Math.ceil(tokens.length / 2)
      lines[0] = tokens.slice(0, divider).join(' ')
      lines[1] = tokens.slice(divider).join(' ')
    }

    // go through the lines of text and draw
    let textWidth, textHeight, textX, textY
    for (let i = 0; i < lines.length; i++) {
      // update metrics
      metrics = this.ctx.measureText(lines[i])

      // calculate bounding box
      textWidth = metrics.width
      textHeight =
        metrics.fontBoundingBoxAscent + metrics.fontBoundingBoxDescent

      if (this.#rotated) {
        textX = this.canvas.height / 2 - metrics.width / 2
      } else {
        textX = this.canvas.width / 2 - metrics.width / 2
      }
      textY = y + height - fontSize - (textHeight + 2 * verticalPadding) * i

      // draw black box
      this.ctx.fillStyle = 'rgb(0 0 0)'
      this.ctx.fillRect(
        textX - horizontalPadding,
        textY - verticalPadding - fontSize,
        textWidth + 2 * horizontalPadding,
        textHeight + 2 * verticalPadding
      )

      // draw text
      this.ctx.fillStyle = 'rgb(255 255 255)'
      this.ctx.fillText(lines[i], textX, textY)
    }
  }

  drawVideo() {
    const { x, y, width, height } = this.getVideoBoundingBox()
    this.ctx.drawImage(this.clipLoader.video, x, y, width, height)
  }

  drawStatic() {
    const { width, height } = this.getVideoDimensions()

    let imageData
    if (this.#rotated) {
      imageData = this.ctx.createImageData(height, width)
    } else {
      imageData = this.ctx.createImageData(width, height)
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
      x = this.canvas.width / 2 - height / 2
      y = this.canvas.height / 2 - width / 2
    } else {
      x = this.canvas.width / 2 - width / 2
      y = this.canvas.height / 2 - height / 2
    }

    this.ctx.putImageData(imageData, x, y)
  }

  getVideoDimensions() {
    let width, height
    // if (this.#rotated) {
    //   width = this.canvas.height
    // } else {
    //   width = this.canvas.width
    // }

    if (this.#rotated) {
      if (this.canvas.height * ASPECT_RATIO > this.canvas.width) {
        height = this.canvas.width
        width = height / ASPECT_RATIO
      } else {
        width = this.canvas.height
        height = width * ASPECT_RATIO
      }
    } else {
      if (this.canvas.width * ASPECT_RATIO > this.canvas.height) {
        height = this.canvas.height
        width = height / ASPECT_RATIO
      } else {
        width = this.canvas.width
        height = width * ASPECT_RATIO
      }
    }

    // const height = width * ASPECT_RATIO

    return { width, height }
  }

  getVideoBoundingBox() {
    const { width, height } = this.getVideoDimensions()

    let x, y

    if (this.#rotated) {
      x = (this.canvas.height - width) / 2
      y = (this.canvas.width - height) / 2
    } else {
      x = (this.canvas.width - width) / 2
      y = (this.canvas.height - height) / 2
    }

    return { x, y, width, height }
  }

  setCanvasSize() {
    const { width, height } = this.getBoundingClientRect()
    this.canvas.setAttribute('width', width + 'px')
    this.canvas.setAttribute('height', height + 'px')
  }

  disconnectedCallback() {
    this.captionLoader.disconnect()
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

  static register() {
    customElements.define('yesterdays-news-player', YesterdaysNewsPlayer)
  }
}
