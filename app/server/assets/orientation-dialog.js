export class OrientationDialog extends HTMLElement {
  constructor() {
    super()
  }

  connectedCallback() {
    // set article content from innerHTML
    const article = document.createElement('article')
    article.innerHTML = this.innerHTML
    this.innerHTML = ''

    // create dialog and append article
    const dialog = document.createElement('dialog')
    dialog.appendChild(article)

    const button = document.createElement('button')
    button.innerHTML = this.closeButtonText
    button.addEventListener('click', () => {
      dialog.close()
    })

    const footer = document.createElement('footer')
    footer.appendChild(button)
    dialog.appendChild(footer)

    this.appendChild(dialog)

    const checkOrientation = () => {
      console.log('checkOrientation', screen.orientation.type)

      if (screen.orientation.type.includes('portrait')) {
        dialog.showModal()
      } else {
        screen.orientation.lock()
        dialog.close()
      }
    }

    screen.orientation.addEventListener('change', checkOrientation)
    checkOrientation()
  }

  get closeButtonText() {
    return this.getAttribute('close-button-label')
  }

  static register() {
    customElements.define('orientation-dialog', OrientationDialog)
  }
}
