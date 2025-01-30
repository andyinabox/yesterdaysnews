export class CaptionLoader {
  constructor(resourceUrl, callback = function () {}) {
    this.resourceUrl = resourceUrl
    this.callback = callback
  }

  connect() {
    this.eventSource = new EventSource(this.resourceUrl)
    this.eventSource.addEventListener('message', this.#handleMessage.bind(this))
  }

  disconnect() {
    this.eventSource.removeEventListener(
      'message',
      this.#handleMessage.bind(this)
    )
  }

  #handleMessage(event) {
    this.callback(event.data)
  }
}
