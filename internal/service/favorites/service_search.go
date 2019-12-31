package favorites

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	favoritesapi "github.com/hellupline/samples-golang-grpc-application/pkg/api/favorites"
)

func (s *Service) Search(ctx context.Context, req *favoritesapi.FavoriteSearchRequest) (*favoritesapi.FavoriteSearchResponse, error) {
	objs, err := s.datastore.FavoriteSearch()
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	logrus.Infof("favorites search")

	favorites := make([]*favoritesapi.FavoriteMessage, len(objs))
	for i, obj := range objs {
		var err error
		favorites[i], err = toPB(obj)
		if err != nil {
			return nil, status.Error(codes.Unknown, err.Error())
		}
	}
	return &favoritesapi.FavoriteSearchResponse{Favorites: favorites}, nil
}
