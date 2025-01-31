import { EventButton } from './js/event-button.js'
import { EventIcon } from './js/event-icon.js'
import { CaptionDisplay } from './js/caption-display.js'
import { VideoPlayer } from './js/video-player.js'

// events
const ABOUT_OPEN_EVENT = 'yn-open-about'
const ABOUT_CLOSE_EVENT = 'yn-close-about'
const ORIENTATION_CLOSE_EVENT = 'yn-close-orientation'
const FULLSCREEN_ENTER_EVENT = 'yn-enter-fullscreen'

// classes
const FULLSCREEN_BODY_CLASS = 'fullscreen'

// ids
const MAIN_VIEWER_ID = 'main'
const ABOUT_MODAL_ID = 'about'
const ORIENTATION_MODAL_ID = 'orientation'

function registerComponents() {
  EventButton.register()
  EventIcon.register()
  CaptionDisplay.register()
  VideoPlayer.register()
}

function isFullScreen() {
  return document.body.classList.contains(FULLSCREEN_BODY_CLASS)
}

function checkOrientation() {
  const orientationModal = document.getElementById(ORIENTATION_MODAL_ID)

  if (!isFullScreen()) {
    orientationModal.close()
    return
  }

  if (screen.orientation.type.includes('portrait')) {
    orientationModal.showModal()
  } else {
    orientationModal.close()
  }
}

function registerEventListeners() {
  const mainViewer = document.getElementById(MAIN_VIEWER_ID)
  const aboutModal = document.getElementById(ABOUT_MODAL_ID)
  const orientationModal = document.getElementById(ORIENTATION_MODAL_ID)

  document.addEventListener(ABOUT_OPEN_EVENT, () => {
    aboutModal.showModal()
  })

  document.addEventListener(ABOUT_CLOSE_EVENT, () => {
    aboutModal.close()
  })

  document.addEventListener(FULLSCREEN_ENTER_EVENT, () => {
    mainViewer.requestFullscreen()
  })

  document.addEventListener(ORIENTATION_CLOSE_EVENT, () => {
    orientationModal.close()
  })

  document.addEventListener('fullscreenchange', () => {
    if (!!document.fullscreenElement) {
      document.body.classList.add(FULLSCREEN_BODY_CLASS)
      checkOrientation()
      // we only want to use the dialog in full screen mode
    } else {
      document.body.classList.remove(FULLSCREEN_BODY_CLASS)
      checkOrientation()
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

  // orientation events
  screen.orientation.addEventListener('change', checkOrientation)
  checkOrientation()
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
