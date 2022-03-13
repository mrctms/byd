package downloader

import (
	"byd/tool"
	"path/filepath"
)

const outputTemplate = "%(title)s-%(id)s.%(ext)s"

func GetDownloadArgs(t string, url string, format string, outputPath string) []string {
	if t == "audio" {
		return downloadAudio(url, format, outputPath)
	} else if t == "video" {
		return downloadVideo(url, format, outputPath)
	}
	return []string{}
}

func downloadAudio(url string, format string, outputPath string) []string {
	var args []string
	args = append(args, url)
	args = append(args, getDefaultArgs(outputPath)...)
	args = append(args, "-x", "--audio-format", format)
	args = append(args, "--extract-audio")
	return args
}

func downloadVideo(url string, format string, outputPath string) []string {
	var args []string
	args = append(args, url)
	args = append(args, "--format", format)
	args = append(args, getDefaultArgs(outputPath)...)
	return args
}

func getDefaultArgs(outputPath string) []string {
	var args []string
	ffmpegPath, err := tool.GetFFmpegPath()
	if err == nil {
		args = append(args, "--ffmpeg-location", ffmpegPath)
	}
	args = append(args,
		"--add-metadata",
		"--geo-bypass",
		"--yes-playlist",
		"--restrict-filenames",
		"--no-abort-on-error",
		"--force-overwrites",
		"--output",
		filepath.Join(outputPath, outputTemplate))
	return args
}
