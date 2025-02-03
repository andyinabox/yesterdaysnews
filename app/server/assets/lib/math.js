const ASPECT_RATIO = 0.5625

export function calcVideoDimensions(el) {
  const { width: containerWidth, height: containerHeight } =
    el.getBoundingClientRect()

  let width, height

  if (containerWidth * ASPECT_RATIO > containerHeight) {
    height = containerHeight
    width = height / ASPECT_RATIO
  } else {
    width = containerWidth
    height = width * ASPECT_RATIO
  }

  return {
    width,
    height,
  }
}
