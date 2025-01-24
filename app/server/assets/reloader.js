class Reloader {
  constructor(reloadUrl) {
    this.eventSource = new EventSource(reloadUrl)
    this.eventSource.addEventListener('message', this.onMessage.bind(this))
  }
  onMessage() {
    this.eventSource.close()
    window.location.reload()
  }
}
