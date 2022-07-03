package item

import (
	"sync"

	pb "github.com/Psykepro/item-storage-protobuf/generated/item"
	domain "github.com/Psykepro/item-storage-server/_domain"
)

type RequestHandler struct {
	service domain.ItemService
	logger  domain.Logger
}

func NewRequestHandler(service domain.ItemService, logger domain.Logger) *RequestHandler {
	return &RequestHandler{
		logger:  logger,
		service: service,
	}
}

func (h *RequestHandler) Handle(requestChan chan *pb.ItemRequest) {
	for request := range requestChan {
		wg := new(sync.WaitGroup)
		wg.Add(1)
		switch request.Command {
		case pb.Command_CREATE:
			h.service.Create(request.Item, wg)
			break
		case pb.Command_DELETE:
			h.service.Delete(request.Item.Uuid, wg)
			break
		case pb.Command_GET:
			h.service.Get(request.Item.Uuid, wg)
			break
		case pb.Command_LIST:
			h.service.List(wg)
			break
		default:
			h.logger.Errorf("Unsupported command: [%s]", request.Command)
			wg.Done()
		}
		wg.Wait()
	}
}
