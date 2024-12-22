class VideoPlayer extends HTMLElement {
  constructor() {
    super()

    // create video element
    this.video = document.createElement('video')
    this.video.setAttribute('muted', true)
    this.video.addEventListener('ended', this.onVideoEnded.bind(this))
    this.video.addEventListener('error', this.onVideoError.bind(this))

    // create source element
    this.videoSource = document.createElement('source')
    this.videoSource.setAttribute('type', 'video/webm')

    // add source to video
    this.video.appendChild(this.videoSource)
  }
  connectedCallback() {
    this.resourceUrl = this.getAttribute('resource-url')

    // load the list of clips
    this.clipsLoading = this.loadClips()

    // append the main video player
    this.appendChild(this.video)

    // load and play the initial clip
    this.loadNewVideo(this.getAttribute('initial-clip'))
  }

  async loadClips() {
    try {
      const resp = await fetch(this.resourceUrl)

      if (!resp.ok) {
        throw new Error(`Response status: ${resp.status}`)
      }

      const data = await resp.json()

      this.clips = data.clips
    } catch (err) {
      console.error(err)
    }
  }

  onVideoError(err) {
    console.error(err)
    this.loadNewVideo(this.getNextClip())
  }

  onVideoEnded() {
    this.loadNewVideo(this.getNextClip())
  }

  getNextClip() {
    if (this.clips.length < 10) {
      this.loadClips()
    }

    const next = this.clips.pop()

    // preload next video somehow?

    return next
  }

  loadNewVideo(url) {
    try {
      this.video.pause()

      this.videoSource.setAttribute('src', url)

      this.video.load()
      this.video.play()
    } catch (err) {
      console.log(`error loading video ${url}`, err)
    }
  }

  disconnectedCallback() {
    this.video.removeEventListener('ended', this.onVideoEnded.bind(this))
    this.video.removeEventListener('error', this.onVideoError.bind(this))
  }
}
customElements.define('video-player', VideoPlayer)
