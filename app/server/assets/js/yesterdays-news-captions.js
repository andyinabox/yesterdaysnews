import { html } from 'https://esm.sh/lit'
import { component, useState, useEffect } from 'https://esm.sh/haunted'

export function YesterdaysNewsCaptions({ resourceUrl }) {
  const [words, setWords] = useState([])

  // add event listener for server-sent event
  useEffect(() => {
    const handleMessage = (event) => {
      console.log('caption: ', event.data)
      setWords(event.data.split(' '))
    }

    const eventSource = new EventSource(resourceUrl)
    eventSource.addEventListener('message', handleMessage)

    return () => eventSource.removeEventListener('message', handleMessage)
  }, [resourceUrl])

  return html`
    ${words.map((word, i) => {
      html`<span class="word">
        ${word}${i == words.length - 1 ? '' : '&nbsp;'}
      </span>`
    })}
  `
}

customElements.define(
  'yesterdays-news-captions',
  component(YesterdaysNewsCaptions, {
    observedAttributes: ['resource-url'],
    useShadowDOM: false,
  })
)

// export class YesterdaysNewsCaptions extends HTMLElement {
//   constructor() {
//     super()
//   }
//   connectedCallback() {
//     this.eventSource = new EventSource(this.resourceUrl)
//     this.eventSource.addEventListener('message', this.handleMessage.bind(this))
//   }
//   handleMessage(evt) {
//     const words = evt.data.split(' ')
//     const elements = words.map(
//       (w, i) =>
//         `<span class="word">${w}${i == words.length - 1 ? '' : '&nbsp;'}</span>`
//     )
//     this.innerHTML = elements.join('')
//   }
//   disconnectedCallback() {
//     this.eventSource.removeEventListener(
//       'message',
//       this.handleMessage.bind(this)
//     )
//   }

//   get resourceUrl() {
//     return this.getAttribute('resource-url')
//   }
// }
