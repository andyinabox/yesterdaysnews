import { html } from 'lit'
import { createRef, ref } from 'lit/directives/ref.js'
import { component, useRef, useEffect, useState } from 'haunted'
import { svgIcon } from '../lib/svg.js'
import { canAutoplayVideoIfMuted, fetchObjectURL } from '../lib/media.js'

export function YesterdaysNewsVideo({ resourceUrl, initialClipUrl }) {
  const [showPlayButton, setShowPlayButton] = useState(false)
  const [hasPlayedOnce, setHasPlayedOnce] = useState(false)

  // react-style data refs
  const clipsRef = useRef([])
  const preloadedRef = useRef([])

  // lit-style dom refs
  const videoEl = createRef()
  const sourceEl = createRef()

  const preloadNextClip = async () => {
    try {
      const objectURL = await fetchObjectURL(
        clipsRef.current.pop(),
        'video/webm'
      )
      preloadedRef.current.push(objectURL)
    } catch (err) {
      console.error('error preloading next clip', err)
    }
  }

  const changeVideoSource = (url) => {
    try {
      videoEl.value.pause()
      sourceEl.value.setAttribute('src', url)
      videoEl.value.load()
      videoEl.value.play()
    } catch (err) {
      console.error(`error changing video source to ${url}`, err)
    }
  }

  const nextVideo = () => {
    // note: it's possible we could end up with
    // multiple outgoing requests if this is called
    // again before the first is completed
    if (clipsRef.current.length < 10) {
      fetchClips()
    }

    let next
    if (preloadedRef.current.length) {
      next = preloadedRef.current.pop()
    } else {
      next = clipsRef.current.pop()
    }

    preloadNextClip()

    return next
  }

  const fetchClips = async () => {
    try {
      const resp = await fetch(resourceUrl)
      if (!resp.ok) {
        throw new Error(`Response status: ${resp.status}`)
      }
      const data = await resp.json()

      clipsRef.current = data.clips
      preloadedRef.current = []
      preloadNextClip()
    } catch (err) {
      console.error('error loading clip urls', err)
    }
  }

  // fetch clips on initial load
  useEffect(fetchClips, [resourceUrl])

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

  const onEnded = () => {
    changeVideoSource(nextVideo())
  }

  const onError = (err) => {
    console.error(err)
    onEnded()
  }

  const onPlay = () => {
    this.dispatchEvent(new Event('play'))
    setHasPlayedOnce(true)
  }

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
      @error=${onError}
      @play=${onPlay}
    >
      <source ${ref(sourceEl)} type="video/webm" src=${initialClipUrl} />
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
