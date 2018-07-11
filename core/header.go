package core

func NewHostHeader(host string) (header map[string][]string) {
	return map[string][]string{
		"WebSocks-Host": {host},
	}
}
