import { html } from 'lit'
import { createRef, ref } from 'lit/directives/ref.js'
import { component, useEffect, useState } from 'haunted'

export function YesterdaysNewsStatic() {
  const canvasEl = createRef()

  useEffect(() => {
    if (!canvasEl.value) return

    let play = true

    const tick = () => {
      const ctx = canvasEl.value.getContext('2d')
      const imageData = ctx.createImageData(ctx.canvas.width, ctx.canvas.height)
      const data = imageData.data

      for (let i = 0; i < data.length; i += 4) {
        const value = Math.random() * 127
        data[i] = value // red
        data[i + 1] = value // green
        data[i + 2] = value // blue
        data[i + 3] = 255 // alpha
      }

      ctx.putImageData(imageData, 0, 0)
      req(tick)
    }

    const req = (fn) => {
      if (play) {
        requestAnimationFrame(tick)
      }
    }

    req(tick)

    return () => {
      play = false
    }
  }, [canvasEl.value])

  return html`<canvas ${ref(canvasEl)}></canvas>`
}

customElements.define(
  'yesterdays-news-static',
  component(YesterdaysNewsStatic, {
    observedAttributes: [],
    useShadowDOM: false,
  })
)
