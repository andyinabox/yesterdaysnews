class Main {
  constructor(elementID) {
    this.el = document.getElementById(elementID)
    this.aboutBtn = document.getElementById('btn-about')
    this.fullscreenBtn = document.getElementById('btn-fullscreen')

    this.el.addEventListener(
      'fullscreenchange',
      this.onFullscreenChange.bind(this)
    )

    this.aboutBtn.addEventListener('click', this.onAboutBtnClick.bind(this))
    this.fullscreenBtn.addEventListener(
      'click',
      this.onFullscreenBtnClick.bind(this)
    )
  }

  onFullscreenChange(evt) {
    if (!!document.fullscreenElement) {
      this.el.classList.add('fullscreen')
    } else {
      this.el.classList.remove('fullscreen')
    }
  }

  onAboutBtnClick(evt) {
    console.log('about button click')
  }

  onFullscreenBtnClick(evt) {
    console.log('fullscreen button click')
    this.el.requestFullscreen()
  }
}
