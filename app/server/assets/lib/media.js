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

export const fetchObjectURL = async (url, type) => {
  const resp = await fetch(url)
  // get the data as array buffer
  const data = await resp.arrayBuffer()
  // create ObjectURL from data
  return URL.createObjectURL(new Blob([data], { type }))
}

export function canAutoplayVideoIfMuted(videoEl) {
  // in this case we have no way of knowing so assume no
  if (!navigator.getAutoplayPolicy) {
    return false
  }

  // https://developer.mozilla.org/en-US/docs/Web/Media/Autoplay_guide
  const policy = navigator.getAutoplayPolicy(videoEl)
  if (!policy) {
    return false
  }

  return policy === 'allowed' || policy === 'allowed-muted'
}

export function canRequestFullscreen(el) {
  return !!el.requestFullscreen
}
