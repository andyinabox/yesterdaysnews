import { html } from 'lit'
import { createRef, ref } from 'lit/directives/ref.js'
import { component, useEffect, useState } from 'haunted'

// not sure if this work the way I think in JS, but in
// theory it prevents re-allocating memory every frame?
let imageData, data

export function YesterdaysNewsStatic({ maxBrightness = 255 }) {
  const canvasEl = createRef()

  useEffect(() => {
    if (!canvasEl.value) return

    let play = true

    const tick = () => {
      const ctx = canvasEl.value.getContext('2d')
      imageData = ctx.createImageData(ctx.canvas.width, ctx.canvas.height)
      data = imageData.data

      for (let i = 0; i < data.length; i += 4) {
        const value = Math.random() * maxBrightness
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
