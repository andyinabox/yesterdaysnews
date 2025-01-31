import { html } from 'https://esm.sh/lit'
import { component } from 'https://esm.sh/haunted'

export function YesterdaysNewNav({ aboutPageUrl }) {
  const onClick = () => {
    this.dispatchEvent(new CustomEvent('fullscreenclick', { bubbles: true }))
  }

  return html`
    <a href=${aboutPageUrl}>About</a>
    <button @click=${onClick}>Fullscreen</button>
  `
}

customElements.define(
  'yesterdays-news-nav',
  component(YesterdaysNewNav, {
    observedAttributes: ['about-page-url'],
    useShadowDOM: false,
  })
)
