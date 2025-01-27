;(function () {
  new VideoPlayer('video')

  const body = document.body
  const mainViewer = document.getElementById('main')
  const aboutModal = document.getElementById('about')

  body.addEventListener('yn-open-about', () => {
    aboutModal.showModal()
  })

  body.addEventListener('yn-close-about', () => {
    aboutModal.close()
  })

  body.addEventListener('yn-enter-fullscreen', () => {
    mainViewer.requestFullscreen()
  })

  body.addEventListener('fullscreenchange', () => {
    if (!!document.fullscreenElement) {
      body.classList.add('fullscreen')
    } else {
      body.classList.remove('fullscreen')
    }
  })

  // register reload event source
  const eventSource = new EventSource('/reload')
  eventSource.addEventListener('message', () => {
    eventSource.close()
    window.location.reload()
  })
})()
