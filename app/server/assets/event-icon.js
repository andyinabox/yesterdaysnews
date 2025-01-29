import { EventElement } from '/assets/event-element.js'

const VIEWBOX_WIDTH = 500
const VIEWBOX_HEIGHT = 500

export class EventIcon extends EventElement {
  constructor() {
    super()
  }

  connectedCallback() {
    this.innerHTML = `<svg version="1.1" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 ${VIEWBOX_WIDTH} ${VIEWBOX_HEIGHT}"><use href="#yn-icon-${this.iconName}" /></svg>`
    this.setAttribute('tabindex', '0')
    this.setAttribute('role', 'button')
    this.addEventListener('click', this.sendEvent.bind(this))
  }

  get iconName() {
    return this.getAttribute('icon-name')
  }

  static register() {
    customElements.define('event-icon', EventIcon)
  }
}
