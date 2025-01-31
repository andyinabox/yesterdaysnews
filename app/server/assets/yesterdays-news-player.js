// import { ClipLoader } from '/assets/clip-loader.js'
// import { CaptionLoader } from '/assets/caption-loader.js'
// import { CaptionFlow } from '/assets/caption-flow.js'

import { CaptionRenderer } from '/assets/caption-renderer.js'
import { VideoRenderer } from '/assets/video-renderer.js'
import { StaticRenderer } from '/assets/static-renderer.js'

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
    const { x, y, width, height } = this.#getVideoBoundingBox()
    this.staticRenderer.draw(ctx, width, height)
    this.videoRenderer.draw(ctx, x, y, width, height)
    this.captionRenderer.draw(ctx, x, y, width, height)
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
    console.log('clips-resource-url', this.getAttribute('clips-resource-url'))
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

  // TODO: I think it should be possible to use transforms better
  // so that this can be eliminated (x, y will just be 0, 0) and
  // the transforms will ensure everything is positioned correctly
  #getVideoBoundingBox() {
    const { width, height } = this.#getVideoDimensions()

    let x, y

    if (this.rotated) {
      x = (this.canvas.height - width) / 2
      y = (this.canvas.width - height) / 2
    } else {
      x = (this.canvas.width - width) / 2
      y = (this.canvas.height - height) / 2
    }

    return { x, y, width, height }
  }
}

customElements.define('yesterdays-news-player', YesterdaysNewsPlayer)

// const ASPECT_RATIO = 720 / 1280
// // const STYLESHEET_URL = '/assets/yesterdays-news-player.css'

// export class YesterdaysNewsPlayer extends HTMLElement {
//   static observedAttributes = ['rotated', 'fullscreen', 'width', 'height']

//   #currentCaption = ''

//   constructor() {
//     super()
//     this._internals = this.attachInternals()
//     // initialize video clip loader
//     this.clipLoader = new ClipLoader(this.clipsResourceURL, this.initialClipURL)
//     // initialize caption loader
//     this.captionLoader = new CaptionLoader(
//       this.captionsResourceURL,
//       (caption) => {
//         this.#currentCaption = caption
//       }
//     )

//     this.captionFlow = new CaptionFlow()
//   }

//   connectedCallback() {
//     // start video player
//     this.clipLoader.video.play()
//     // connect to caption loader
//     this.captionLoader.connect()

//     // set up shadow dom
//     this.shadow = this.attachShadow({ mode: 'open' })
//     this.loadStyles()

//     // set styles
//     const sheet = new CSSStyleSheet()
//     sheet.replaceSync(`@import url("${STYLESHEET_URL}");`)
//     this.shadow.adoptedStyleSheets = [sheet]

//     // setup canvas
//     this.canvas = document.createElement('canvas')
//     this.ctx = this.canvas.getContext('2d')
//     this.setCanvasSize()
//     this.shadow.appendChild(this.canvas)

//     // setup nav
//     this.nav = document.createElement('nav')
//     this.aboutBtn = document.createElement('a')
//     this.aboutBtn.setAttribute('target', '_blank')
//     this.aboutBtn.href = '/about'
//     this.aboutBtn.innerHTML = 'about'
//     this.fullscreenBtn = document.createElement('button')
//     this.fullscreenBtn.addEventListener('click', () => {
//       this.requestFullscreen()
//     })
//     this.fullscreenBtn.innerHTML = 'fullscreen'
//     this.nav.appendChild(this.aboutBtn)
//     this.nav.appendChild(this.fullscreenBtn)
//     this.shadow.append(this.nav)
//     // add window resize listener
//     // TODO: add throttling
//     window.addEventListener('resize', this.setCanvasSize.bind(this))

//     // start canvas animation
//     window.requestAnimationFrame(this.drawCanvas.bind(this))

//     document.addEventListener(
//       'fullscreenchange',
//       this.handleFullscreenChange.bind(this)
//     )

//     // handle orientation change
//     screen.orientation.addEventListener(
//       'change',
//       this.checkOrientation.bind(this)
//     )
//     this.checkOrientation()
//   }

//   disconnectedCallback() {}

//   attributeChangedCallback(name, oldValue, newValue) {
//     switch (name) {
//       case 'rotated':
//         break
//       case 'fullscreen':
//         break
//       // CAUTION: intentional switch fallthrough
//       case 'width':
//       case 'height':
//         // handle change to width or height
//         break
//     }
//   }

//   loadStyles() {
//     const link = document.createElement('link')
//     link.setAttribute('href', STYLESHEET_URL)
//     link.setAttribute('rel', 'stylesheet')
//     link.setAttribute('type', 'text/css')
//     this.shadow.appendChild(link)
//   }

//   drawCanvas() {
//     this.clearCanvas()
//     this.setRotation()
//     const video = this.clipLoader.video
//     if (video.paused || video.ended) {
//       this.drawStatic()
//     } else {
//       this.drawVideo()
//     }
//     this.drawCaption()
//     window.requestAnimationFrame(this.drawCanvas.bind(this))
//   }

//   clearCanvas() {
//     this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height)
//   }

//   setRotation() {
//     this.ctx.resetTransform()
//     if (this.rotated) {
//       this.ctx.translate(this.canvas.width, 0)
//       this.ctx.rotate((90 * Math.PI) / 180)
//     }
//   }

