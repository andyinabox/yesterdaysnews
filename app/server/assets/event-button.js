class EventButton extends HTMLElement {
  constructor() {
    super()
  }

  connectedCallback() {
    this.eventName = this.getAttribute('event-name')
    this.addEventListener('click', this.sendEvent.bind(this))
  }

  sendEvent(originalEvent) {
    this.dispatchEvent(new CustomEvent(this.eventName, { bubbles: true }))
  }
}

customElements.define('event-button', EventButton)
