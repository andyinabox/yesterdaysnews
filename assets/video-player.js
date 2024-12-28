class VideoPlayer extends HTMLElement {
  clips = []
  preloaded = []

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

  async preloadNextClip() {
    // get the next url
    const url = this.clips.pop()
    // fetch the video
    const resp = await fetch(url)
    // get the video data as array buffer
    const data = await resp.arrayBuffer()
    // create an ObjectURL from teh array buffer
    const objectURL = URL.createObjectURL(data)
    // add to array of preloaded videos
    this.preloaded.push(objectURL)
  }

  async loadClips() {
    try {
      const resp = await fetch(this.resourceUrl)

      if (!resp.ok) {
        throw new Error(`Response status: ${resp.status}`)
      }

      const data = await resp.json()

      this.clips = data.clips
      this.preloadNextClip()
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

    let next
    if (this.preloaded.length) {
      console.log('getting preloaded ObjectURL')
      next = this.preloaded.pop()
    } else {
      console.log('get next clip URL')
      next = this.clips.pop()
    }

    this.preloadNextClip()

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