//   drawCaption() {
//     // get video bounding box
//     const { x, y, width, height } = this.getVideoBoundingBox()

//     // set canvas based on video width
//     const fontSize = 0.04 * width
//     const fontFamily = 'monospace'
//     const verticalPadding = 0.1 * fontSize
//     const horizontalPadding = 0.2 * fontSize
//     const maxCaptionWidth = width - 4 * horizontalPadding

//     const lines = this.captionFlow.get(
//       this.ctx,
//       this.rotated,
//       fontSize,
//       fontFamily,
//       maxCaptionWidth,
//       this.#currentCaption
//     )

//     // go through the lines of text and draw
//     let textWidth, textHeight, textX, textY, metrics
//     for (let i = 0; i < lines.length; i++) {
//       // update metrics
//       metrics = this.ctx.measureText(lines[i])

//       // calculate bounding box
//       textWidth = metrics.width
//       textHeight =
//         metrics.fontBoundingBoxAscent + metrics.fontBoundingBoxDescent

//       if (this.rotated) {
//         textX = this.canvas.height / 2 - metrics.width / 2
//       } else {
//         textX = this.canvas.width / 2 - metrics.width / 2
//       }
//       textY = y + height - fontSize - (textHeight + 2 * verticalPadding) * i

//       // draw black box
//       this.ctx.fillStyle = 'rgb(0 0 0)'
//       this.ctx.fillRect(
//         textX - horizontalPadding,
//         textY - verticalPadding - fontSize,
//         textWidth + 2 * horizontalPadding,
//         textHeight + 2 * verticalPadding
//       )

//       // draw text
//       this.ctx.fillStyle = 'rgb(255 255 255)'
//       this.ctx.fillText(lines[i], textX, textY)
//     }
//   }

//   drawVideo() {
//     const { x, y, width, height } = this.getVideoBoundingBox()
//     this.ctx.drawImage(this.clipLoader.video, x, y, width, height)
//   }

//   drawStatic() {
//     const { width, height } = this.getVideoDimensions()

//     let imageData
//     if (this.rotated) {
//       imageData = this.ctx.createImageData(height, width)
//     } else {
//       imageData = this.ctx.createImageData(width, height)
//     }
//     const data = imageData.data

//     for (let i = 0; i < data.length; i += 4) {
//       const value = Math.random() * 64
//       data[i] = value // red
//       data[i + 1] = value // green
//       data[i + 2] = value // blue
//       data[i + 3] = 255 // alpha
//     }

//     let x, y
//     if (this.rotated) {
//       x = this.canvas.width / 2 - height / 2
//       y = this.canvas.height / 2 - width / 2
//     } else {
//       x = this.canvas.width / 2 - width / 2
//       y = this.canvas.height / 2 - height / 2
//     }

//     this.ctx.putImageData(imageData, x, y)
//   }

//   getVideoDimensions() {
//     let width, height
//     // if (this.rotated) {
//     //   width = this.canvas.height
//     // } else {
//     //   width = this.canvas.width
//     // }

//     if (this.rotated) {
//       if (this.canvas.height * ASPECT_RATIO > this.canvas.width) {
//         height = this.canvas.width
//         width = height / ASPECT_RATIO
//       } else {
//         width = this.canvas.height
//         height = width * ASPECT_RATIO
//       }
//     } else {
//       if (this.canvas.width * ASPECT_RATIO > this.canvas.height) {
//         height = this.canvas.height
//         width = height / ASPECT_RATIO
//       } else {
//         width = this.canvas.width
//         height = width * ASPECT_RATIO
//       }
//     }

//     // const height = width * ASPECT_RATIO

//     return { width, height }
//   }

//   getVideoBoundingBox() {
//     const { width, height } = this.getVideoDimensions()

//     let x, y

//     if (this.rotated) {
//       x = (this.canvas.height - width) / 2
//       y = (this.canvas.width - height) / 2
//     } else {
//       x = (this.canvas.width - width) / 2
//       y = (this.canvas.height - height) / 2
//     }

//     return { x, y, width, height }
//   }

//   handleFullscreenChange() {
//     if (document.fullscreenElement === this) {
//       this.fullscreen = true
//     } else {
//       this.fullscreen = false
//     }

//     this.checkOrientation()
//   }

//   checkOrientation() {
//     if (!this.fullscreen) {
//       this.rotated = false
//       return
//     }

//     if (screen.orientation.type.includes('portrait')) {
//       this.rotated = true
//     } else {
//       this.rotated = false
//     }
//   }

//   setCanvasSize() {
//     const { width, height } = this.getBoundingClientRect()
//     this.canvas.setAttribute('width', width + 'px')
//     this.canvas.setAttribute('height', height + 'px')
//   }

//   disconnectedCallback() {
//     this.captionLoader.disconnect()
//   }

//   get width() {
//     return this.getAttribute('width')
//   }

//   get height() {
//     return this.getAttribute('height')
//   }

//   get rotated() {
//     return this.getAttribute('rotated')
//   }

//   get fullscreen() {
//     return this.getAttribute('fullscreen')
//   }

//   get clipsResourceURL() {
//     return this.getAttribute('clips-resource-url')
//   }

//   get captionsResourceURL() {
//     return this.getAttribute('captions-resource-url')
//   }

//   get initialClipURL() {
//     return this.getAttribute('initial-clip-url')
//   }

//   static register() {
//     customElements.define('yesterdays-news-player', YesterdaysNewsPlayer)
//   }
// }
