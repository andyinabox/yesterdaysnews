class VideoPlayer extends HTMLElement {
  constructor() {
    super()

    // create video element
    this.video = document.createElement('video')
    this.video.setAttribute('muted', true)
    this.video.addEventListener('ended', this.onVideoEnded.bind(this))
    this.video.addEventListener('error', this.onVideoEnded.bind(this))

    // create source element
    this.videoSource = document.createElement('source')
    this.videoSource.setAttribute('type', 'video/webm')

    // add source to video
    this.video.appendChild(this.videoSource)
  }
  connectedCallback() {
    this.resourceUrl = this.getAttribute('resource-url')
    this.clipsLoading = this.loadClips()

    const initialClip = this.getAttribute('initial-clip')

    this.appendChild(this.video)

    this.loadNewVideo(initialClip)
  }

  async loadClips() {
    try {
      const resp = await fetch(this.resourceUrl)
      if (!resp.ok) {
        throw new Error(`Response status: ${resp.status}`)
      }

      const data = await resp.json()

      this.clips = data.clips
      console.log('clips loaded', this.clips)
    } catch (err) {
      console.error(err)
    }
  }

  onVideoEnded() {
    this.loadNewVideo(this.getNextClip())
  }

  getNextClip() {
    return this.clips.pop()
  }

  loadNewVideo(url) {
    console.log('loadNewVideo', url)
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
    this.video.removeEventListener('error', this.onVideoEnded.bind(this))
  }
}
customElements.define('video-player', VideoPlayer)
