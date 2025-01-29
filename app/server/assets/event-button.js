export class EventButton extends HTMLElement {
  constructor() {
    super()
  }

  connectedCallback() {
    this.innerHTML = `<svg version="1.1" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 500 500"><use href="#yn-icon-${this.iconName}" /></svg> `

    this.addEventListener('click', this.sendEvent.bind(this))
  }

  sendEvent(originalEvent) {
    this.dispatchEvent(new CustomEvent(this.eventName, { bubbles: true }))
  }

  get eventName() {
    return this.getAttribute('event-name')
  }

  get iconName() {
    return this.getAttribute('icon-name')
  }

  static register() {
    customElements.define('event-button', EventButton)
  }
}
