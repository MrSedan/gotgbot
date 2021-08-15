package dispatcher

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type disp struct {
	server    *DispServer
	handler   *ConvHandler
	curStage  int
	id        int
	NextStage chan bool
}

func CreateDisp(id int, ds *DispServer, h *ConvHandler) *disp {
	return &disp{
		server:    ds,
		handler:   h,
		curStage:  0,
		id:        id,
		NextStage: make(chan bool, 10),
	}
}

func (d *disp) Run(bot *tgbotapi.BotAPI) {
	defer func() {
		d.server.RemDisp <- d
	}()
	for {
		if <-d.NextStage {
			d.curStage += 1
			msg := tgbotapi.NewMessage(d.server.Id, d.handler.Stages[d.curStage-1])
			bot.Send(msg)
			if d.curStage == d.handler.StagesCount {
				return
			}
		}
	}
}
