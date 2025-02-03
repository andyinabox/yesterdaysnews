import { html } from 'lit'

export const svgIcon = (name) => {
  return html`
    <svg version="1.1" xmlns="http://www.w3.org/2000/svg" viewbox="0 0 500 500">
      <use href="#yn-icon-${name}" />
    </svg>
  `
}
