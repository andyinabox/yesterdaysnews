import { html } from 'lit'
import { createRef, ref } from 'lit/directives/ref.js'
import { component, useRef, useEffect, useState } from 'haunted'
import { svgIcon } from '../lib/svg.js'
import { canAutoplayVideoIfMuted } from '../lib/media.js'
import { usePlaybackPolling } from '../hooks/use-playback-polling.js'
import { VideoLoader } from '../lib/video-loader.js'

export function YesterdaysNewsVideo({
  resourceUrl,
  initialClipUrl,
  fetchClipsWhenLowerThan,
}) {
  const [showPlayButton, setShowPlayButton] = useState(false)
  const [hasPlayedOnce, setHasPlayedOnce] = useState(false)
  const [videoLoader, setVideoLoader] = useState(null)

  // lit-style dom refs
  const videoEl = createRef()
  const sourceEl = createRef()

  const loadNextVideo = async () => {
    try {
      const next = await videoLoader.next()

      if (!videoEl.value) {
        throw new Error('video element is not available')
      }

      videoEl.value.pause()
      sourceEl.value.setAttribute('src', next)
      videoEl.value.load()
      await videoEl.value.play()
    } catch (err) {
      if (err.name !== 'AbortError') {
        console.error(err)
      }
    }
  }

  // this should only be necessary in safari, but doesn't seem to
  // cause any issues in other browsers I've tested
  usePlaybackPolling(videoEl, loadNextVideo)

  // fetch clips on initial load
  useEffect(() => {
    setVideoLoader(new VideoLoader(resourceUrl, fetchClipsWhenLowerThan))
  }, [resourceUrl, fetchClipsWhenLowerThan])

  // by default show video button if autoplay is disabled
  useEffect(() => {
    if (hasPlayedOnce) {
      setShowPlayButton(false)
      return
    }

    if (!videoEl.value) return

    if (!canAutoplayVideoIfMuted(videoEl.value)) {
      setShowPlayButton(true)
    }
  }, [videoEl.value, hasPlayedOnce])

  // event handlers
  const onEnded = () => {
    loadNextVideo()
  }

  const onPlay = () => {
    this.dispatchEvent(new Event('play'))
  }

  const onPlaying = () => {
    setHasPlayedOnce(true)
  }

  const onSourceError = () => {
    loadNextVideo()
  }

  // render functions
  const renderPlayButton = () => {
    const onPlayClick = () => {
      videoEl.value.play()
    }
    if (showPlayButton) {
      return html`<button @click=${onPlayClick} class="play-button">
        ${svgIcon('play')}
      </button>`
    }
  }

  return html`
    <video
      ${ref(videoEl)}
      muted
      autoplay
      playsinline
      tabindex="-1"
      @ended=${onEnded}
      @play=${onPlay}
      @playing=${onPlaying}
    >
      <source
        ${ref(sourceEl)}
        type="video/webm"
        src=${initialClipUrl}
        @error=${onSourceError}
      />
    </video>
    ${renderPlayButton()}
  `
}
customElements.define(
  'yesterdays-news-video',
  component(YesterdaysNewsVideo, {
    observedAttributes: ['resource-url', 'initial-clip-url'],
    useShadowDOM: false,
  })
)
