import { EventButton } from '/assets/event-button.js'
import { EventIcon } from '/assets/event-icon.js'
import { CaptionDisplay } from '/assets/caption-display.js'
import { VideoPlayer } from '/assets/video-player.js'
import { OrientationDialog } from '/assets/orientation-dialog.js'

// events
const ABOUT_OPEN_EVENT = 'yn-open-about'
const ABOUT_CLOSE_EVENT = 'yn-close-about'
const FULLSCREEN_ENTER_EVENT = 'yn-enter-fullscreen'

// classes
const FULLSCREEN_BODY_CLASS = 'fullscreen'

// ids
const MAIN_VIEWER_ID = 'main'
const ABOUT_MODAL_ID = 'about'

function registerComponents() {
  EventButton.register()
  EventIcon.register()
  CaptionDisplay.register()
  VideoPlayer.register()
  OrientationDialog.register()
}

function registerEventListeners() {
  const mainViewer = document.getElementById(MAIN_VIEWER_ID)
  const aboutModal = document.getElementById(ABOUT_MODAL_ID)

  document.addEventListener(ABOUT_OPEN_EVENT, () => {
    aboutModal.showModal()
  })

  document.addEventListener(ABOUT_CLOSE_EVENT, () => {
    aboutModal.close()
  })

  document.addEventListener(FULLSCREEN_ENTER_EVENT, () => {
    mainViewer.requestFullscreen()
  })

  document.addEventListener('fullscreenchange', () => {
    if (!!document.fullscreenElement) {
      document.body.classList.add(FULLSCREEN_BODY_CLASS)
    } else {
      document.body.classList.remove(FULLSCREEN_BODY_CLASS)
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
}

function registerReloadEventSource() {
  const eventSource = new EventSource('/reload')
  eventSource.addEventListener('message', () => {
    eventSource.close()
    window.location.reload()
  })
}

// go!
;(function () {
  registerComponents()
  registerEventListeners()
  registerReloadEventSource()
})()
