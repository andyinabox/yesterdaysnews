import { html } from 'lit'
import { component, useState, useEffect } from 'haunted'

import { calcVideoDimensions, canRequestFullscreen } from '../lib/media.js'
import { useConstructableStylesheets } from '../hooks/use-constructable-stylesheets.js'
import { svgSymbols } from '../lib/svg.js'

import styles from './yesterdays-news.styles.js'

function YesterdaysNews({
  clipsResourceUrl,
  captionsResourceUrl,
  initialClipUrl,
  aboutPageUrl,
  fetchClipsWhenLowerThan,
}) {
  const [fullscreen, setFullscreen] = useState(false)

  useConstructableStylesheets(this, [styles])

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

  const onFullscreenClick = () => {
    this.requestFullscreen()
  }

  const renderNav = () => {
    // hide nav in fullscreen mode
    if (fullscreen) {
      return
    }

    return html`<yesterdays-news-nav
      @fullscreenclick=${onFullscreenClick}
      about-page-url=${aboutPageUrl}
      .showFullscreenButton=${canRequestFullscreen(this)}
    ></yesterdays-news-nav>`
  }

  return html`
    ${svgSymbols()}
    <yesterdays-news-player
      clips-resource-url=${clipsResourceUrl}
      initial-clip-url=${initialClipUrl}
      captions-resource-url=${captionsResourceUrl}
      .fetchClipsWhenLowerThan=${parseInt(fetchClipsWhenLowerThan)}
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
      'fetch-clips-when-lower-than',
    ],
  })
)
