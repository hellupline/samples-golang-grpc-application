package main

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/goware/statik/fs"
	"github.com/hellupline/samples-golang-grpc-application/internal/static/tlsdata"
	favoritesapi "github.com/hellupline/samples-golang-grpc-application/pkg/api/favorites"
	"github.com/hellupline/samples-golang-grpc-server/tlsconfig"
	"github.com/jedib0t/go-pretty/table"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	cmdFavoritesSearch = &cobra.Command{
		Use: "search",
		Run: favoritesSearchRun,
	}
	cmdFavoritesCreate = &cobra.Command{
		Use:  "create",
		Args: cobra.MinimumNArgs(1),
		Run:  favoritesCreateRun,
	}
	cmdFavoritesRead = &cobra.Command{
		Use:  "read",
		Args: cobra.MinimumNArgs(1),
		Run:  favoritesReadRun,
	}
	cmdFavoritesUpdate = &cobra.Command{
		Use:  "update",
		Args: cobra.MinimumNArgs(1),
		Run:  favoritesUpdateRun,
	}
	cmdFavoritesDelete = &cobra.Command{
		Use:  "delete",
		Args: cobra.MinimumNArgs(1),
		Run:  favoritesDeleteRun,
	}

	cmdFavorites = &cobra.Command{Use: "favorites"}

	rootCmd = &cobra.Command{Use: "app"}

	grpcAddr = rootCmd.PersistentFlags().String("grpc-addr", "localhost:50051", "endpoint of the gRPC service")
	useTLS   = rootCmd.PersistentFlags().Bool("tls", false, "Whether to use a secure connection")
)

func init() {
	cmdFavorites.AddCommand(
		cmdFavoritesSearch,
		cmdFavoritesCreate,
		cmdFavoritesRead,
		cmdFavoritesUpdate,
		cmdFavoritesDelete,
	)
	rootCmd.AddCommand(cmdFavorites)
}

func main() {
	rootCmd.Execute()
}

func getFavoritesClient(ctx context.Context) (favoritesapi.FavoritesClient, error) {
	tlsFileSystem, err := fs.New(tlsdata.Asset)
	if err != nil {
		return nil, fmt.Errorf("error opening tlsFileSystem: %w", err)
	}
	tlsConfig, err := tlsconfig.LoadKeyPair(tlsFileSystem)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.DialContext(
		ctx,
		*grpcAddr,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}
	return favoritesapi.NewFavoritesClient(conn), nil

}

func favoritesSearchRun(cmd *cobra.Command, args []string) {
	if err := favoritesSearch(); err != nil {
		logrus.WithError(err).Error("failed to create favorites")
	}
}

func favoritesSearch() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := getFavoritesClient(ctx)
	if err != nil {
		return err
	}
	msg := &favoritesapi.FavoriteSearchRequest{}
	resp, err := client.Search(ctx, msg)
	if err != nil {
		return err
	}

	tw := table.NewWriter()
	tw.AppendHeader(table.Row{"name", "created-at", "updated-at", "count"})
	for _, obj := range resp.GetFavorites() {
		row, err := tableRow(obj)
		if err != nil {
			return err
		}
		tw.AppendRow(row)
	}
	fmt.Println(tw.Render())
	return nil
}

func favoritesCreateRun(cmd *cobra.Command, args []string) {
	if err := favoritesCreate(args); err != nil {
		logrus.WithError(err).Error("failed to create favorites")
	}
}

func favoritesCreate(names []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := getFavoritesClient(ctx)
	if err != nil {
		return err
	}
	tw := table.NewWriter()
	tw.AppendHeader(table.Row{"name", "created-at", "updated-at", "count"})
	for _, name := range names {
		msg := &favoritesapi.FavoriteCreateRequest{Favorite: &favoritesapi.FavoriteMessage{
			Name: name,
		}}
		resp, err := client.Create(ctx, msg)
		if err != nil {
			return err
		}
		row, err := tableRow(resp.GetFavorite())
		if err != nil {
			return err
		}
		tw.AppendRow(row)
	}
	fmt.Println(tw.Render())
	return nil
}

func favoritesReadRun(cmd *cobra.Command, args []string) {
	if err := favoritesRead(args); err != nil {
		logrus.WithError(err).Error("failed to read favorites")
	}
}

func favoritesRead(names []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := getFavoritesClient(ctx)
	if err != nil {
		return err
	}
	tw := table.NewWriter()
	tw.AppendHeader(table.Row{"name", "created-at", "updated-at", "count"})
	for _, name := range names {
		msg := &favoritesapi.FavoriteReadRequest{Name: name}
		resp, err := client.Read(ctx, msg)
		if err != nil {
			return err
		}
		row, err := tableRow(resp.GetFavorite())
		if err != nil {
			return err
		}
		tw.AppendRow(row)
	}
	fmt.Println(tw.Render())
	return nil
}

func favoritesUpdateRun(cmd *cobra.Command, args []string) {
	if err := favoritesUpdate(args); err != nil {
		logrus.WithError(err).Error("failed to update favorites")
	}
}

func favoritesUpdate(names []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := getFavoritesClient(ctx)
	if err != nil {
		return err
	}
	tw := table.NewWriter()
	tw.AppendHeader(table.Row{"name", "created-at", "updated-at", "count"})
	for _, name := range names {
		msg := &favoritesapi.FavoriteUpdateRequest{Name: name}
		resp, err := client.Update(ctx, msg)
		if err != nil {
			return err
		}
		row, err := tableRow(resp.GetFavorite())
		if err != nil {
			return err
		}
		tw.AppendRow(row)
	}
	fmt.Println(tw.Render())
	return nil
}

func favoritesDeleteRun(cmd *cobra.Command, args []string) {
	if err := favoritesDelete(args); err != nil {
		logrus.WithError(err).Error("failed to delete favorites")
	}
}

func favoritesDelete(names []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := getFavoritesClient(ctx)
	if err != nil {
		return err
	}
	for _, name := range names {
		msg := &favoritesapi.FavoriteDeleteRequest{Name: name}
		if _, err := client.Delete(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

func tableRow(obj *favoritesapi.FavoriteMessage) (table.Row, error) {
	createdAt, err := ptypes.Timestamp(obj.GetCreatedAt())
	if err != nil {
		return nil, err
	}
	updatedAt, err := ptypes.Timestamp(obj.GetUpdatedAt())
	if err != nil {
		return nil, err
	}
	return table.Row{obj.GetName(), createdAt, updatedAt, obj.GetCount()}, nil
}
