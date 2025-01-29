import { EventButton } from '/assets/event-button.js'
import { EventIcon } from '/assets/event-icon.js'
import { CaptionDisplay } from '/assets/caption-display.js'
import { VideoPlayer } from '/assets/video-player.js'
;(function () {
  EventButton.register()
  EventIcon.register()
  CaptionDisplay.register()
  VideoPlayer.register()

  const mainViewer = document.getElementById('main')
  const aboutModal = document.getElementById('about')

  document.addEventListener('yn-open-about', () => {
    aboutModal.showModal()
  })

  document.addEventListener('yn-close-about', () => {
    aboutModal.close()
  })

  document.addEventListener('yn-enter-fullscreen', () => {
    mainViewer.requestFullscreen()
  })

  document.addEventListener('fullscreenchange', () => {
    if (!!document.fullscreenElement) {
      document.body.classList.add('fullscreen')
    } else {
      document.body.classList.remove('fullscreen')
    }
  })

  document.addEventListener('keydown', (event) => {
    if (
      document.activeElement &&
      document.activeElement.getAttribute('role') === 'button'
    ) {
      if (event.key === 'Enter') {
        event.preventDefault()
        document.activeElement.click()
      }
    }
  })

  // register reload event source
  const eventSource = new EventSource('/reload')
  eventSource.addEventListener('message', () => {
    eventSource.close()
    window.location.reload()
  })
})()
