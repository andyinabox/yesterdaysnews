import { ClipLoader } from '/assets/clip-loader.js'
import { CaptionLoader } from '/assets/caption-loader.js'

const ASPECT_RATIO = 720 / 1280

export class YesterdaysNewsPlayer extends HTMLElement {
  static css = `
    canvas {
      width: 100%;
      height: 100%;
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
    // this.shadow = this.attachShadow({ mode: 'open' })

    // set styles
    const sheet = new CSSStyleSheet()
    sheet.replaceSync(YesterdaysNewsPlayer.css)
    this.adoptedStyleSheets = [sheet]

    this.canvas = document.createElement('canvas')
    const { width, height } = this.getBoundingClientRect()
    console.log('width, height', width, height)
    this.canvas.setAttribute('width', width + 'px')
    this.canvas.setAttribute('height', height + 'px')
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

    this.appendChild(this.canvas)

    window.requestAnimationFrame(this.drawCanvas.bind(this))
  }

  drawCanvas() {
    this.clearCanvas()
    this.drawVideo()
    this.drawCaption()
    window.requestAnimationFrame(this.drawCanvas.bind(this))
  }

  clearCanvas() {
    this.ctx.fillStyle = 'rgb(0 0 0)'
    this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height)
    // this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height)
  }

  drawCaption() {
    const fontSize = 0.035 * this.canvas.width
    const padding = 0.5 * fontSize
    const [videoX, videoY, videoWidth, videoHeight] =
      this.calcVideoDimensionsAndLocation()

    this.ctx.font = `${fontSize}px monospace`
    const metrics = this.ctx.measureText(this.#currentCaption)

    const x = this.canvas.width / 2 - metrics.width / 2
    const y = videoY + videoHeight - fontSize

    this.ctx.fillStyle = 'rgb(0 0 0)'
    this.ctx.fillRect(x, y - fontSize, metrics.width, fontSize)

    this.ctx.fillStyle = 'rgb(255 255 255)'
    this.ctx.fillText(this.#currentCaption, x, y)
  }

  drawVideo() {
    const [x, y, width, height] = this.calcVideoDimensionsAndLocation()

    this.ctx.drawImage(this.video, x, y, width, height)
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
