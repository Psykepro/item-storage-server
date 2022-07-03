package _domain

import (
	"sync"

	pb "github.com/Psykepro/item-storage-protobuf/generated/item"
)

type ItemRequestHandler interface {
	Handle(requestChannel chan *pb.ItemRequest)
}

type ItemService interface {
	Create(item *pb.Item, wg *sync.WaitGroup)
	Delete(uuid string, wg *sync.WaitGroup)
	Get(uuid string, wg *sync.WaitGroup)
	List(wg *sync.WaitGroup)
}
