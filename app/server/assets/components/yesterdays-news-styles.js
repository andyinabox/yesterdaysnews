import { css } from 'lit'
// import { virtual } from 'haunted'

export default css`
:host {
  --yn-nav-btn-size: 10vmin;
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

yesterdays-news-player {
  position: relative;
  width: var(--yn-video-width);
  height: var(--yn-video-height);
  background-color: #111;
}
yesterdays-news-video {
  display: block;
  z-index: 0;
}
yesterdays-news-video > video {
  width: var(--yn-video-width);
  height: var(--yn-video-height);
}
yesterdays-news-video .play-button {
  z-index: 5;
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  background: transparent;
  border: none;
  padding: 0;
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

// export const YesterdaysNewsStyles = virtual(() => {
//   return html`
//     <style>
//       :host {
//         --yn-nav-btn-size: 10vmin;
//         --yn-nav-btn-margin: 2vmin;

//         /* these will be set with javascript */
//         --yn-video-width: 0px;
//         --yn-video-height: 0px
//         --yn-caption-font-size: var(--yn-nav-font-size);

//         display: flex;
//         align-items: center;
//         justify-content: center;
//       }

//       button {
//         cursor: pointer;
//       }

//       svg > use {
//         fill: white;
//         stroke: none;
//         display: block;
//         width: 100%;
//         height: 100%;
//       }

//       yesterdays-news-player {
//         position: relative;
//         width: var(--yn-video-width);
//         height: var(--yn-video-height);
//         background-color: #111;
//       }
//       yesterdays-news-video {
//         display: block;
//         z-index: 0;
//       }
//       yesterdays-news-video > video {
//         width: var(--yn-video-width);
//         height: var(--yn-video-height);
//       }
//       yesterdays-news-video .play-button {
//         z-index: 5;
//         position: absolute;
//         left: 50%;
//         top: 50%;
//         transform: translate(-50%, -50%);
//         background: transparent;
//         border: none;
//         padding: 0;
//         width: var(--yn-nav-btn-size);
//         height: var(--yn-nav-btn-size);
//       }
//       yesterdays-news-captions {
//         position: absolute;
//         top: 0px;
//         left: 0px;
//         z-index: 1;
//         width: var(--yn-video-width);
//         height: var(--yn-video-height);
//         font-size: var(--yn-caption-font-size);
//         display: flex;
//         align-items: flex-end;
//         justify-content: center;
//       }
//       yesterdays-news-captions > .caption {
//         font-family: monospace;
//         text-align: center;
//         margin-bottom: 0.5em;
//       }
//       yesterdays-news-captions > .caption > .word {
//         background-color: black;
//         color: white;
//       }

//       yesterdays-news-nav {
//         z-index: 3;
//         position: fixed;
//         bottom: 0px;
//         left: 0px;
//         width: 100%;
//         display: flex;
//         align-items: center;
//         justify-content: space-between;
//       }
//       yesterdays-news-nav > .btn {
//         position: absolute;
//         bottom: 0px;
//         cursor: pointer;
//         display: block;
//         border: none;
//         padding: 0;
//         background-color: transparent;
//         width: var(--yn-nav-btn-size);
//         height: var(--yn-nav-btn-size);
//         margin: var(--yn-nav-btn-margin);
//       }
//       yesterdays-news-nav > .btn.btn-about {
//         left: 0px;
//       }
//       yesterdays-news-nav > .btn.btn-fullscreen {
//         right: 0px;
//       }
//     </style>
//   `
// })
