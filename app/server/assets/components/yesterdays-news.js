import { html } from 'lit'
import { component, useState, useEffect } from 'haunted'

const ASPECT_RATIO = 0.5625

function calcVideoDimensions(el, rotated = false) {
  const { width: containerWidth, height: containerHeight } =
    el.getBoundingClientRect()

  let width, height

  if (rotated) {
    if (containerHeight * ASPECT_RATIO > containerWidth) {
      height = containerWidth
      width = height / ASPECT_RATIO
    } else {
      width = containerHeight
      height = width * ASPECT_RATIO
    }
  } else {
    if (containerWidth * ASPECT_RATIO > containerHeight) {
      height = containerHeight
      width = height / ASPECT_RATIO
    } else {
      width = containerWidth
      height = width * ASPECT_RATIO
    }
  }

  return {
    width,
    height,
  }
}

function YesterdaysNews({
  clipsResourceUrl,
  captionsResourceUrl,
  initialClipUrl,
  aboutPageUrl,
}) {
  const [rotated, setRotated] = useState(false)
  const [fullscreen, setFullscreen] = useState(false)

  // update video dimensions when window resizes
  const handleResize = () => {
    const { width, height } = calcVideoDimensions(this, rotated)
    this.style.setProperty('--yn-video-width', width + 'px')
    this.style.setProperty('--yn-video-height', height + 'px')
    this.style.setProperty('--yn-caption-font-size', width * 0.04 + 'px')
  }
  useEffect(() => {
    window.addEventListener('resize', handleResize)
    handleResize()
    return () => window.removeEventListener('resize', handleResize)
  }, [])

  // handle fullscreen event
  const handleFullScreenChange = () => {
    if (document.fullscreenElement === this) {
      setFullscreen(true)
      // we only want to use the dialog in full screen mode
    } else {
      setFullscreen(false)
    }
  }
  useEffect(() => {
    document.addEventListener('fullscreenchange', handleFullScreenChange)
    return () =>
      document.removeEventListener('fullscreenchange', handleFullScreenChange)
  }, [])

  function handleOrientationChange() {
    // we only want to do rotation in fullscreen mode
    if (!fullscreen) {
      setRotated(false)
      return
    }

    setRotated(screen.orientation.type.includes('portrait'))
  }
  useEffect(() => {
    screen.orientation.addEventListener('change', handleOrientationChange)
    return () =>
      screen.orientation.removeEventListener('change', handleOrientationChange)
  }, [])

  // handle change to fullscreen
  useEffect(() => {
    this.classList.toggle('fullscreen', fullscreen)
  }, [fullscreen])

  // handle change to rotated
  useEffect(() => {
    this.classList.toggle('rotated', rotated)
    this.style.setProperty('--yn-transform', rotated ? 'rotate(90deg)' : 'none')
  }, [rotated])

  useEffect(handleOrientationChange, [fullscreen, rotated])

  const onFullscreenClick = () => {
    this.requestFullscreen()
  }

  return html`
    <style>
        :host {
          /* set defaults */
          --yn-transform: none;
          --yn-video-width: 0px;
          --yn-video-height: 0px
          --yb-font-size: 0px;

          display: flex;
          align-items: center;
          justify-content: center;
        }
        #svg-symbols {
          display: none;
        }
        yesterdays-news-player {
          position: relative;
          transform: var(--yn-transform);
          width: var(--yn-video-width);
          height: var(--yn-video-height);
        }
        yesterdays-news-video {
          display: block;
          z-index: 0;
        }
        yesterdays-news-video > video {
          width: var(--yn-video-width);
          height: var(--yn-video-height);
        }
        yesterdays-news-captions {
          position: absolute;
          top: 0px;
          left: 0px;
          z-index: 1;
          width: var(--yn-video-width);
          height: var(--yn-video-height);
          font-size: var(--yn-caption-font-size);
          display: flex;
          align-items: flex-end;
          justify-content: center;
        }
        yesterdays-news-captions > .caption {
          font-family: monospace;
          text-align: center;
          margin-bottom: 0.5em;
        }
        yesterdays-news-captions > .caption > .word {
          background-color: black;
          color: white;
        }

        yesterdays-news-nav {
          z-index: 3;
          position: fixed;
          bottom: 0px;
          left: 0px;
          width: 100%;
          display: flex;
          align-items: center;
          justify-content: space-between;
        }

        yesterdays-news-nav > a, yesterdays-news-nav > button {
        cursor: pointer;
        font-size: calc(75% + 3vmin);
        display: block;
        border: none;
        background-color: transparent;
        width: 2em;
        height: 2em;
        margin: 0.5em;
      }
      yesterdays-news.fullscreen > yesterdays-news-nav {
        display: none;
      }

      yesterdays-news-nav svg {
        fill: white;
        stroke: none;
        display: block;
        width: 100%;
        height: 100%;
      }
    </style>
    <svg id="svg-symbols" xmlns="http://www.w3.org/2000/svg">
      <defs>
        <!-- https://thenounproject.com/icon/fullscreen-6938590/ -->
        <symbol id="yn-icon-fullscreen">
          <path
            d="m 322.47037,87.302654 h 75.22698 v 75.226976 h 51.13133 V 36.171317 H 322.47037 Z"
          />
          <path
            d="M 92.302654,162.52963 V 87.302654 H 167.52963 V 36.171317 H 41.171317 V 162.52963 Z"
          />
          <path
            d="M 167.52963,392.69735 H 92.302654 V 317.47037 H 41.171317 V 443.82868 H 167.52963 Z"
          />
          <path
            d="m 397.69735,317.47037 v 75.22698 h -75.22698 v 51.13133 H 448.82868 V 317.47037 Z"
          />
        </symbol>

        <!-- https://thenounproject.com/icon/about-6264304/ -->
        <symbol id="yn-icon-about">
          <path
            d="M 250,30.14007 C 129.07704,30.14007 30.14007,129.07704 30.14007,250 30.14007,370.92296 129.07704,469.85993 250,469.85993 370.92296,469.85993 469.85993,370.92296 469.85993,250 469.85993,129.07704 370.92296,30.14007 250,30.14007 Z m 0,395.74787 C 153.26163,425.88794 74.112056,346.73837 74.112056,250 74.112056,153.26163 153.26163,74.112056 250,74.112056 c 96.73837,0 175.88794,79.149574 175.88794,175.887944 0,96.73837 -79.14957,175.88794 -175.88794,175.88794 z"
          />
          <path
            d="m 265.3902,120.28264 c -26.3832,-4.3972 -52.76639,2.1986 -72.55378,19.78739 -19.7874,15.3902 -30.78039,39.57479 -30.78039,65.95798 h 43.97198 c 0,-13.19159 6.5958,-26.38319 15.3902,-32.97899 10.993,-8.79439 21.98599,-10.99299 37.37619,-8.79439 17.58879,2.1986 32.97899,17.58879 35.17759,35.17759 4.39719,21.98599 -6.5958,41.77338 -26.3832,48.36918 -24.18459,8.7944 -37.37618,30.78039 -37.37618,52.76638 v 15.3902 h 43.97198 v -15.3902 c 0,-6.59579 4.3972,-10.99299 10.993,-13.19159 39.57479,-15.3902 61.56078,-54.96498 54.96498,-96.73837 -8.7944,-35.17759 -39.57479,-65.95798 -74.75237,-70.35518 z"
          />
          <path d="m 228.01401,337.94397 h 43.97198 v 43.97198 h -43.97198 z" />
        </symbol>

        <symbol id="yn-icon-close">
          <path
            d="M 83.188924,42.278411 C 199.19023,158.27972 315.19155,274.28103 431.19286,390.28235 419.04742,402.42778 406.90199,414.57322 394.75655,426.71866 278.75524,310.71735 162.75392,194.71603 46.752613,78.714722 58.89805,66.569285 71.043486,54.423848 83.188924,42.278411 Z"
          />
          <path
            d="M 394.75657,42.278416 C 278.75525,158.27973 162.75394,274.28104 46.752628,390.28235 58.898063,402.42779 71.043499,414.57323 83.188935,426.71866 199.19025,310.71735 315.19156,194.71604 431.19287,78.714724 419.04744,66.569288 406.902,54.423852 394.75657,42.278416 Z"
          />
        </symbol>
      </defs>
    </svg>
    <yesterdays-news-player
      clips-resource-url=${clipsResourceUrl}
      initial-clip-url=${initialClipUrl}
      captions-resource-url=${captionsResourceUrl}
    ></yesterdays-news-player>
    <yesterdays-news-nav
      @fullscreenclick=${onFullscreenClick}
      about-page-url=${aboutPageUrl}
      .fullscreen=${fullscreen}
    ></yesterdays-news-nav>
  `
}

customElements.define(
  'yesterdays-news',
  component(YesterdaysNews, {
    observedAttributes: [
      'clips-resource-url',
      'initial-clip-url',
      'captions-resource-url',
      'about-page-url',
    ],
  })
)
