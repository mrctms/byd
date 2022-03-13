package workers

const DOWNLOAD_WORKER = "download-worker"

type Worker interface {
	Start() error
}
