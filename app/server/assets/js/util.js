const passthoughTemplateLiteral = (strings, ...values) =>
  String.raw({ raw: strings }, ...values)

// these are no different then the default `` template literal right now,
// but they allow for easier code highlighting :D
export const css = passthoughTemplateLiteral
export const svg = passthoughTemplateLiteral
