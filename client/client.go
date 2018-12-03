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
	if client.Mux {
		client.muxGroup = mux.NewGroup(true)
		log.Println("group created")
		go func() {
			//todo
			for {
				if len(client.muxGroup.MuxWSs) == 0 {
					err := client.OpenMux()
					if err != nil {
						log.Printf(err.Error())
						continue
					}
				}
				//这个弱智BUG折腾了我一天
				time.Sleep(time.Second)
			}
		}()
	}

	listener, err := net.ListenTCP("tcp", client.ListenAddr)
	if err != nil {
		return err
	}

	log.Printf("Start to listen at %s", client.ListenAddr.String())

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
	//debug log
	//log.Println("new socks5 conn")

	lc, err := NewLocalConn(conn)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	host := lc.Host

	if client.Mux {
		muxConn, err := client.muxGroup.NewMuxConn(host)
		if err != nil {
			log.Printf(err.Error())
			return
		}

		//debug log
		log.Printf("created new mux conn: %x", muxConn.ID)

		muxConn.Run(conn)
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
	return
}
