import { fetchObjectURL } from '../lib/media.js'

export class VideoLoader {
  #resourceUrl = ''
  #fetchWhenLowerThan = 5

  #clips = []
  #preloaded = []

  #fetching = Promise.resolve()
  #preloading = Promise.resolve()

  constructor(resourceUrl, fetchWhenLowerThan = 5) {
    this.#resourceUrl = resourceUrl
    this.#fetchWhenLowerThan = fetchWhenLowerThan
    this.fetch()
  }

  async next() {
    await this.fetching

    if (this.#clips.length < this.#fetchWhenLowerThan) {
      this.fetch()
    }

    let next
    if (this.#preloaded.length) {
      next = this.#preloaded.pop()
    } else {
      next = this.#clips.pop()
    }

    this.#preload()

    return next
  }

  fetch() {
    // want to set the fetch promise before returning
    this.#fetching = (async () => {
      const resp = await fetch(this.#resourceUrl)
      if (!resp.ok) {
        throw new Error(`Response status: ${resp.status}`)
      }
      const data = await resp.json()

      this.#clips = data.clips

      // await this.preloading
      this.#preload()
    })()
  }

  #preload() {
    this.#preloading = (async () => {
      await this.fetching

      if (!this.#clips.length) {
        throw new Error('attempted to preload with no clips')
      }
      const clipUrl = this.#clips.pop()
      const mimeType = clipUrl.endsWith('.mp4') ? 'video/mp4' : 'video/webm'
      const objectURL = await fetchObjectURL(clipUrl, mimeType)
      this.#preloaded.push(objectURL)
    })()
  }

  get fetching() {
    return this.#fetching
  }

  get preloading() {
    return this.#preloading
  }
}
