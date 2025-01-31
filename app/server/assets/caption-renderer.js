import { SingleObjectCache } from '/assets/single-object-cache.js'

export class CaptionRenderer {
  #resourceURL
  #eventSource

  #caption
  #captionLinesCache

  #rotated = false

  constructor(resourceURL) {
    this.#resourceURL = resourceURL
    this.#captionLinesCache = new SingleObjectCache()
  }

  setup() {
    this.#eventSource = new EventSource(this.#resourceURL)
    this.#eventSource.addEventListener('message', this.handleMessage.bind(this))
  }

  update(rotated = false) {
    this.#rotated = rotated
  }

  draw(ctx, x, y, width, height) {
    // set canvas based on video width
    const fontSize = 0.04 * width
    const fontFamily = 'monospace'
    const verticalPadding = 0.1 * fontSize
    const horizontalPadding = 0.2 * fontSize
    const maxCaptionWidth = width - 4 * horizontalPadding

    const lineData = this.#getCaptionLineData(
      ctx,
      this.#caption,
      maxCaptionWidth,
      {
        rotated: this.#rotated,
        caption: this.#caption,
        fontSize,
        fontFamily,
        maxCaptionWidth,
      }
    )

    // go through the lines of text and draw
    let textWidth, textHeight, textX, textY
    for (let i = 0; i < lineData.length; i++) {
      // update metrics
      const { line, metrics } = lineData[i]

      // calculate bounding box
      textWidth = metrics.width
      textHeight =
        metrics.fontBoundingBoxAscent + metrics.fontBoundingBoxDescent

      if (this.rotated) {
        textX = this.canvas.height / 2 - metrics.width / 2
      } else {
        textX = this.canvas.width / 2 - metrics.width / 2
      }
      textY = y + height - fontSize - (textHeight + 2 * verticalPadding) * i

      // draw black box
      this.ctx.fillStyle = 'rgb(0 0 0)'
      this.ctx.fillRect(
        textX - horizontalPadding,
        textY - verticalPadding - fontSize,
        textWidth + 2 * horizontalPadding,
        textHeight + 2 * verticalPadding
      )

      // draw text
      this.ctx.fillStyle = 'rgb(255 255 255)'
      this.ctx.fillText(line, textX, textY)
    }
  }

  destroy() {
    this.#eventSource.removeEventListener(
      'message',
      this.handleMessage.bind(this)
    )
  }

  handleMessage(event) {
    this.caption = event.data
  }

  set caption(str) {
    this.#caption = str
  }

  get caption() {
    return this.#caption
  }

  #getCaptionLineData(ctx, caption, maxWidth, cachedAttributes = {}) {
    if (this.#captionLinesCache.has(cachedAttributes)) {
      return this.#captionLinesCache.get(cachedAttributes)
    }

    const lineData = []
    const tokens = caption.split(' ')

    let line, metrics
    while (tokens.length) {
      // take next 1-3 words to start the next line
      line = tokens.splice(0, 3).join(' ')
      // get dimensions
      metrics = ctx.measureText(line)

      // keep adding tokens until no tokens are left or the line
      // is larger than maxWidth
      var nextToken, nextLine
      while (tokens.length) {
        // add another word to the line and get metrics
        nextToken = tokens.shift()
        nextLine = line + ' ' + nextToken
        metrics = ctx.measureText(nextLine)

        // if too wide, replace the last word and finish loop
        if (metrics.width > maxWidth) {
          tokens.unshift(nextToken)
          break
        }

        // if not too wide, update value of line and run loop agein
        line = nextLine
      }

      lineData.unshift({ line, metrics })
    }

    this.#captionLinesCache.set(cachedAttributes, lineData)

    return lineData
  }
}
