import { html } from 'lit'
import { component, useEffect } from 'haunted'

export function YesterdaysNewsPlayer({
  clipsResourceUrl,
  captionsResourceUrl,
  initialClipUrl,
}) {
  return html`
    <yesterdays-news-video
      resource-url=${clipsResourceUrl}
      initial-clip-url=${initialClipUrl}
    ></yesterdays-news-video>
    <yesterdays-news-captions
      resource-url=${captionsResourceUrl}
    ></yesterdays-news-captions>
  `
}

YesterdaysNewsPlayer.observedAttributes = [
  'clips-resource-url',
  'initial-clip-url',
  'captions-resource-url',
]

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
