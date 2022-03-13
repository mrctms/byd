package repo

type Repo interface {
	Insert(model interface{}) error
	Update(model interface{}, data map[string]interface{}) error
	GetFirst(model interface{}, id uint) error
	Get(model interface{}, data map[string]interface{}) error
	GetDownloadFromHub(model *Download, hubId uint, downloadId uint) error
	GetDownloadsFromHub(model *[]Download, hubId uint) error
}
