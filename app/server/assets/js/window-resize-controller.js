export class WindowResizeController {
  host
  parentEl
  width
  height

  _onWindowResize = () => {
    this.host.requestUpdate()
  }

  constructor(host, parentEl) {
    this.host = host
    host.addController(this)

    this.parentEl = parentEl
  }

  hostConnected() {
    window.addEventListener('resize', this._onMouseMove)
  }

  hostDisconnected() {
    window.removeEventListener('resize', this._onMouseMove)
  }
}
