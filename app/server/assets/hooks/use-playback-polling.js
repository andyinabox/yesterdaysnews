import { useState, useEffect, useRef } from 'haunted'

const PLAYBACK_FAIL_LIMIT = 1

// usePlaybackPolling will poll the video element for playback status
// NOTE: this does not account for all playback events, only the ones
// relevant to this application
export function usePlaybackPolling(videoLitRef, callback, interval = 200) {
  const [needsReload, setNeedsReload] = useState(false)
  const isPlaying = useRef(false)

  const onEnded = () => {
    isPlaying.current = false
  }

  const onPlaying = () => {
    isPlaying.current = true
  }

  // add listeners to video element
  useEffect(() => {
    if (!videoLitRef.value) return

    videoLitRef.value.addEventListener('ended', onEnded)
    videoLitRef.value.addEventListener('playing', onPlaying)

    return () => {
      videoLitRef.value.removeEventListener('ended', onEnded)
      videoLitRef.value.removeEventListener('playing', onPlaying)
    }
  }, [videoLitRef.value])

  // poll playback
  useEffect(() => {
    let count = 0
    const int = setInterval(() => {
      if (isPlaying.current) {
        count = 0
      } else {
        count++
      }

      if (count > PLAYBACK_FAIL_LIMIT) {
        setNeedsReload(true)
      }
    }, interval)
    return () => clearInterval(int)
  }, [])

  useEffect(() => {
    if (needsReload) {
      callback()
    }
    setNeedsReload(false)
  }, [needsReload])
}
