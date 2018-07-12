package client

func (client *WebSocksClient) AddDownloaded(downloaded uint64) {
	client.downloadMutex.Lock()
	client.downloadSpeedA += downloaded
	client.downloadMutex.Unlock()
	return
}

func (client *WebSocksClient) AddUploaded(uploaded uint64) {
	client.uploadMutex.Lock()
	client.uploadSpeedA += uploaded
	client.uploadMutex.Unlock()
	return
}

//todo
func (client *WebSocksClient) Status() {

}
