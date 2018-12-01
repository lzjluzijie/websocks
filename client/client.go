package client

import (
	"log"
	"net"
	"time"

	"github.com/lzjluzijie/websocks/core/mux"

	"net/url"

	"github.com/gorilla/websocket"
	"github.com/lzjluzijie/websocks/core"
)

type WebSocksClient struct {
	ServerURL  *url.URL
	ListenAddr *net.TCPAddr

	dialer *websocket.Dialer

	Mux bool

	muxGroup *mux.Group

	//todo
	//control
	stopC chan int

	CreatedAt time.Time
	Stats     *core.Stats
}

func (client *WebSocksClient) Run() (err error) {
	listener, err := net.ListenTCP("tcp", client.ListenAddr)
	if err != nil {
		return err
	}

	log.Printf("Start to listen at %s", client.ListenAddr.String())

	if client.Mux {
		group := mux.NewGroup(true)
		go func() {
			//todo
			for {
				if len(group.MuxWSs) == 0 {
					err := client.OpenMux()
					if err != nil {
						log.Printf(err.Error())
						continue
					}
				}
			}
		}()
	}

	go func() {
		client.stopC = make(chan int)
		<-client.stopC
		err = listener.Close()
		if err != nil {
			log.Printf(err.Error())
			return
		}

		log.Print("stopped")
	}()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf(err.Error())
			break
		}

		go client.HandleConn(conn)
	}
	return nil
}

func (client *WebSocksClient) Stop() {
	client.stopC <- 1911
	return
}

func (client *WebSocksClient) HandleConn(conn *net.TCPConn) {
	lc, err := NewLocalConn(conn)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	host := lc.Host

	if client.Mux {
		err = client.muxGroup.NewMuxConn(host)
		if err != nil {
			log.Printf(err.Error())
			return
		}
		return
	}

	ws, err := client.DialWebSocket(core.NewHostHeader(host))
	if err != nil {
		log.Printf(err.Error())
		return
	}

	lc.Run(ws)
	return
}

func (client *WebSocksClient) DialWebSocket(header map[string][]string) (ws *core.WebSocket, err error) {
	wsConn, _, err := client.dialer.Dial(client.ServerURL.String(), header)
	if err != nil {
		return
	}

	ws = core.NewWebSocket(wsConn, client.Stats)
	//client.connMutex.Lock()
	//client.wsConns = append(client.wsConns, ws)
	//client.connMutex.Unlock()
	return
}
