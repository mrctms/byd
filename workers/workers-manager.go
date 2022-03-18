package workers

import (
	"byd/downloader"
	"byd/messaging"
	"byd/repo"
	"byd/tool"
)

type WorkersManager struct {
	workersMap map[string]func() (Worker, error)
	bus        messaging.MessageBus
}

func NewWorkersManager(bus messaging.MessageBus, repo repo.Repo) *WorkersManager {
	wm := new(WorkersManager)
	wm.bus = bus
	wm.workersMap = map[string]func() (Worker, error){}
	wm.workersMap[DOWNLOAD_WORKER] = func() (Worker, error) {
		d := downloader.NewYouTubeDownloader(tool.GetYtdPath())
		return NewDownloadWorker(wm.bus, repo, d), nil
	}
	return wm
}

func (w *WorkersManager) SpawnWorker(id string) error {
	worker, err := w.workersMap[id]()
	if err != nil {
		return err
	}
	worker.Start()
	return nil
}
