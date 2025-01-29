export class VideoPlayer extends HTMLElement {
  clips = []
  preloaded = []

  constructor() {
    super()
    this.loadClips()
  }

  connectedCallback() {
    // create video element
    const video = document.createElement('video')
    video.setAttribute('muted', true)
    video.setAttribute('autoplay', true)
    video.setAttribute('tabindex', '-1')
    video.addEventListener('ended', this.onVideoEnded.bind(this))
    video.addEventListener('error', this.onVideoError.bind(this))

    // create source element
    const videoSource = document.createElement('source')
    videoSource.setAttribute('type', 'video/webm')
    videoSource.setAttribute('src', this.initialClipUrl)

    video.appendChild(videoSource)

    this.video = video
    this.videoSource = videoSource
    this.appendChild(video)
  }

  async preloadNextClip() {
    // get the next url
    const url = this.clips.pop()
    // console.log('start preloading clip ' + url)
    // fetch the video
    try {
      const resp = await fetch(url)
      // get the video data as array buffer
      const data = await resp.arrayBuffer()
      // add to array of preloaded videos
      this.preloaded.push({
        url,
        objectURL: URL.createObjectURL(
          new Blob([data], { type: 'video/webm' })
        ),
      })
      // console.log('done preloading clip ' + url)
    } catch (err) {
      console.error('error preloading next clip: ', err)
      this.preloadNextClip()
    }
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
      // console.log('getting preloaded ObjectURL')
      const nextData = this.preloaded.pop()
      // console.log('returning ObjectURL for video ' + nextData.url)
      next = nextData.objectURL
    } else {
      // console.log('get next clip URL')
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

  get resourceUrl() {
    return this.getAttribute('resource-url')
  }

  get initialClipUrl() {
    return this.getAttribute('initial-clip-url')
  }

  static register() {
    customElements.define('video-player', VideoPlayer)
  }
}
