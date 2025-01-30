import { ClipLoader } from '/assets/clip-loader.js'
import { CaptionLoader } from '/assets/caption-loader.js'

const ASPECT_RATIO = 720 / 1280

export class YesterdaysNewsPlayer extends HTMLElement {
  static css = `
    :host {
      display: flex;
      justify-content: center;
      align-items: center;
      background-color: black;
    }
  `

  #currentCaption = ''

  constructor() {
    super()
    this.clipLoader = new ClipLoader(this.clipsResourceURL)
    this.captionLoader = new CaptionLoader(
      this.captionsResourceURL,
      (caption) => {
        this.#currentCaption = caption
      }
    )
  }

  connectedCallback() {
    // connect to caption loader
    this.captionLoader.connect()
    this.shadow = this.attachShadow({ mode: 'open' })

    // set styles
    const sheet = new CSSStyleSheet()
    sheet.replaceSync(YesterdaysNewsPlayer.css)
    this.shadow.adoptedStyleSheets = [sheet]

    this.canvas = document.createElement('canvas')
    this.onResize()
    this.ctx = this.canvas.getContext('2d')

    this.video = document.createElement('video')
    this.video.setAttribute('muted', true)
    this.video.setAttribute('tabindex', '-1')
    this.video.addEventListener('ended', this.onVideoEnded.bind(this))
    this.video.addEventListener('error', this.onVideoError.bind(this))

    this.videoSource = document.createElement('source')
    this.videoSource.setAttribute('type', 'video/webm')
    this.videoSource.setAttribute('src', this.initialClipURL)

    this.video.appendChild(this.videoSource)

    this.video.play()

    this.shadow.appendChild(this.canvas)

    window.addEventListener('resize', this.onResize.bind(this))

    window.requestAnimationFrame(this.drawCanvas.bind(this))
  }

  drawCanvas() {
    this.clearCanvas()
    // this.setRotation()
    if (this.video.paused || this.video.ended) {
      this.drawStatic()
    } else {
      this.drawVideo()
    }
    this.drawCaption()
    window.requestAnimationFrame(this.drawCanvas.bind(this))
  }

  clearCanvas() {
    this.ctx.fillStyle = 'rgb(0 0 0)'
    this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height)
    // this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height)
  }

  setRotation() {
    this.ctx.resetTransform()
    if (true) {
      this.ctx.translate(this.canvas.width / 2, -this.canvas.height / 2)
      this.ctx.rotate((90 * Math.PI) / 180)
    }
  }

  drawCaption() {
    const fontSize = 0.035 * this.canvas.width
    const padding = 0.5 * fontSize
    const [videoX, videoY, videoWidth, videoHeight] =
      this.calcVideoDimensionsAndLocation()
    const caption = this.#currentCaption

    this.ctx.font = `${fontSize}px monospace`
    let metrics = this.ctx.measureText(caption)

    let lines = []
    if (metrics.width > videoWidth) {
      const tokens = caption.split(' ')
      const divider = Math.ceil(tokens.length / 2)
      lines[0] = tokens.slice(0, divider).join(' ')
      lines[1] = tokens.slice(divider).join(' ')
    } else {
      lines = [caption]
    }

    for (let i = 0; i < lines.length; i++) {
      metrics = this.ctx.measureText(lines[i])
      const width = metrics.width
      const height =
        metrics.fontBoundingBoxAscent + metrics.fontBoundingBoxDescent

      const x = this.canvas.width / 2 - metrics.width / 2
      const y = videoY + videoHeight - fontSize - height * i

      this.ctx.fillStyle = 'rgb(0 0 0)'
      this.ctx.fillRect(x, y - fontSize, width, height)

      this.ctx.fillStyle = 'rgb(255 255 255)'
      this.ctx.fillText(lines[i], x, y)
    }
  }

  drawVideo() {
    const [x, y, width, height] = this.calcVideoDimensionsAndLocation()

    this.ctx.drawImage(this.video, x, y, width, height)
  }

  drawStatic() {
    const [videoX, videoY, videoWidth, videoHeight] =
      this.calcVideoDimensionsAndLocation()
    const imageData = this.ctx.createImageData(videoWidth, videoHeight)
    const data = imageData.data

    for (let i = 0; i < data.length; i += 4) {
      const value = Math.random() * 64
      data[i] = value // red
      data[i + 1] = value // green
      data[i + 2] = value // blue
      data[i + 3] = 255 // alpha
    }

    this.ctx.putImageData(imageData, videoX, videoY)
  }

  calcVideoDimensions() {
    const width = this.canvas.width
    const height = this.canvas.width * ASPECT_RATIO
    return [width, height]
  }

  calcVideoDimensionsAndLocation() {
    const [width, height] = this.calcVideoDimensions()
    const x = 0
    const y = (this.canvas.height - height) / 2
    return [x, y, width, height]
  }

  onVideoEnded() {
    this.loadNewVideo()
  }

  onVideoError() {
    console.error('video error')
    this.loadNewVideo()
  }

  onResize() {
    const { width, height } = this.getBoundingClientRect()
    this.canvas.setAttribute('width', width + 'px')
    this.canvas.setAttribute('height', width * ASPECT_RATIO)
  }

  loadNewVideo() {
    try {
      this.video.pause()

      this.videoSource.setAttribute('src', this.clipLoader.next())

      this.video.load()
      this.video.play()
    } catch (err) {
      console.log(`error loading video ${url}`, err)
    }
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
