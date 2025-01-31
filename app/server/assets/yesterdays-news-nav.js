export class YesterdaysNewsNav extends HTMLElement {
  static observedAttributes = ['rotated', 'fullscreen']
  constructor() {
    super()
  }
  connectedCallback() {
    // TODO: add SVG symbols

    // about button
    this.aboutBtn = document.createElement('a')
    this.aboutBtn.setAttribute('target', '_blank')
    this.aboutBtn.href = this.aboutURL
    this.aboutBtn.innerHTML = 'about'

    // fullscreen button
    this.fullscreenBtn = document.createElement('button')
    this.fullscreenBtn.addEventListener('click', this.handleFullscreenClick)
    this.fullscreenBtn.innerHTML = 'fullscreen'

    // add to dom
    this.appendChild(this.aboutBtn)
    this.appendChild(this.fullscreenBtn)
  }

  disconnectedCallback() {
    this.fullscreenBtn.removeEventListener('click', this.handleFullscreenClick)
  }

  attributeChangedCallback(name, oldValue, newValue) {
    switch (name) {
      case 'rotated':
        break
      case 'fullscreen':
        break
    }
  }

  handleFullscreenClick() {
    this.dispatchEvent(
      new CustomEvent('yn-fullscreenrequest', { bubbles: true })
    )
  }

  get rotated() {
    return this.hasAttribute('rotated') && this.getAttribute('rotated') !== null
  }

  get fullscreen() {
    return (
      this.hasAttribute('fullscreen') &&
      this.getAttribute('fullscreen') !== null
    )
  }

  get aboutURL() {
    return this.getAttribute('about-url')
  }
}

customElements.define('yesterdays-news-nav', YesterdaysNewsNav)
