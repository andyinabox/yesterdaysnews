export function canRequestFullscreen(el) {
  return !!el.requestFullscreen
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
