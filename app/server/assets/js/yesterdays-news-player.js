import { html } from 'https://esm.sh/lit'
import { component } from 'https://esm.sh/haunted'

export function YesterdaysNewsPlayer({
  width,
  height,
  clipsResourceUrl,
  captionsResourceUrl,
  initialClipUrl,
}) {
  return html`
    <yesterdays-news-video
      resource-url=${clipsResourceUrl}
      initial-clip-url=${initialClipUrl}
      .width=${width}
      .height=${height}
    ></yesterdays-news-video>
    <yesterdays-news-captions
      resource-url=${captionsResourceUrl}
      .width=${width}
      .heigth=${height}
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
