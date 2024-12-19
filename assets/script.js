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

  // webmSourceEl.setAttribute('type', 'video/webm')
  // webmSourceEl.setAttribute('src', getNextVideoURL())

  // videoEl.load()
  // videoEl.play()

  function loadNewVideo() {
    videoEl.pause()

    webmSourceEl.setAttribute('type', 'video/webm')
    webmSourceEl.setAttribute('src', getNextVideoURL())

    videoEl.load()
    videoEl.play()

    console.log({
      src: webmSourceEl.getAttribute('src'),
      type: webmSourceEl.getAttribute('type'),
    })
  }

  videoEl.addEventListener('ended', loadNewVideo)
  videoEl.addEventListener('error', loadNewVideo)
  loadNewVideo()
}

run()
