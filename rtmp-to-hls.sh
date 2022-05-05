#!/bin/bash -e

echo "$gopSize"
echo "$segmentSize"

gopSize=$1
segmentSize=$2
fps=$3

if [[ $gopSize == "" ]]
then
    echo "GOP Size not provided"
    exit 1
fi

if [[ $segmentSize == "" ]]
then
    echo "Segment size not provided"
    exit 1
fi

if [[ $fps == "" ]]
then
    echo "FPS not provided"
    exit 1
fi

# -maxrate 400k -bufsize 1835k -pix_fmt yuv420p\
# MPEG-2-compressed video dat
# -c:v libx264 -c:a aac -ac 1 -crf 18 \
# segments/playlist_master.m3u8

# ffmpeg -listen 1 -v verbose -i rtmp://127.0.0.1:1937/live/movie\
#     -preset veryfast -sc_threshold 0 \
#     -map [vstream001] -c:v:0 libx264 -g:v:0 gopSize \
#     # -map [vstream002] -c:v:1 libx264 -g:v:1 fps  \
#     -map a:0 -map a:0 -c:a aac -b:a 128k -ac 2 \
#     -f hls -flags -global_header -hls_delete_threshold 10 -hls_time $segmentSize \
#     -master_pl_name master.m3u8 \
#     # -hls_segment_filename segments/data%06d.ts \
#     -var_stream_map "v:0,a:0 v:1,a:1" stream_%v.m3u8

ffmpeg -listen 1 -v verbose -i rtmp://127.0.0.1:1937/live/movie \
    -c:v libx264 -c:a aac -ac 1 -strict -2 -crf 18 \
    -preset veryfast -g $gopSize -sc_threshold 0 \
    -profile:v baseline -maxrate 400k -bufsize 1835k -pix_fmt yuv420p \
    -flags -global_header -hls_time $segmentSize -hls_delete_threshold 10\
    -hls_segment_filename segments/partial/data%06d.ts \
    -hls_flags 'independent_segments;program_date_time' \
    -hls_base_url "https://app-fenil-files.loca.lt/partial/" \
    segments/partial/playlist_partial_master.m3u8 \
    -f flv rtmp://127.0.0.1:1938/live/movie

    # -f hls -c:v libx264 -c:a aac -ac 1 -strict -2 -crf 18 \
    # -preset veryfast -g $segmentSize -sc_threshold 0 \
    # -profile:v baseline -maxrate 400k -bufsize 1835k -pix_fmt yuv420p \
    # -flags -global_header -hls_time 2 -hls_delete_threshold 10\
    # -hls_segment_filename segments/full/data%06d.ts -start_number 1 \
    # segments/full/playlist_full_master.m3u8 \

