import { html } from 'https://esm.sh/lit'
import { component } from 'https://esm.sh/haunted'

export function YesterdaysNewNav({ aboutPageUrl }) {
  const onClick = () => {
    this.dispatchEvent(new CustomEvent('fullscreenclick', { bubbles: true }))
  }

  return html`
    <style>
      :host {
        display: block;
        position: fixed;
        bottom: 0px;
        left: 0px;
        width: 100%;
      }
      nav {
        display: flex;
        align-items: center;
        justify-content: space-between;
      }
    </style>
    <nav>
      <a href=${aboutPageUrl}>About</a>
      <button @click=${onClick}>Fullscreen</button>
    </nav>
  `
}

customElements.define(
  'yesterdays-news-nav',
  component(YesterdaysNewNav, {
    observedAttributes: ['about-page-url'],
    // useShadowDOM: false,
  })
)
