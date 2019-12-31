package favorites

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	favoritesapi "github.com/hellupline/samples-golang-grpc-application/pkg/api/favorites"
)

func (s *Service) Update(ctx context.Context, req *favoritesapi.FavoriteUpdateRequest) (*favoritesapi.FavoriteUpdateResponse, error) {
	result, err := s.datastore.FavoriteUpdate(req.GetName(), fromPB(req.GetFavorite()))
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	logrus.Infof("favorite %s updated", result.Name)

	favorite, err := toPB(result)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &favoritesapi.FavoriteUpdateResponse{Favorite: favorite}, nil
}
