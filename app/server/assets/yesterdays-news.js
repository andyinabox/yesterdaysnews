export class YesterdaysNews extends HTMLElement {
  static stylesheetURL = '/assets/yesterdays-news.css'

  constructor() {
    super()
    this._internals = this.attachInternals()
  }

  connectedCallback() {
    this.shadow = this.attachShadow({ mode: 'open' })

    // load stylesheet
    const link = document.createElement('link')
    link.setAttribute('href', YesterdaysNews.stylesheetURL)
    link.setAttribute('rel', 'stylesheet')
    link.setAttribute('type', 'text/css')
    this.shadow.appendChild(link)

    // setup player
    this.player = document.createElement('yesterdays-news-player')
    this.player.setAttribute('clips-resource-url', this.clipsResourceURL)
    this.player.setAttribute('captions-resource-url', this.captionsResourceURL)
    this.player.setAttribute('initial-clip-url', this.initialClipURL)
    this.player.setup()

    // setup nav
    this.nav = docyment.createElement('yesterdays-news-nav')
    this.nav.setup()

    // add to shadow dom
    this.shadow.appendChild(this.player)
    this.shadow.appendChild(this.nav)

    // set event listeners
    // todo: add debounce for resize
    window.addEventListener('resize', this.handleWindowResize.bind(this))
    document.addEventListener(
      'fullscreenchange',
      this.handleFullscreenChange.bind(this)
    )
    screen.orientation.addEventListener(
      'change',
      this.handleOrientationChange.bind(this)
    )

    // do initial orientation check
    this.handleOrientationChange()
  }

  disconnectedCallback() {
    window.removeEventListener('resize', this.handleWindowResize.bind(this))
    document.removeEventListener(
      'fullscreenchange',
      this.handleFullscreenChange.bind(this)
    )
    screen.orientation.removeEventListener(
      'change',
      this.handleOrientationChange.bind(this)
    )
  }

  handleWindowResize() {
    const { width, height } = this.getBoundingClientRect()
    this.canvas.setAttribute('width', width)
    this.canvas.setAttribute('height', height)
  }

  handleOrientationChange() {
    if (!this.fullscreen) {
      this.rotated = false
      return
    }

    if (screen.orientation.type.includes('portrait')) {
      this.rotated = true
    } else {
      this.rotated = false
    }
  }

  handleFullscreenChange() {
    if (document.fullscreenElement === this) {
      this.fullscreen = true
    } else {
      this.fullscreen = false
    }

    this.handleOrientationChange()
  }

  get rotated() {
    return this._internals.states.has('rotated')
  }

  set rotated(value) {
    if (value) {
      // Existence of identifier corresponds to "true"
      this._internals.states.add('rotated')
    } else {
      // Absence of identifier corresponds to "false"
      this._internals.states.delete('rotated')
    }

    this.player.setAttribute('rotated', value)
    this.nav.setAttribute('rotated', value)
  }

  get fullscreen() {
    return this._internals.states.has('fullscreen')
  }

  set fullscreen(value) {
    if (value) {
      // Existence of identifier corresponds to "true"
      this._internals.states.add('fullscreen')
    } else {
      // Absence of identifier corresponds to "false"
      this._internals.states.delete('fullscreen')
    }

    this.player.setAttribute('fullscreen', value)
    this.nav.setAttribute('fullscreen', value)
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
    customElements.define('yesterdays-news', YesterdaysNews)
  }
}
