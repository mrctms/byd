package repo

type Download struct {
	DownloadId uint   `gorm:"primary_key;column:download_id"`
	ETA        string `gorm:"column:eta"`
	Speed      string `gorm:"column:speed"`
	Total      string `gorm:"column:total"`
	Progress   string `gorm:"column:progress"`
	URL        string `gorm:"column:url"`
	Type       string `gorm:"column:type"`
	Format     string `gorm:"column:format"`
	Dir        string `gorm:"column:dir"`
	Details    string `gorm:"column:details"`
}

type Hub struct {
	HubId   uint   `gorm:"primary_key;column:hub_id"`
	Result  string `gorm:"column:result"`
	Details string `gorm:"column:details"`
	ZipName string `gorm:"column:zip_name"`
}

type DownloadHub struct {
	HubID      uint `gorm:"column:hub_id;index"`
	Hub        Hub
	DownloadID uint `gorm:"column:download_id;index"`
	Download   Download
}
