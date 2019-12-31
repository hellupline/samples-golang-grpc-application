package favorites

import (
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/hellupline/samples-golang-grpc-server/storage"
	"github.com/sirupsen/logrus"

	favoritesapi "github.com/hellupline/samples-golang-grpc-application/pkg/api/favorites"
)

var logger = logrus.WithField("module", "service/favorites")

type Service struct {
	datastore *Datastore

	favoritesapi.UnimplementedFavoritesServer
}

func New(s *storage.Storage) *Service {
	return &Service{datastore: NewDatastore(s)}
}

func toPB(o *Favorite) (*favoritesapi.FavoriteMessage, error) {
	createdAt, err := ptypes.TimestampProto(o.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to convert time to protobuf: %+v: %w", createdAt, err)
	}
	updatedAt, err := ptypes.TimestampProto(o.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to convert time to protobuf: %+v: %w", createdAt, err)
	}
	msg := &favoritesapi.FavoriteMessage{
		Name:      o.Name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Count:     uint64(o.Count),
	}
	return msg, nil
}

func fromPB(o *favoritesapi.FavoriteMessage) Favorite {
	return Favorite{
		Name:  o.GetName(),
		Count: uint(o.GetCount()),
	}
}
