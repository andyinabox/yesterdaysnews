import { html } from 'lit'
import { component } from 'haunted'
import { svgIcon } from '../lib/svg.js'

export function YesterdaysNewNav({ aboutPageUrl, showFullscreenButton }) {
  const onFullscreenClick = () => {
    this.dispatchEvent(new CustomEvent('fullscreenclick', { bubbles: true }))
  }

  const renderFullscreenButton = () => {
    if (!showFullscreenButton) return
    return html` <button
      class="btn btn-fullscreen"
      tabindex="0"
      title="Enter fullscreen"
      @click=${onFullscreenClick}
    >
      ${svgIcon('fullscreen')}
    </button>`
  }

  return html`
    <a
      class="btn btn-about"
      tabindex="0"
      title="Go to about page"
      href=${aboutPageUrl}
      >${svgIcon('about')}
    </a>
    ${renderFullscreenButton()}
  `
}

customElements.define(
  'yesterdays-news-nav',
  component(YesterdaysNewNav, {
    observedAttributes: ['about-page-url'],
    useShadowDOM: false,
  })
)
