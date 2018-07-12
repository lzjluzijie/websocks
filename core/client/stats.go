package client

func (client *WebSocksClient) AddDownloaded(downloaded uint64) {
	client.downloadMutex.Lock()
	client.Downloaded += downloaded
	client.downloadSpeedA += downloaded
	client.downloadMutex.Unlock()
	return
}

func (client *WebSocksClient) AddUploaded(uploaded uint64) {
	client.uploadMutex.Lock()
	client.Uploaded += uploaded
	client.uploadSpeedA += uploaded
	client.uploadMutex.Unlock()
	return
}

type Stats struct {
	Downloaded    uint64
	DownloadSpeed uint64
	Uploaded      uint64
	UploadSpeed   uint64

	//todo conns
}

func (client *WebSocksClient) Status() (stats *Stats) {
	return &Stats{
		Downloaded:    client.Downloaded,
		DownloadSpeed: client.DownloadSpeed,
		Uploaded:      client.Uploaded,
		UploadSpeed:   client.UploadSpeed,
	}
}
