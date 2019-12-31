package favorites

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	favoritesapi "github.com/hellupline/samples-golang-grpc-application/pkg/api/favorites"
)

func (s *Service) Create(ctx context.Context, req *favoritesapi.FavoriteCreateRequest) (*favoritesapi.FavoriteCreateResponse, error) {
	result, err := s.datastore.FavoriteCreate(fromPB(req.GetFavorite()))
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	logrus.Infof("favorite %s created", result.Name)

	favorite, err := toPB(result)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &favoritesapi.FavoriteCreateResponse{Favorite: favorite}, nil
}
