class Main {
  constructor(
    mainElementID,
    aboutDialogID,
    aboutBtnID,
    closeAboutBtnID,
    fullscreenBtnID
  ) {
    this.body = document.body
    this.main = document.getElementById(mainElementID)
    this.about = document.getElementById(aboutDialogID)
    this.aboutBtn = document.getElementById(aboutBtnID)
    this.closeAboutBtn = document.getElementById(closeAboutBtnID)
    this.fullscreenBtn = document.getElementById(fullscreenBtnID)

    this.main.addEventListener(
      'fullscreenchange',
      this.onFullscreenChange.bind(this)
    )

    this.aboutBtn.addEventListener('click', this.onAboutBtnClick.bind(this))
    this.closeAboutBtn.addEventListener(
      'click',
      this.onCloseAboutBtnClick.bind(this)
    )

    this.fullscreenBtn.addEventListener(
      'click',
      this.onFullscreenBtnClick.bind(this)
    )
  }

  onFullscreenChange(evt) {
    if (!!document.fullscreenElement) {
      this.body.classList.add('fullscreen')
    } else {
      this.body.classList.remove('fullscreen')
    }
  }

  onAboutBtnClick(evt) {
    this.about.showModal()
  }

  onCloseAboutBtnClick(evt) {
    console.log('close about button click')
    this.about.close()
  }

  onFullscreenBtnClick(evt) {
    console.log('fullscreen button click')
    this.main.requestFullscreen()
  }
}
