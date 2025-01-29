export class EventButton extends HTMLElement {
  constructor() {
    super()
  }

  connectedCallback() {
    this.addEventListener('click', this.sendEvent.bind(this))
  }

  sendEvent(originalEvent) {
    this.dispatchEvent(new CustomEvent(this.eventName, { bubbles: true }))
  }

  get eventName() {
    return this.getAttribute('event-name')
  }

  static register() {
    customElements.define('event-button', EventButton)
  }
}
