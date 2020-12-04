package server

type WorkerService struct {
	messageHandler MessageHandler
}

type MessageHandler interface {
	HandleMessage(message []byte)
}

type defaultMessageHandler struct {
}

func NewWorkerService() *WorkerService {
	return &WorkerService{
		messageHandler: &defaultMessageHandler{},
	}
}

func (w *WorkerService) GetMessageHandler() MessageHandler {
	return w.messageHandler
}

func (w *WorkerService) SetMessageHandler(h MessageHandler) {
	w.messageHandler = h
}

func (h *defaultMessageHandler) HandleMessage(message []byte) {

}
