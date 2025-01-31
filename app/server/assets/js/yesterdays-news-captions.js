import { html } from 'https://esm.sh/lit'
import { component, useState, useEffect, useMemo } from 'https://esm.sh/haunted'

export function YesterdaysNewsCaptions({ resourceUrl, width, height }) {
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

  const wordElements = useMemo(() => {
    return words.map(
      (word, i) => html`<span class="word">
        ${word}${i == words.length - 1 ? '' : ' '}
      </span>`
    )
  }, [words])

  useEffect(() => {
    // this.style.setProperty('--width', width + 'px')
    // this.style.setProperty('--height', height + 'px')
    this.style.setProperty('--font-size', width * 0.04 + 'px')
  }, [width, height])

  return html`
    <style>
      :host {
        width: var(--width);
        height: var(--height);
        font-size: var(--font-size);
        display: flex;
        align-items: flex-end;
        justify-content: center;
      }

      .caption {
        font-family: monospace;
        text-align: center;
        margin-bottom: 0.5em;
      }
      .word {
        background-color: black;
        color: white;
      }
    </style>
    <div class="caption">${wordElements}</div>
  `
}

customElements.define(
  'yesterdays-news-captions',
  component(YesterdaysNewsCaptions, {
    observedAttributes: ['resource-url'],
    // useShadowDOM: false,
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
