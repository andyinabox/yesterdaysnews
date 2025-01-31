import { html } from 'https://esm.sh/lit'
import { component, useEffect } from 'https://esm.sh/haunted'

export function YesterdaysNewsPlayer({
  width,
  height,
  clipsResourceUrl,
  captionsResourceUrl,
  initialClipUrl,
}) {
  useEffect(() => {
    this.style.setProperty('--width', width + 'px')
    this.style.setProperty('--height', height + 'px')
  }, [width, height])

  return html`
    <style>
      :host {
        position: relative;
        width: var(--width);
        height: var(--height);
      }
      yesterdays-news-video {
        /* position: absolute;
        top: 0px;
        left: 0px; */
        z-index: 0;
      }
      yesterdays-news-captions {
        position: absolute;
        top: 0px;
        left: 0px;
        z-index: 1;
      }
    </style>
    <yesterdays-news-video
      resource-url=${clipsResourceUrl}
      initial-clip-url=${initialClipUrl}
      .width=${width}
      .height=${height}
    ></yesterdays-news-video>
    <yesterdays-news-captions
      resource-url=${captionsResourceUrl}
      .width=${width}
      .height=${height}
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
    // useShadowDOM: false,
  })
)
