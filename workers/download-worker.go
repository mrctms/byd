package workers

import (
	"byd/downloader"
	"byd/messaging"
	"byd/repo"
	"byd/tool"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/mholt/archiver"
)

type DownloadWorker struct {
	downloader   *downloader.YouTubeDownloader
	busWorkQueue messaging.MessageBusConsumer
	repo         repo.Repo
}

func NewDownloadWorker(bus messaging.MessageBus, repo repo.Repo, downloader *downloader.YouTubeDownloader) *DownloadWorker {
	d := new(DownloadWorker)
	d.downloader = downloader
	d.busWorkQueue = bus.CreateRabbitMQWorker()
	d.repo = repo
	return d
}

func (d *DownloadWorker) Start() error {
	d.busWorkQueue.OnMessageReceived(d.execute)
	err := d.busWorkQueue.StartListening()
	if err != nil {
		return err
	}
	return nil
}

func (d *DownloadWorker) execute(msg string) error {
	var downloadMsg downloader.DownloadMsg
	err := json.Unmarshal([]byte(msg), &downloadMsg)
	if err != nil {
		return err
	}
	hub := repo.Hub{HubId: downloadMsg.HubId}
	var downloads []repo.Download
	d.repo.GetDownloadsFromHub(&downloads, downloadMsg.HubId)
	outputPath, err := tool.GetOutputPath()
	if err != nil {
		d.repo.Update(&hub, map[string]interface{}{"result": "fatal-error", "details": err.Error()})
		return err
	}
	mainOutputPath := filepath.Join(os.TempDir(), uuid.NewString())
	os.Mkdir(mainOutputPath, 0777)
	defer os.RemoveAll(mainOutputPath)
	var anyError bool
	for _, v := range downloads {
		c := make(chan downloader.DownloaderOutput)
		args := downloader.GetDownloadArgs(v.Type, v.URL, v.Format, filepath.Join(mainOutputPath, v.Dir))
		go func() {
			d.downloader.Download(c, args...)
		}()
		for o := range c {
			o.Error = strings.Replace(o.Error, "\x00", "", -1) // do better
			if !anyError {
				anyError = o.Error != ""
			}
			d.repo.Update(&repo.Download{DownloadId: v.DownloadId}, map[string]interface{}{
				"eta":      o.ETA,
				"speed":    o.Speed,
				"progress": o.Progress,
				"total":    o.Total,
				"details":  o.Error,
			})
		}
	}
	zipName := uuid.NewString()
	err = archiver.Archive([]string{mainOutputPath}, filepath.Join(outputPath, fmt.Sprintf("%s.zip", zipName)))
	if err != nil {
		d.repo.Update(&hub, map[string]interface{}{"status": "fatal-error", "details": err.Error()})
		return err
	}
	result := make(map[string]interface{})
	if anyError {
		result["result"] = "errors"
	} else {
		result["result"] = "ok"
	}
	result["zip_name"] = zipName
	d.repo.Update(&hub, result)
	return nil
}
