package item

import (
	"fmt"
	"sync"

	pb "github.com/Psykepro/item-storage-protobuf/generated/item"
	domain "github.com/Psykepro/item-storage-server/_domain"
	"github.com/Psykepro/item-storage-server/pkg/collections"
	"google.golang.org/protobuf/proto"
)

type Service struct {
	storage      *collections.OrderedDict[*pb.Item]
	stdOutLogger domain.Logger
	fileLogger   domain.Logger
}

func NewService(stdOutLogger domain.Logger, fileLogger domain.Logger) *Service {
	return &Service{
		storage:      collections.NewOrderedDict[*pb.Item](),
		stdOutLogger: stdOutLogger,
		fileLogger:   fileLogger,
	}
}

func (s *Service) Create(item *pb.Item, wg *sync.WaitGroup) {
	defer wg.Done()
	s.stdOutLogger.Debug("Creating new item ...")

	if item.Uuid == "" {
		errMsg := fmt.Sprintf("Failed to create Item. Err: [invalid item uuid - %s]", item.Uuid)
		s.logCreateErrorResponse(errMsg)
		return
	}

	if _, ok := s.storage.Get(item.Uuid); ok {
		errMsg := fmt.Sprintf("Failed to create Item. Err: [item with uuid - [%s] already exists]", item.Uuid)
		s.logCreateErrorResponse(errMsg)
		return
	}

	s.storage.Set(item.Uuid, item)
	response := &pb.CreateItemResponse{
		Item:  item,
		Error: nil,
	}
	asBytes, _ := proto.Marshal(response)
	s.fileLogger.Infof("%#v", asBytes)
	s.stdOutLogger.Debugf("Successfully created item with uuid: [%s]", item.Uuid)
}

func (s *Service) Delete(uuid string, wg *sync.WaitGroup) {
	defer wg.Done()
	s.stdOutLogger.Debugf("Initiating [DELETE] item with uuid: [%s]", uuid)

	if uuid == "" {
		errMsg := fmt.Sprintf("Failed to [DELETE] item. Err: [invalid item uuid - %s]", uuid)
		s.logDeleteErrorResponse(uuid, errMsg)
		return
	}

	ok := s.storage.Remove(uuid)
	if !ok {
		errMsg := fmt.Sprintf("Failed to [DELETE] item with uuid: [%s]. Err: [item not exist]", uuid)
		s.logDeleteErrorResponse(uuid, errMsg)
		return
	}

	response := &pb.DeleteItemResponse{
		ItemUuid: uuid,
		Error:    nil,
	}
	asBytes, _ := proto.Marshal(response)
	s.fileLogger.Infof("%#v", asBytes)
	s.stdOutLogger.Debugf("Successful [DELETE] item with uuid: [%s]", uuid)
}

func (s *Service) Get(uuid string, wg *sync.WaitGroup) {
	defer wg.Done()
	s.stdOutLogger.Debugf("Initiating [GET] item with uuid: [%s] ...", uuid)
	if uuid == "" {
		s.stdOutLogger.Errorf("Failed to [GET] Item. Err: [invalid item uuid - %s]", uuid)
	}

	item, ok := s.storage.Get(uuid)
	if !ok {
		errMsg := fmt.Sprintf("Failed to [GET] item. Err: [item not exist]")
		s.logGetErrorResponse(errMsg)
		return
	}

	response := &pb.GetItemResponse{
		Item:  item,
		Error: nil,
	}
	asBytes, err := proto.Marshal(response)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to marshal item on [GET]. Err: %s", err)
		s.logGetErrorResponse(errMsg)
		return
	}

	s.fileLogger.Infof("%#v", asBytes)
	s.stdOutLogger.Debugf("Successful [GET] item with uuid: [%s].", uuid)
}

func (s *Service) List(wg *sync.WaitGroup) {
	defer wg.Done()
	s.stdOutLogger.Debugf("Initiating [LIST] items ...")
	i := 0
	response := &pb.ListItemsResponse{Items: make([]*pb.Item, s.storage.Count())}
	for item := range s.storage.Iterate() {
		response.Items[i] = item
		i += 1
	}

	asBytes, err := proto.Marshal(response)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to marshal items on [LIST]. Err: %s", err)
		s.stdOutLogger.Error(errMsg)
		response = &pb.ListItemsResponse{
			Items: nil,
			Error: &errMsg,
		}
		asBytes, _ = proto.Marshal(response)
		s.fileLogger.Errorf("%#v", asBytes)
	}
	s.fileLogger.Infof("%#v", asBytes)
	s.stdOutLogger.Debugf("Successful [LIST] of items.")
}

func (s *Service) logCreateErrorResponse(errMsg string) {
	s.stdOutLogger.Error(errMsg)
	response := &pb.CreateItemResponse{
		Item:  nil,
		Error: &errMsg,
	}
	asBytes, _ := proto.Marshal(response)
	s.fileLogger.Errorf("%#v", asBytes)
}

func (s *Service) logDeleteErrorResponse(uuid string, errMsg string) {
	s.stdOutLogger.Error(errMsg)
	response := &pb.DeleteItemResponse{
		ItemUuid: uuid,
		Error:    &errMsg,
	}
	asBytes, _ := proto.Marshal(response)
	s.fileLogger.Errorf("%#v", asBytes)
}

func (s *Service) logGetErrorResponse(errMsg string) {
	s.stdOutLogger.Error(errMsg)
	response := &pb.GetItemResponse{
		Item:  nil,
		Error: &errMsg,
	}
	asBytes, _ := proto.Marshal(response)
	s.fileLogger.Errorf("%#v", asBytes)
}
