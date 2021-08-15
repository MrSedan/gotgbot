package dispatcher

import "log"

type DispServer struct {
	Hub     *Hub
	Id      int64
	Disps   map[int]*disp
	NewDisp chan *disp
	RemDisp chan *disp
}

func NewServer(id int64, h *Hub) *DispServer {
	return &DispServer{
		Hub:     h,
		Id:      id,
		Disps:   make(map[int]*disp),
		NewDisp: make(chan *disp, 100),
		RemDisp: make(chan *disp, 100),
	}
}

func (s *DispServer) Run() {
	defer func() {
		s.Hub.RemServer <- s
	}()
	for {
		select {
		case d := <-s.NewDisp:
			if _, ok := s.Disps[d.id]; !ok {
				s.Disps[d.id] = d
				log.Println("\nNew disp:", d.id)
			}
		case d := <-s.RemDisp:
			if _, ok := s.Disps[d.id]; ok {
				delete(s.Disps, d.id)
				log.Println("\nDisp ended:", d.id)
			}
			if len(s.Disps) == 0 {
				return
			}
		}
	}
}
