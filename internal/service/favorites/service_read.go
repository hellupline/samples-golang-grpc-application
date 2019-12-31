package favorites

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	favoritesapi "github.com/hellupline/samples-golang-grpc-application/pkg/api/favorites"
)

func (s *Service) Read(ctx context.Context, req *favoritesapi.FavoriteReadRequest) (*favoritesapi.FavoriteReadResponse, error) {
	result, err := s.datastore.FavoriteRead(req.GetName())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	logrus.Infof("favorite %s read", result.Name)

	favorite, err := toPB(result)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &favoritesapi.FavoriteReadResponse{Favorite: favorite}, nil
}
