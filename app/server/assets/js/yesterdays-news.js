import { html } from 'https://esm.sh/lit'
import { component, useState, useEffect } from 'https://esm.sh/haunted'

const ASPECT_RATIO = 0.5625

function calcVideoDimensions(el, rotated = true) {
  const { width: containerWidth, height: containerHeight } =
    el.getBoundingClientRect()

  let width, height

  if (rotated) {
    if (containerHeight * ASPECT_RATIO > containerWidth) {
      height = containerWidth
      width = height / ASPECT_RATIO
    } else {
      width = containerHeight
      height = width * ASPECT_RATIO
    }
  } else {
    if (containerWidth * ASPECT_RATIO > containerHeight) {
      height = containerHeight
      width = height / ASPECT_RATIO
    } else {
      width = containerWidth
      height = width * ASPECT_RATIO
    }
  }

  return {
    width,
    height,
  }
}

function YesterdaysNews({
  clipsResourceUrl,
  captionsResourceUrl,
  initialClipUrl,
  aboutPageUrl,
}) {
  const [videoWidth, setVideoWidth] = useState(0)
  const [videoHeight, setVideoHeight] = useState(0)
  const [rotated, setRotated] = useState(false)

  // update video dimensions when window resizes
  useEffect(() => {
    const handleResize = () => {
      const { width, height } = calcVideoDimensions(this, rotated)
      console.log('videoDimensions', width, height)
      setVideoWidth(width)
      setVideoHeight(height)
    }
    window.addEventListener('resize', handleResize)
    handleResize()
    return () => window.removeEventListener('resize', handleResize)
  }, [])

  useEffect(() => {
    if (rotated) {
      this.style.setProperty('--yn-transform', 'rotate(90deg)')
    } else {
      this.style.setProperty('--yn-transform', 'none')
    }
  }, [rotated])

  return html`
    <style>
      :host {
        display: flex;
        align-items: center;
        justify-content: center;
      }
      yesterdays-news-player {
        transform: var(--yn-transform);
      }
    </style>
    <yesterdays-news-player
      clips-resource-url=${clipsResourceUrl}
      initial-clip-url=${initialClipUrl}
      captions-resource-url=${captionsResourceUrl}
      .width=${videoWidth}
      .height=${videoHeight}
    ></yesterdays-news-player>
    <yesterdays-news-nav about-page-url=${aboutPageUrl}></yesterdays-news-nav>
  `
}

customElements.define(
  'yesterdays-news',
  component(YesterdaysNews, {
    observedAttributes: [
      'clips-resource-url',
      'initial-clip-url',
      'captions-resource-url',
      'about-page-url',
    ],
  })
)
