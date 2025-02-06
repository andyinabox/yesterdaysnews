import { html } from 'lit'
import { component, useState } from 'haunted'

export function YesterdaysNewsPlayer({
  clipsResourceUrl,
  captionsResourceUrl,
  initialClipUrl,
  fetchClipsWhenLowerThan,
}) {
  const [showCaptions, setShowCaptions] = useState(false)

  const onVideoPlay = () => {
    setShowCaptions(true)
  }

  const renderCaptions = () => {
    if (!showCaptions) return
    return html`<yesterdays-news-captions
      resource-url=${captionsResourceUrl}
    ></yesterdays-news-captions>`
  }

  return html`
    <yesterdays-news-video
      @play=${onVideoPlay}
      resource-url=${clipsResourceUrl}
      initial-clip-url=${initialClipUrl}
      .fetchClipsWhenLowerThan=${fetchClipsWhenLowerThan}
    ></yesterdays-news-video>
    ${renderCaptions()}
  `
}

customElements.define(
  'yesterdays-news-player',
  component(YesterdaysNewsPlayer, {
    observedAttributes: [
      'clips-resource-url',
      'initial-clip-url',
      'captions-resource-url',
    ],
    useShadowDOM: false,
  })
)
