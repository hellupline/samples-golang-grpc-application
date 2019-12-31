package favorites

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	favoritesapi "github.com/hellupline/samples-golang-grpc-application/pkg/api/favorites"
)

func (s *Service) Delete(ctx context.Context, req *favoritesapi.FavoriteDeleteRequest) (*favoritesapi.FavoriteDeleteResponse, error) {
	if err := s.datastore.FavoriteDelete(req.GetName()); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	logrus.Infof("favorite %s deleted", req.GetName())

	return &favoritesapi.FavoriteDeleteResponse{}, nil
}
