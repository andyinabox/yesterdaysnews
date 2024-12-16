# Yesterday's News Downloader

A Go port of [yesterdays-news-py](https://github.com/andyinabox/yesterdays-news-py/).

## Todo

- [x] Fetch video list from Google for channel
- [x] Download videos from video list
  - [x] Download videos
  - [ ] FIX: Ensure only portrait videos are downloaded
- [x] Extract subtitles from videos and convert to text file
- [ ] Cut videos into smaller clips
- [ ] Generate markov model from subtitles and journal text
- [ ] Combine 
- [ ] Export all relevant artifacts to S3
  - [ ] Models
  - [ ] Clips
  - [ ] Output video
  - [ ] Output subtitles file (`.srt`)
  - [ ] Manifest file