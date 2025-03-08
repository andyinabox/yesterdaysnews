import { css } from 'lit'

export default css`
:host {
  --yn-nav-btn-size: 10vmin;
  --yn-loading-icon-size: 20vmin;
  --yn-nav-btn-margin: 2vmin;
  /* these will be set with javascript */
  --yn-video-width: 0px;
  --yn-video-height: 0px
  --yn-caption-font-size: var(--yn-nav-font-size);

  display: flex;
  align-items: center;
  justify-content: center;
}

button {
  cursor: pointer;
}

svg > use {
  fill: white;
  stroke: none;
  display: block;
  width: 100%;
  height: 100%;
}

@keyframes changeOpacity1 {
  0% {
    opacity: 0%
  }
  25% {
    opacity: 100%;
  }
  50% {
    opacity: 0%;
  }
}

@keyframes changeOpacity2 {
  0% {
    opacity: 0%;
  }
  25% {
    opacity: 0%
  }
  50% {
    opacity: 100%;
  }
  75% {
    opacity: 0%;
  }
}

@keyframes changeOpacity3 {
  0% {
    opacity: 0%;
  }
  50% {
    opacity: 0%;
  }
  75% {
    opacity: 100%
  }
  100% {
    opacity: 0%;
  }
}


.loading-icon {
  fill: #fff;
  width: var(--yn-loading-icon-size);
  height: var(--yn-loading-icon-size);
}

.loading-icon .wave1,
.loading-icon .wave2,
.loading-icon .wave3 {
  animation-duration: 2s;
  animation-iteration-count: infinite;
  animation-timing-function: step-end;
}

.loading-icon .wave1 {
  animation-name: changeOpacity3;
}
.loading-icon .wave2 {
  animation-name: changeOpacity2;
}
.loading-icon .wave3 {
  animation-name: changeOpacity1;
}



yesterdays-news-player {
  position: relative;
  width: var(--yn-video-width);
  height: var(--yn-video-height);
  background-color: #111;
}

yesterdays-news-static {
  display: block;
}
yesterdays-news-static > canvas {
  width: var(--yn-video-width);
  height: var(--yn-video-height);
}

yesterdays-news-video {
  display: block;
  position: absolute;
  z-index: 0;
  top: 0px;
  left: 0px;
}
yesterdays-news-video > video {
  width: var(--yn-video-width);
  height: var(--yn-video-height);
}

yesterdays-news-video .centered-icon {
  z-index: 5;
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  background: transparent;
  border: none;
  padding: 0;
  // width: var(--yn-nav-btn-size);
  // height: var(--yn-nav-btn-size);
}


yesterdays-news-video .play-button {
  width: var(--yn-nav-btn-size);
  height: var(--yn-nav-btn-size);
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
yesterdays-news-nav > .btn {
  position: absolute;
  bottom: 0px;
  cursor: pointer;
  display: block;
  border: none;
  padding: 0;
  background-color: transparent;
  width: var(--yn-nav-btn-size);
  height: var(--yn-nav-btn-size);
  margin: var(--yn-nav-btn-margin);
}
yesterdays-news-nav > .btn.btn-about {
  left: 0px;
}
yesterdays-news-nav > .btn.btn-fullscreen {
  right: 0px;
}
`
