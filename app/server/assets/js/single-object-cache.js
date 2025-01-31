export class SingleObjectCache {
  #key
  #value

  constructor() {}

  has(obj = {}) {
    const key = this.#objectToCacheKey(obj)
    return key === this.#key
  }

  get(obj = {}) {
    if (this.has(obj)) {
      return this.#value
    }
    return null
  }

  set(key = {}, value) {
    this.#key = this.#objectToCacheKey(key)
    this.#value = value
  }

  #objectToCacheKey(obj = {}) {
    let cacheKey = ''
    Object.keys(obj)
      .sort()
      .forEach((key) => {
        cacheKey += `${obj[key]}`
      })
    return cacheKey
  }
}
