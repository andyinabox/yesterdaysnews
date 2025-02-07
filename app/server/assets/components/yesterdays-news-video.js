import { html } from 'lit'
import { createRef, ref } from 'lit/directives/ref.js'
import { component, useEffect, useState, useRef } from 'haunted'
import { svgIcon } from '../lib/svg.js'
import { canAutoplayVideoIfMuted } from '../lib/media.js'
import { VideoLoader } from '../lib/video-loader.js'

const PLAYBACK_POLLING_INTERVAL = 200
const PLAYBACK_FAIL_LIMIT = 1

export function YesterdaysNewsVideo({
  resourceUrl,
  initialClipUrl,
  fetchClipsWhenLowerThan,
}) {
  // used to hide play button after first play
  const [hasPlayedOnce, setHasPlayedOnce] = useState(false)
  // stores a VideoLoader instance. maybe this should be a ref?
  const [videoLoader, setVideoLoader] = useState(null)
  // indicates we need to load a new video
  const [needsReload, setNeedsReload] = useState(false)
  // determines whether we need to show the play button
  const [showPlayButton, setShowPlayButton] = useState(false)

  // react-style refs
  const isPlaying = useRef(false)

  // lit-style dom refs
  const videoEl = createRef()
  const sourceEl = createRef()

  // get next video from VideoLoader and set as next video
  const loadNextVideo = async () => {
    try {
      const next = await videoLoader.next()

      if (!videoEl.value) return

      videoEl.value.pause()
      sourceEl.value.setAttribute('src', next)
      videoEl.value.load()

      // play returns a promise, and we need to use
      // await here to catch any playback errors
      await videoEl.value.play()
    } catch (err) {
      // abort error happens when we cancel playback because
      // the video was unable to load. this happens normally
      // when recovering from a decoding error so we want
      // to suppress it here
      if (err.name !== 'AbortError') {
        console.error(err)
      }
    }
  }

  //
  // side-effects
  //

  // create new VideoLoader
  useEffect(() => {
    setVideoLoader(new VideoLoader(resourceUrl, fetchClipsWhenLowerThan))
  }, [resourceUrl, fetchClipsWhenLowerThan])

  // by default show video button if autoplay is disabled
  useEffect(() => {
    // once it's played we shouldn't need the button
    if (hasPlayedOnce) {
      setShowPlayButton(false)
      return
    }

    if (!videoEl.value) return

    // only show if autoplay-ability is disabled (or unknown)
    if (!canAutoplayVideoIfMuted(videoEl.value)) {
      setShowPlayButton(true)
    }
  }, [videoEl.value, hasPlayedOnce])

  // poll video to see if video is actually playing
  useEffect(() => {
    let count = 0
    const int = setInterval(() => {
      if (isPlaying.current) {
        count = 0
      } else {
        count++
      }

      if (count > PLAYBACK_FAIL_LIMIT) {
        emitStopped()
        setNeedsReload(true)
      }
    }, PLAYBACK_POLLING_INTERVAL)
    return () => clearInterval(int)
  }, [])

  // respond to needsReload and load next video
  useEffect(() => {
    if (needsReload) loadNextVideo()
    setNeedsReload(false)
  }, [needsReload])

  //
  // event emitters
  //
  const emitPlaying = () => {
    this.dispatchEvent(new Event('playing'))
  }

  const emitStopped = () => {
    this.dispatchEvent(new Event('stopped'))
  }

  //
  // event handlers
  //

  // video has ended
  // https://developer.mozilla.org/en-US/docs/Web/API/HTMLMediaElement/ended_event
  const onEnded = () => {
    isPlaying.current = false
    // note that we aren not emitting the "stopped" event
    // because during normal playback another video should
    // play immediately
    setNeedsReload(true)
  }

  // video has started autplay, or the .play method is triggered
  // https://developer.mozilla.org/en-US/docs/Web/API/HTMLMediaElement/play_event
  const onPlay = () => {}

  // this is fired when the video has *actually* started playing, whereas "play"
  // indicates playback has been requested.
  // https://developer.mozilla.org/en-US/docs/Web/API/HTMLMediaElement/playing_event
  const onPlaying = () => {
    isPlaying.current = true
    setHasPlayedOnce(true)
    emitPlaying()
  }

  // this event is fired on the <source> element when the video cannot be loaded
  // note that it does not seem to work on Safari, or at least some versions of Safari
  // https://developer.mozilla.org/en-US/docs/Web/API/HTMLMediaElement/playing_event
  const onSourceError = () => {
    isPlaying.current = false
    emitStopped()
    loadNextVideo()
  }

  //
  // rendering
  //
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
