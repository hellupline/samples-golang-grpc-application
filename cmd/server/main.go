package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/goware/statik/fs"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
	"google.golang.org/grpc"

	"github.com/hellupline/samples-golang-grpc-server/server"
	"github.com/hellupline/samples-golang-grpc-server/storage"
	"github.com/hellupline/samples-golang-grpc-server/tlsconfig"

	"github.com/hellupline/samples-golang-grpc-application/internal/service/favorites"
	"github.com/hellupline/samples-golang-grpc-application/internal/static/openapidata"
	"github.com/hellupline/samples-golang-grpc-application/internal/static/tlsdata"
	favoritesapi "github.com/hellupline/samples-golang-grpc-application/pkg/api/favorites"
)

const (
	dbName = "application.db"
)

var (
	grpcAddr = flag.String("grpc-addr", "localhost:50051", "endpoint of the gRPC service")
	httpAddr = flag.String("http-addr", "localhost:8080", "endpoint of the http service")
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		logrus.Error(err)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		signal.Notify(c, syscall.SIGTERM)
		signal.Notify(c, syscall.SIGINT)
		<-c
		cancel()
	}()

	db, err := bbolt.Open(dbName, 0600, nil)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()
	favoritesService := favorites.New(storage.New(db))

	tlsConfig, err := loadkeys()
	if err != nil {
		return err
	}
	openapiData, err := loadOpenApiData()
	if err != nil {
		return err
	}
	s := server.New(*grpcAddr, *httpAddr, tlsConfig, openapiData, http.DefaultServeMux)
	if err := s.StartGrpcServer(func(s *grpc.Server) error {
		favoritesapi.RegisterFavoritesServer(s, favoritesService)
		return nil
	}); err != nil {
		return err
	}
	if err := s.StartHttpServer(func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
		if err := favoritesapi.RegisterFavoritesHandler(ctx, mux, conn); err != nil {
			return fmt.Errorf("error registering Posts handler: %w", err)
		}
		return nil
	}); err != nil {
		return err
	}
	defer s.Close()

	<-ctx.Done()
	return ctx.Err()
}

func loadkeys() (*tls.Config, error) {
	tlsFileSystem, err := fs.New(tlsdata.Asset)
	if err != nil {
		return nil, fmt.Errorf("error opening tlsFileSystem: %w", err)
	}
	tlsConfig, err := tlsconfig.LoadKeyPair(tlsFileSystem)
	if err != nil {
		return nil, err
	}
	return tlsConfig, nil
}

func loadOpenApiData() ([]byte, error) {
	fileSystem, err := fs.New(openapidata.Asset)
	if err != nil {
		return nil, fmt.Errorf("error opening openapiFileSystem: %w", err)
	}
	swaggerFile, err := fileSystem.Open("/apidocs.swagger.json")
	if err != nil {
		return nil, fmt.Errorf("error opening swagger file: %w", err)
	}
	data, err := ioutil.ReadAll(swaggerFile)
	if err != nil {
		return nil, fmt.Errorf("error reading swagger file: %w", err)
	}
	return data, nil
}
