export function canRequestFullscreen(el) {
  return !!el.requestFullscreen
}

// https://developer.mozilla.org/en-US/docs/Web/Media/Autoplay_guide
export function canAutoplayVideoIfMuted(videoEl) {
  if (!navigator.getAutoplayPolicy) {
    return false
  }

  const policy = navigator.getAutoplayPolicy(videoEl)
  if (!policy) {
    return false
  }

  console.log('policy', policy)

  return policy === 'allowed' || policy === 'allowed-muted'
}
