package downloader

import "regexp"

const downloadRegex = `\[download\]\s+(?:(?P<percent>[\d\.]+)%(?:\s+of\s+\~?(?P<total>[\d\.\w]+))?\s+at\s+(?:(?P<speed>[\d\.\w]+\/s)|[\w\s]+)\s+ETA\s(?P<eta>[\d\:]+))?`

type DownloaderOutput struct {
	ETA      string
	Speed    string
	Total    string
	Progress string
	Error    string
}

type YouTubeDownloader struct {
	ydPath string
}

func NewYouTubeDownloader(ydPath string) *YouTubeDownloader {
	yd := new(YouTubeDownloader)
	yd.ydPath = ydPath
	return yd
}

func (y *YouTubeDownloader) Download(progress chan<- DownloaderOutput, args ...string) {
	defer close(progress)
	err := runDownloadProcess(y.ydPath, func(output string) {
		reg := regexp.MustCompile(downloadRegex)
		if reg.MatchString(output) {
			match := reg.FindStringSubmatch(output)
			if len(match) >= 4 {
				out := DownloaderOutput{
					Progress: match[1],
					Total:    match[2],
					Speed:    match[3],
					ETA:      match[4],
				}
				if (out != DownloaderOutput{}) {
					progress <- out
				}
			}
		}
	}, func(output string) {
		progress <- DownloaderOutput{Error: output}
	}, args...)
	if err != nil {
		progress <- DownloaderOutput{Error: err.Error()}
	}
}
