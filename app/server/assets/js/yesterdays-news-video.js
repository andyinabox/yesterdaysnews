import { html } from 'https://esm.sh/lit'
import { createRef, ref } from 'https://esm.sh/lit/directives/ref.js'
import { component, useRef, useEffect } from 'https://esm.sh/haunted'

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

  // we want to avoid re-rendering so i think this will work?
  useEffect(() => {
    videoEl.value.setAttribute('width', width)
    videoEl.value.setAttribute('height', height)
  }, [width, height])

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

// export class YesterdaysNewsVideo extends HTMLElement {
//   clips = []
//   preloaded = []

//   constructor() {
//     super()
//     this.loadClips()
//   }

//   connectedCallback() {
//     // create video element
//     const video = document.createElement('video')
//     video.setAttribute('muted', true)
//     video.setAttribute('autoplay', true)
//     video.setAttribute('tabindex', '-1')
//     video.addEventListener('ended', this.onVideoEnded.bind(this))
//     video.addEventListener('error', this.onVideoError.bind(this))

//     // create source element
//     const videoSource = document.createElement('source')
//     videoSource.setAttribute('type', 'video/webm')
//     videoSource.setAttribute('src', this.initialClipUrl)

//     video.appendChild(videoSource)

//     this.video = video
//     this.videoSource = videoSource
//     this.appendChild(video)
//   }

//   async preloadNextClip() {
//     // get the next url
//     const url = this.clips.pop()
//     // console.log('start preloading clip ' + url)
//     // fetch the video
//     try {
//       const resp = await fetch(url)
//       // get the video data as array buffer
//       const data = await resp.arrayBuffer()
//       // add to array of preloaded videos
//       this.preloaded.push({
//         url,
//         objectURL: URL.createObjectURL(
//           new Blob([data], { type: 'video/webm' })
//         ),
//       })
//       // console.log('done preloading clip ' + url)
//     } catch (err) {
//       console.error('error preloading next clip: ', err)
//       this.preloadNextClip()
//     }
//   }

//   async loadClips() {
//     try {
//       const resp = await fetch(this.resourceUrl)

//       if (!resp.ok) {
//         throw new Error(`Response status: ${resp.status}`)
//       }

//       const data = await resp.json()

//       this.clips = data.clips
//       this.preloadNextClip()
//     } catch (err) {
//       console.error(err)
//     }
//   }

//   onVideoError(err) {
//     console.error(err)
//     this.loadNewVideo(this.getNextClip())
//   }

//   onVideoEnded() {
//     this.loadNewVideo(this.getNextClip())
//   }

//   getNextClip() {
//     if (this.clips.length < 10) {
//       this.loadClips()
//     }

//     let next
//     if (this.preloaded.length) {
//       // console.log('getting preloaded ObjectURL')
//       const nextData = this.preloaded.pop()
//       // console.log('returning ObjectURL for video ' + nextData.url)
//       next = nextData.objectURL
//     } else {
//       // console.log('get next clip URL')
//       next = this.clips.pop()
//     }

//     this.preloadNextClip()

//     return next
//   }

//   loadNewVideo(url) {
//     try {
//       this.video.pause()

//       this.videoSource.setAttribute('src', url)

//       this.video.load()
//       this.video.play()
//     } catch (err) {
//       console.log(`error loading video ${url}`, err)
//     }
//   }

//   get resourceUrl() {
//     return this.getAttribute('resource-url')
//   }

//   get initialClipUrl() {
//     return this.getAttribute('initial-clip-url')
//   }
// }
// customElements.define('yesterdays-news-video', YesterdaysNewsVideo)
