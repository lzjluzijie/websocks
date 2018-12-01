package core

import (
	"sync"
	"time"
)

//todo better stats
type Stats struct {
	Downloaded     uint64
	DownloadSpeed  uint64
	downloadMutex  sync.Mutex
	downloadSpeedA uint64

	Uploaded     uint64
	UploadSpeed  uint64
	uploadMutex  sync.Mutex
	uploadSpeedA uint64
}

func NewStats() (stats *Stats) {
	stats = &Stats{}
	go func() {
		for range time.Tick(time.Second) {
			stats.downloadMutex.Lock()
			stats.DownloadSpeed = stats.downloadSpeedA
			stats.downloadSpeedA = 0
			stats.downloadMutex.Unlock()
		}
	}()

	go func() {
		for range time.Tick(time.Second) {
			stats.uploadMutex.Lock()
			stats.UploadSpeed = stats.uploadSpeedA
			stats.uploadSpeedA = 0
			stats.uploadMutex.Unlock()
		}
	}()
	return
}

func (stats *Stats) AddDownloaded(downloaded uint64) {
	stats.downloadMutex.Lock()
	stats.Downloaded += downloaded
	stats.downloadSpeedA += downloaded
	stats.downloadMutex.Unlock()
	return
}

func (stats *Stats) AddUploaded(uploaded uint64) {
	stats.uploadMutex.Lock()
	stats.Uploaded += uploaded
	stats.uploadSpeedA += uploaded
	stats.uploadMutex.Unlock()
	return
}
