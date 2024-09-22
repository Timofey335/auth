package app

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/Timofey335/platform_common/pkg/closer"
	"github.com/fatih/color"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/Timofey335/auth/internal/config"
	"github.com/Timofey335/auth/internal/interceptor"
	desc "github.com/Timofey335/auth/pkg/auth_v1"
)

// App - структура с объектами serviceProvider и grpcServer
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
}

// NewApp - создает объект структуры App и вызывает функцию initDeps
func NewApp(ctx context.Context, cfg string) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx, cfg); err != nil {
		return nil, err
	}

	return a, nil
}

// Run - запускает GRPC и HTTP серверы
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			log.Fatalf("failed to run GRPC server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context, cfg string) error {
	inits := []func(context.Context, string) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
	}

	for _, f := range inits {
		if err := f(ctx, cfg); err != nil {
			return err
		}

	}

	return nil
}

func (a *App) initConfig(_ context.Context, cfg string) error {
	if err := config.Load(cfg); err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context, _ string) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context, _ string) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(interceptor.ValidateInterceptor),
	)

	reflection.Register(a.grpcServer)

	desc.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.ServImplementation(ctx))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context, _ string) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := desc.RegisterAuthV1HandlerFromEndpoint(ctx, mux, a.serviceProvider.GRPCConfig().Address(), opts)
	if err != nil {
		return err
	}

	a.httpServer = &http.Server{
		Addr:    a.serviceProvider.HTTPConfig().Address(),
		Handler: mux,
	}

	return nil
}

func (a *App) runGRPCServer() error {
	log.Printf(color.BlueString("GRPC server is running on %s", a.serviceProvider.GRPCConfig().Address()))

	list, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	if err = a.grpcServer.Serve(list); err != nil {
		return err
	}

	return nil
}

func (a *App) runHTTPServer() error {
	log.Printf(color.BlueString("HTTP server is running on %s", a.serviceProvider.HTTPConfig().Address()))

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
