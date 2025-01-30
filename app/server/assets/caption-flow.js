export class CaptionFlow {
  #fontSize = 0
  #fontFamily = ''
  #maxWidth = 0
  #text = ''

  #cachedLines = []
  constructor() {}

  get(ctx, fontSize, fontFamily, maxWidth, text) {
    if (
      this.#fontSize === fontSize &&
      this.#fontFamily === fontFamily &&
      this.#maxWidth === maxWidth &&
      this.#text === text
    ) {
      return this.#cachedLines
    }

    ctx.font = `${fontSize}px ${fontFamily}`

    // skip all this if we've cached the sentence
    const lines = []
    const tokens = text.split(' ')

    let currentLine, metrics
    while (tokens.length) {
      // take next 1-3 words to start the next line
      currentLine = tokens.splice(0, 3).join(' ')
      // get dimensions
      metrics = ctx.measureText(currentLine)

      // keep adding tokens until no tokens are left or the line
      // is larger than maxWidth
      var nextToken, nextLine
      while (tokens.length) {
        nextToken = tokens.shift()
        nextLine = currentLine + ' ' + nextToken
        metrics = ctx.measureText(nextLine)

        if (metrics.width > maxWidth) {
          tokens.unshift(nextToken)
          break
        }

        currentLine = nextLine
      }

      lines.unshift(currentLine)
    }

    this.#fontSize = fontSize
    this.#fontFamily = fontFamily
    this.#maxWidth = maxWidth
    this.#text = text

    this.#cachedLines = lines

    return this.#cachedLines
  }
}
