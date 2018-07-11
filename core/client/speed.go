package client

func (client *WebSocksClient) AddDownloaded(downloaded uint64) {
	client.downloadedC <- downloaded

	return
}

func (client *WebSocksClient) AddUploaded(uploaded uint64) {
	client.uploadedC <- uploaded
	return
}

//todo
func (client *WebSocksClient) Status() {

}
