package websocket

import (
	"log"
	"net/http"

	"gitlab.com/telegram_clone/websocket_service/config"
	grpcPkg "gitlab.com/telegram_clone/websocket_service/pkg/grpc_client"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "websocket/home.html")
}

func Run(cfg config.Config, grpcClient grpcPkg.GrpcClientI) {
	hub := newHub(grpcClient)
	go hub.run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r, grpcClient)
	})

	log.Println("Websocket server started in port ", cfg.WsPort)
	err := http.ListenAndServe(cfg.WsPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
