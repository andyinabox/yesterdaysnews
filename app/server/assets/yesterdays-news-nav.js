import svgSymbols from '/assets/svg-symbols.js'
export class YesterdaysNewsNav extends HTMLElement {
  static observedAttributes = ['rotated', 'fullscreen']
  constructor() {
    super()
  }
  connectedCallback() {
    // TODO: add SVG symbols

    this.svgSymbols = document.createElement('div')
    this.svgSymbols.setAttribute('style', 'display: none;')
    this.svgSymbols.innerHTML = svgSymbols
    this.appendChild(this.svgSymbols)

    // about button
    this.aboutBtn = document.createElement('a')
    this.aboutBtn.setAttribute('target', '_blank')
    this.aboutBtn.href = this.aboutURL
    this.aboutBtn.innerHTML = this.#svgIconString('about')

    // fullscreen button
    this.fullscreenBtn = document.createElement('button')
    this.fullscreenBtn.addEventListener('click', this.handleFullscreenClick)
    this.fullscreenBtn.innerHTML = this.#svgIconString('fullscreen')

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
      new CustomEvent('yn-fullscreenrequest', { bubbles: true, composed: true })
    )
  }

  get rotated() {
    return (
      this.hasAttribute('rotated') && this.getAttribute('rotated') !== 'false'
    )
  }

  get fullscreen() {
    return (
      this.hasAttribute('fullscreen') &&
      this.getAttribute('fullscreen') !== 'false'
    )
  }

  get aboutURL() {
    return this.getAttribute('about-url')
  }

  #svgIconString(name) {
    return `
      <svg version="1.1" xmlns="http://www.w3.org/2000/svg" viewbox="0 0 500 500">
        <use href="#yn-icon-${name}" />
      </svg>
    `
  }
}

customElements.define('yesterdays-news-nav', YesterdaysNewsNav)
