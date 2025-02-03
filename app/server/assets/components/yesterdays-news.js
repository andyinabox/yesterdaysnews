import { html } from 'lit'
import { component, useState, useEffect } from 'haunted'

import { calcVideoDimensions } from '../lib/math.js'

import { YesterdaysNewsStyles } from './yesterdays-news-styles.js'
import { YesterdaysNewsIcons } from './yesterdays-news-icons.js'

function YesterdaysNews({
  clipsResourceUrl,
  captionsResourceUrl,
  initialClipUrl,
  aboutPageUrl,
}) {
  const [fullscreen, setFullscreen] = useState(false)

  // update video dimensions when window resizes
  const handleResize = () => {
    const { width, height } = calcVideoDimensions(this)
    this.style.setProperty('--yn-video-width', width + 'px')
    this.style.setProperty('--yn-video-height', height + 'px')
    this.style.setProperty('--yn-caption-font-size', width * 0.04 + 'px')
  }
  useEffect(() => {
    window.addEventListener('resize', handleResize)
    handleResize()
    return () => window.removeEventListener('resize', handleResize)
  }, [])

  // handle fullscreen event
  const handleFullScreenChange = () => {
    setFullscreen(document.fullscreenElement === this)
  }
  useEffect(() => {
    document.addEventListener('fullscreenchange', handleFullScreenChange)
    return () =>
      document.removeEventListener('fullscreenchange', handleFullScreenChange)
  }, [])

  // handle change to fullscreen
  useEffect(() => {
    this.classList.toggle('fullscreen', fullscreen)
  }, [fullscreen])

  const onFullscreenClick = () => {
    this.requestFullscreen()
  }

  const renderNav = () => {
    if (fullscreen) {
      return
    }

    return html`<yesterdays-news-nav
      @fullscreenclick=${onFullscreenClick}
      about-page-url=${aboutPageUrl}
      .showFullscreenButton=${!!this.requestFullscreen}
    ></yesterdays-news-nav>`
  }

  return html`
    ${YesterdaysNewsStyles()} ${YesterdaysNewsIcons()}
    <yesterdays-news-player
      clips-resource-url=${clipsResourceUrl}
      initial-clip-url=${initialClipUrl}
      captions-resource-url=${captionsResourceUrl}
    ></yesterdays-news-player>
    ${renderNav()}
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
