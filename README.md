# NOTE: Sorry but this doesn't work anymore!

# Running instructions:

- Make a `segments/` dir in the root dir of the project
- Make `segments/full/` and `segments/partial/` dir
- Take the `sample-playlist-master.m3u8` and place it in `segments/` dir as `playlist_master.m3u8`(You can make changes to this file according to the needs, this creation and updation will be automated in the future)
- Start a static file server from the segments/ dir, now replace the URL in `rtmp-to-hls.sh`, `main.go` and `playlist_master.m3u8` file
- Start the ll hls server provided by apple in `segments/` dir
- Go to the playback URL and enjoy the video!
