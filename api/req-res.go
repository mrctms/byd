package api

type downloadArgs struct {
	ArgsData []argsData `json:"data"`
	Format   string     `json:"format"`
	Type     string     `json:"type"`
}

type argsData struct {
	URL string `json:"url"`
	Dir string `json:"dir"`
}

type genericResponse struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

type downloadQueuedReponse struct {
	Message string `json:"message"`
	HubId   uint   `json:"hubId"`
}

type hubReponse struct {
	HubId   uint   `json:"hubId"`
	Result  string `json:"result"`
	Details string `json:"details"`
	ZipName string `json:"zip_name"`
}

type hubDownloadsResponse struct {
	HubId     uint          `json:"hubId"`
	Downloads []hubDownload `json:"downloads"`
}

type hubDownload struct {
	DownloadId uint   `json:"downloadId"`
	ETA        string `json:"eta"`
	Speed      string `json:"speed"`
	Total      string `json:"total"`
	Progress   string `json:"progress"`
	URL        string `json:"url"`
	Type       string `json:"type"`
	Format     string `json:"format"`
	Dir        string `json:"dir"`
	Details    string `json:"details"`
}
