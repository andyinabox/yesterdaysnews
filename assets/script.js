const YN_OBJECT_STORE_URL = window.YN_ENV.objectStoreUrl

const manifestUrl = `${YN_OBJECT_STORE_URL}/manifest.json`

async function getManifest() {
  try {
    const response = await fetch(manifestUrl)
    if (!response.ok) {
      throw new Error(`Response status: ${response.status}`)
    }

    const json = await response.json()
    return json
  } catch (error) {
    console.error(error.message)
  }
}
async function run() {
  // handle captions
  const caption = document.getElementById('caption')
  const es = new EventSource('/captions')
  es.addEventListener('message', (evt) => {
    // console.log(evt.data)
    caption.innerText = evt.data
  })

  // handle videos
  const manifest = await getManifest()
  console.log(manifest)
  var clips = []

  function getNextVideoURL() {
    if (clips.length === 0) {
      clips = manifest.files.clips
        .map((path) => `${YN_OBJECT_STORE_URL}/${path}`)
        .sort(() => Math.random() * 0.5)
    }
    const n = Math.floor(Math.random() * clips.length)
    return clips.splice(n, 1)
  }

  const videoEl = document.getElementById('video')
  const webmSourceEl = document.getElementById('video-src-webm')

  function loadNewVideo() {
    try {
      videoEl.pause()

      const url = getNextVideoURL()

      webmSourceEl.setAttribute('type', 'video/webm')
      webmSourceEl.setAttribute('src', url)

      videoEl.load()
      videoEl.play()
    } catch (err) {
      console.log(`error loading video ${url}`, err)
      loadNewVideo()
    }
  }

  videoEl.addEventListener('ended', (evt) => {
    // console.log('video ended', evt)
    loadNewVideo()
  })
  videoEl.addEventListener('error', (evt) => {
    console.log('video error', evt)
    loadNewVideo()
  })

  loadNewVideo()
}

run()
