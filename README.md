# Yesterday's News

This is really two projects in one

 - A Go port of [yesterdays-news-py](https://github.com/andyinabox/yesterdays-news-py/). The goal is to build assets and them push them to an object store.
 - A web project that will be like [yesterdays-news-of](https://github.com/andyinabox/yesterdays-news-of/) in a browser.

## Todo

 - [ ] Handle video loading/encoding errors in frontend
 - [ ] Set CORS on object store so videos can be preloaded using `fetch`  - I think I can do this by hitting the OpenStack API, but having some issues getting that to work (see `openstack` branch)
 - [ ] Containerize both projects
 - [ ] Deploy builder remotely
   - [ ] Set up cron process
 - [ ] Deploy server remotely