export class CaptionDisplay extends HTMLElement {
  constructor() {
    super()
  }
  connectedCallback() {
    this.eventSource = new EventSource(this.resourceUrl)
    this.eventSource.addEventListener('message', this.handleMessage.bind(this))
  }
  handleMessage(evt) {
    const words = evt.data.split(' ')
    const elements = words.map(
      (w, i) =>
        `<span class="word">${w}${i == words.length - 1 ? '' : '&nbsp;'}</span>`
    )
    this.innerHTML = elements.join('')
  }
  disconnectedCallback() {
    this.eventSource.removeEventListener(
      'message',
      this.handleMessage.bind(this)
    )
  }

  get resourceUrl() {
    return this.getAttribute('resource-url')
  }

  static register() {
    customElements.define('caption-display', CaptionDisplay)
  }
}
