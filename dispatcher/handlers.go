package dispatcher

type ConvHandler struct {
	//*Придумать, как привязать к тг
	StagesCount int
	Stages      []string
}

func HelloHandler() *ConvHandler {
	ch := &ConvHandler{
		StagesCount: 4,
	}
	ch.Stages = append(ch.Stages, "What's your name?", "Ok, what's your surname?", "Your age?", "Amogus")
	return ch
}
