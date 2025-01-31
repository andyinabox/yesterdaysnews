export class EventElement extends HTMLElement {
  sendEvent() {
    this.dispatchEvent(new CustomEvent(this.eventName, { bubbles: true }))
  }

  get eventName() {
    return this.getAttribute('event-name')
  }
}
