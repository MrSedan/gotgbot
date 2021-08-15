package dispatcher

import "log"

type Hub struct {
	Servers   map[int64]*DispServer
	NewServer chan *DispServer
	RemServer chan *DispServer
}

func (h *Hub) Run() {
	for {
		select {
		case s := <-h.NewServer:
			if _, ok := h.Servers[s.Id]; !ok {
				h.Servers[s.Id] = s
				log.Println("\nNew server:", s.Id)
			}
		case s := <-h.RemServer:
			if _, ok := h.Servers[s.Id]; ok {
				delete(h.Servers, s.Id)
				log.Println("\nServer ended:", s.Id)
			}
		}
	}
}
