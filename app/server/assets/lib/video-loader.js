import { fetchObjectURL } from '../lib/media.js'

const FETCH_WHEN_LOWER_THAN = 5

export class VideoLoader {
  #resourceUrl = ''

  #clips = []
  #preloaded = []

  #fetching = Promise.resolve()
  #preloading = Promise.resolve()

  constructor(resourceUrl) {
    this.#resourceUrl = resourceUrl
    this.fetch()
  }

  async next() {
    await this.fetching

    if (this.#clips.length < FETCH_WHEN_LOWER_THAN) {
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
      const objectURL = await fetchObjectURL(this.#clips.pop(), 'video/webm')
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
