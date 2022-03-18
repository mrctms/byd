package api

import (
	"byd/downloader"
	"byd/messaging"
	"byd/repo"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	url  string
	bus  messaging.MessageBus
	repo repo.Repo
}

func NewApiServer(url string, bus messaging.MessageBus, repo repo.Repo) *ApiServer {
	api := new(ApiServer)
	api.url = url
	api.bus = bus
	api.repo = repo
	return api
}

func (a *ApiServer) StartApi() error {
	r := mux.NewRouter()
	r.HandleFunc("/hubs/downloads", a.postDownloads).Methods(http.MethodPost)
	r.HandleFunc("/hubs/{hubId}/downloads", a.getHubDownloads).Methods(http.MethodGet)
	r.HandleFunc("/hubs/{hubId}", a.getHubs).Methods(http.MethodGet)
	r.HandleFunc("/hubs/{hubId}/downloads/{downloadId}", a.getDownloads).Methods(http.MethodGet)
	http.Handle("/", r)
	srv := &http.Server{
		Handler: r,
		Addr:    a.url,
	}
	return srv.ListenAndServe()
}

func (a *ApiServer) postDownloads(w http.ResponseWriter, r *http.Request) {
	var dArgs downloadArgs
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response(w, genericResponse{Message: "error while reading body", Details: err.Error()}, http.StatusInternalServerError)
		return
	}
	json.Unmarshal(content, &dArgs)
	// if _, ok := AudioTypes[dArgs.Format]; !ok {
	// 	response(w, GenericResponse{Message: "format not supported", Details: dArgs.Format}, http.StatusNotFound)
	// 	return
	// }
	var hub repo.Hub
	a.repo.Insert(&hub)
	for _, v := range dArgs.ArgsData {
		d := repo.Download{
			Type:   dArgs.Type,
			Format: dArgs.Format,
			URL:    v.URL,
			Dir:    v.Dir,
		}
		a.repo.Insert(&d)
		a.repo.Insert(repo.DownloadHub{HubID: hub.HubId, DownloadID: d.DownloadId})
	}
	msg, _ := json.Marshal(downloader.DownloadMsg{HubId: hub.HubId})
	a.bus.SendMessageToRabbitMQWorkQueue(string(msg))
	response(w, downloadQueuedReponse{Message: "Download queued", HubId: hub.HubId}, http.StatusOK)
}

func (a *ApiServer) getHubDownloads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hubId, ok := vars["hubId"]
	if !ok {
		response(w, nil, http.StatusBadRequest)
		return
	}
	var downloads []repo.Download
	hubIdConv, _ := strconv.Atoi(hubId)
	a.repo.GetDownloadsFromHub(&downloads, uint(hubIdConv))
	if len(downloads) != 0 {
		resp := hubDownloadsResponse{HubId: uint(hubIdConv)}
		for _, v := range downloads {
			resp.Downloads = append(resp.Downloads, hubDownload{DownloadId: v.DownloadId,
				ETA: v.ETA, Details: v.Details, Dir: v.Dir,
				Format: v.Format, Progress: v.Progress, Type: v.Type,
				Speed: v.Speed, Total: v.Total, URL: v.URL})
		}
		response(w, resp, http.StatusOK)
		return
	}
	response(w, nil, http.StatusNotFound)
}

func (a *ApiServer) getHubs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hubId, ok := vars["hubId"]
	if !ok {
		response(w, nil, http.StatusBadRequest)
		return
	}
	var hub repo.Hub
	err := a.repo.Get(&hub, map[string]interface{}{"hub_id": hubId})
	if (hub == repo.Hub{}) {
		response(w, genericResponse{Message: fmt.Sprintf("hub with id %s not found", hubId), Details: err.Error()}, http.StatusNotFound)
		return
	}
	if err != nil {
		response(w, genericResponse{Details: err.Error()}, http.StatusInternalServerError)
		return
	}
	response(w, hubReponse{HubId: hub.HubId, Result: hub.Result, ZipName: hub.ZipName, Details: hub.Details}, http.StatusOK)
}

func (a *ApiServer) getDownloads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hubId, ok := vars["hubId"]
	if !ok {
		response(w, nil, http.StatusBadRequest)
		return
	}
	downloadId, ok := vars["downloadId"]
	if !ok {
		response(w, nil, http.StatusBadRequest)
		return
	}
	hubIdConv, _ := strconv.Atoi(hubId)
	downloadIdConv, _ := strconv.Atoi(downloadId)
	var download repo.Download
	a.repo.GetDownloadFromHub(&download, uint(hubIdConv), uint(downloadIdConv))
	if (download != repo.Download{}) {
		response(w, hubDownload{DownloadId: download.DownloadId,
			ETA: download.ETA, Details: download.Details, Dir: download.Dir,
			Format: download.Format, Progress: download.Progress, Type: download.Type,
			Speed: download.Speed, Total: download.Total, URL: download.URL}, http.StatusOK)
	} else {
		response(w, nil, http.StatusNotFound)
	}
}

func response(w http.ResponseWriter, body interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp, _ := json.Marshal(body)
	w.Write(resp)
}
