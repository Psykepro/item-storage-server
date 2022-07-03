package _domain

import (
	pb "github.com/Psykepro/item-storage-protobuf/generated/item"
)

type ItemRequestHandler interface {
	Handle(requestChannel chan *pb.ItemRequest)
}

type ItemService interface {
	Create(item *pb.Item)
	Delete(uuid string)
	Get(uuid string)
	List()
}
