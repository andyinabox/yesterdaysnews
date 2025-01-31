import { html } from 'lit'
import { component } from 'haunted'

const svgIcon = (name) => {
  return html`
    <svg version="1.1" xmlns="http://www.w3.org/2000/svg" viewbox="0 0 500 500">
      <use href="#yn-icon-${name}" />
    </svg>
  `
}

export function YesterdaysNewNav({ aboutPageUrl, fullscreen }) {
  if (fullscreen) {
    return null
  }

  const onClick = () => {
    this.dispatchEvent(new CustomEvent('fullscreenclick', { bubbles: true }))
  }

  return html`
    <a href=${aboutPageUrl}>${svgIcon('about')}</a>
    <button @click=${onClick}>${svgIcon('fullscreen')}</button>
  `
}

customElements.define(
  'yesterdays-news-nav',
  component(YesterdaysNewNav, {
    observedAttributes: ['about-page-url'],
    useShadowDOM: false,
  })
)
