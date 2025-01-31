import { EventElement } from './event-element.js'
export class EventButton extends EventElement {
  constructor() {
    super()
  }

  connectedCallback() {
    const button = document.createElement('button')
    button.innerHTML = this.innerHTML

    this.innerHTML = ''
    if (this.hasAttribute('autofocus')) {
      this.removeAttribute('autofocus')
      button.toggleAttribute('autofocus', true)
    }

    this.appendChild(button)

    button.addEventListener('click', this.sendEvent.bind(this))
  }

  static register() {
    customElements.define('event-button', EventButton)
  }
}
