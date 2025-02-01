import { html } from 'lit'
import { createRef, ref } from 'lit/directives/ref.js'
import { component, useRef, useEffect } from 'haunted'

const fetchObjectURL = async (url, type) => {
  const resp = await fetch(url)
  // get the data as array buffer
  const data = await resp.arrayBuffer()
  // create ObjectURL from data
  return URL.createObjectURL(new Blob([data], { type }))
}

export function YesterdaysNewsVideo({
  resourceUrl,
  initialClipUrl,
  width,
  height,
}) {
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
      console.log(`error changing video source to ${url}`, err)
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

  const onEnded = () => {
    changeVideoSource(nextVideo())
  }
  const onError = (err) => {
    console.error(err)
    onEnded()
  }

  return html`
    <video
      ${ref(videoEl)}
      muted
      autoplay
      tabindex="-1"
      @ended=${onEnded}
      @error=${onError}
    >
      <source ${ref(sourceEl)} type="video/webm" src=${initialClipUrl} />
    </video>
  `
}
customElements.define(
  'yesterdays-news-video',
  component(YesterdaysNewsVideo, {
    observedAttributes: ['resource-url', 'initial-clip-url'],
    useShadowDOM: false,
  })
)
