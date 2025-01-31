export class VideoRenderer {
  #resourceURL

  #clips = []
  #preloaded = []

  #video = null
  #videoSource = null

  #rotated = false
  #fullscreen = false

  constructor(resourceURL, initialClipURL) {
    this.#resourceURL = resourceURL

    this.#fetchClipURLs()

    // set up video element for loading videos
    this.#video = document.createElement('video')
    this.#video.setAttribute('muted', true) // this is required to allow autoplay
    this.#video.addEventListener('ended', this.handleVideoEnded.bind(this))
    this.#video.addEventListener('error', this.handleVideoError.bind(this))

    // set up video source
    this.#videoSource = document.createElement('source')
    this.#videoSource.setAttribute('type', 'video/webm')
    this.#videoSource.setAttribute('src', initialClipURL)

    // add source to video
    this.#video.appendChild(this.#videoSource)
  }

  setup() {
    this.#video.play()
  }

  update(rotated = false, fullscreen = false) {
    this.#rotated = rotated
    this.#fullscreen = fullscreen
  }

  draw(ctx, x, y, width, height) {
    ctx.drawImage(this.#video, x, y, width, height)
  }

  destroy() {
    this.#video.pause()
    this.#video.removeEventListener('ended', this.handleVideoEnded.bind(this))
    this.#video.removeEventListener('error', this.handleVideoError.bind(this))
  }

  handleVideoEnded() {
    this.#loadNewVideo()
  }

  handleVideoError(err) {
    console.error(err)
    this.#loadNewVideo()
  }

  get playing() {
    return !(this.#video.paused || this.#video.ended)
  }

  #nextClip() {
    if (this.#clips.length < 10) {
      this.#fetchClipURLs()
    }

    let next

    if (this.#preloaded.length) {
      // console.log('getting preloaded ObjectURL')
      const nextData = this.#preloaded.pop()
      // console.log('returning ObjectURL for video ' + nextData.url)
      next = nextData.objectURL
    } else {
      // console.log('get next clip URL')
      next = this.#clips.pop()
    }

    this.#preloadNextClip()

    return next
  }

  async #fetchClipURLs() {
    try {
      const resp = await fetch(this.#resourceURL)

      if (!resp.ok) {
        throw new Error(`Response status: ${resp.status}`)
      }

      const data = await resp.json()

      this.#clips = data.clips
      this.#preloadNextClip()
    } catch (err) {
      console.error(err)
    }
  }

  #loadNewVideo() {
    try {
      this.#video.pause()

      this.#videoSource.setAttribute('src', this.#nextClip())

      this.#video.load()
      this.#video.play()
    } catch (err) {
      console.log(`error loading video ${url}`, err)
    }
  }

  async #preloadNextClip() {
    // get the next url
    const url = this.#clips.pop()
    // console.log('start preloading clip ' + url)
    // fetch the video
    try {
      const resp = await fetch(url)
      // get the video data as array buffer
      const data = await resp.arrayBuffer()
      // add to array of preloaded videos
      this.#preloaded.push({
        url,
        objectURL: URL.createObjectURL(
          new Blob([data], { type: 'video/webm' })
        ),
      })
      // console.log('done preloading clip ' + url)
    } catch (err) {
      console.error('error preloading next clip: ', err)
      this.#preloadNextClip()
    }
  }
}
