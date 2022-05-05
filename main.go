package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/sync/errgroup"
)

func main() {
	partialSegmentSize := 0.2

	fps, err := getFPS()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Fetched FPS!", fps)

	// startRTMPServers(errs)
	// fmt.Println("RTMP server started!")

	// time.Sleep(7 * time.Second)

	errs, _ := errgroup.WithContext(context.Background())

	startHLSServer(fps, partialSegmentSize, errs)
	fmt.Println("Started Server!")

	err = errs.Wait()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Completed process")
}

func startHLSServer(fps int, partialSegmentSize float64, errs *errgroup.Group) {

	errs.Go(func() error {
		out, err, stderr := execCommand("ffmpeg", "-listen", "1", "-v", "verbose", "-i", "rtmp://127.0.0.1:1938/live/movie", "-c:v", "libx264", "-c:a", "aac", "-ac", "1", "-strict", "-2", "-crf", "18", "-preset", "veryfast", "-g", fmt.Sprint(fps), "-sc_threshold", "0", "-profile:v", "baseline", "-maxrate", "400k", "-bufsize", "1835k", "-pix_fmt", "yuv420p", "-flags", "-global_header", "-hls_time", "2", "-hls_delete_threshold", "10", "-hls_segment_filename", "segments/full/data%06d.ts", "-hls_flags", "independent_segments;program_date_time", "-hls_base_url", "https://app-fenil-files.loca.lt/full/", "segments/full/playlist_full_master.m3u8")
		if err != nil {
			return errors.New(fmt.Sprint(err) + ": " + stderr)
		}

		fmt.Println("out: ", out)

		return nil
	})

	gopSize := float64(fps) * partialSegmentSize

	errs.Go(func() error {
		out, err, stderr := execCommand("./rtmp-to-hls.sh", fmt.Sprint(gopSize), fmt.Sprint(partialSegmentSize), fmt.Sprint(fps))
		if err != nil {
			return errors.New(fmt.Sprint(err) + ": " + stderr)
		}

		fmt.Println("out: ", out)

		return nil
	})
}

func getFPS() (int, error) {
	out, err, stderr := execCommand("ffprobe", "-listen", "1", "-v", "error", "-select_streams", "v", "-of", "default=noprint_wrappers=1:nokey=1", "-show_entries", "stream=avg_frame_rate", "rtmp://127.0.0.1:1937/live/movie")
	if err != nil {
		return 0, errors.New(fmt.Sprint(err) + ": " + stderr)
	}

	ratio := strings.Split(out, "/")

	numerator, err := strconv.Atoi(ratio[0])
	if err != nil {
		return 0, errors.New("Error converting numerator to number")
	}

	denominator, err := strconv.Atoi(strings.TrimSuffix(ratio[1], "\n"))
	if err != nil {
		return 0, errors.New("Error converting denominator to number")
	}

	fps := numerator / denominator

	return fps, nil
}

func execCommand(mainCommand string, flags ...string) (string, error, string) {
	cmd := exec.Command(mainCommand, flags...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", err, stderr.String()
	}

	return out.String(), nil, ""
}

// func startRTMPServers(errs *errgroup.Group) {
// 	// errs.Go(func() error {
// 	// 	_, err, stderr := execCommand("ffmpeg", "-listen", "1", "-i", "rtmp://127.0.0.1:1937/live/movie", "-f", "null", "/dev/null")
// 	// 	if err != nil {
// 	// 		return errors.New(fmt.Sprint(err) + ": " + stderr)
// 	// 	}

// 	// 	return nil
// 	// })

// 	// errs.Go(func() error {
// 	// 	_, err, stderr := execCommand("ffmpeg", "-listen", "1", "-i", "rtmp://127.0.0.1:1936/live/movie", "-f", "flv", "rtmp://127.0.0.1:1937/live/movie")
// 	// 	if err != nil {
// 	// 		return errors.New(fmt.Sprint(err) + ": " + stderr)
// 	// 	}

// 	// 	return nil
// 	// })
// }
// out, err, stderr := execCommand("ffmpeg", "-listen", "1", "-v", "verbose", "-i", "rtmp://127.0.0.1:1937/live/movie", "-c:v", "libx264", "-c:a", "aac", "-ac", "1", "-strict", "-2", "-crf", "18", "-preset", "veryfast", "-g", fmt.Sprint(gopSize), "-sc_threshold", "0", "-profile:v", "baseline", "-maxrate", "400k", "-bufsize", "1835k", "-pix_fmt", "yuv420p", "-flags", "-global_header", "-hls_time", fmt.Sprint(segmentSize), "-hls_delete_threshold", "10", "-hls_segment_filename", "segments/data%06d.ts", "-start_number", "1", "segments/playlist_master.m3u8")
