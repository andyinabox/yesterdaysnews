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

  return html` <div class="caption">${wordElements}</div> `
}

customElements.define(
  'yesterdays-news-captions',
  component(YesterdaysNewsCaptions, {
    observedAttributes: ['resource-url'],
    useShadowDOM: false,
  })
)
