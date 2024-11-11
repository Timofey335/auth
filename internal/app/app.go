package app

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/Timofey335/auth/internal/config"
	"github.com/Timofey335/auth/internal/interceptor"
	"github.com/Timofey335/auth/internal/logger"
	descAccess "github.com/Timofey335/auth/pkg/access_v1"
	descAuth "github.com/Timofey335/auth/pkg/auth_v1"
	_ "github.com/Timofey335/auth/statik"
	"github.com/Timofey335/platform_common/pkg/closer"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// App - структура App
type App struct {
	serviceProvider  *serviceProvider
	grpcServer       *grpc.Server
	httpServer       *http.Server
	prometheusServer *http.Server
	swaggerServer    *http.Server
}

// NewApp - создает объект структуры App и вызывает функцию initDeps
func NewApp(ctx context.Context, cfg string) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx, cfg); err != nil {
		return nil, err
	}

	return a, nil
}

// Run - запускает GRPC, HTTP и Swagger серверы
func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := &sync.WaitGroup{}
	wg.Add(5)

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

	go func() {
		defer wg.Done()

		err := a.serviceProvider.UserSaverConsumer(ctx).RunConsumer(ctx)
		if err != nil {
			log.Printf("failed to run consumer: %s", err.Error())
		}

	}()

	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			log.Fatalf("failed to run Swagger server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runPrometheusServer()
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context, cfg string) error {
	inits := []func(context.Context, string) error{
		a.initConfig,
		a.initServiceProvider,
		a.initLogger,
		a.initMetrics,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
		a.initPrometheusServer,
	}

	for _, f := range inits {
		if err := f(ctx, cfg); err != nil {
			return err
		}

	}

	return nil
}

func (a *App) initLogger(_ context.Context, _ string) error {
	a.serviceProvider.Logger()

	return nil
}

func (a *App) initConfig(_ context.Context, cfg string) error {
	if err := config.Load(cfg); err != nil {
		return err
	}

	return nil
}

func (a *App) initMetrics(ctx context.Context, _ string) error {
	a.serviceProvider.Metrics(ctx)

	return nil
}

func (a *App) initServiceProvider(_ context.Context, _ string) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context, _ string) error {
	creds, err := credentials.NewServerTLSFromFile("cert/service.pem", "cert/service.key")
	if err != nil {
		log.Fatalf("failed to load TLS keys: %v", err)
	}

	a.grpcServer = grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(
		grpcMiddleware.ChainUnaryServer(
			interceptor.ValidateInterceptor,
			interceptor.MetricsInterceptor,
		),
	),
	// grpc.UnaryInterceptor(
	// 	grpcMidd
	// 	interceptor.ValidateInterceptor),
	// grpc.UnaryInterceptor(interceptor.MetricsInterceptor),
	)

	reflection.Register(a.grpcServer)

	descAuth.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.ServImplementation(ctx))
	descAccess.RegisterAccessV1Server(a.grpcServer, a.serviceProvider.AccessServImplementation(ctx))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context, _ string) error {
	mux := runtime.NewServeMux()

	creds, err := credentials.NewClientTLSFromFile("cert/service.pem", "")
	if err != nil {
		log.Fatalf("could not process the credentials: %v", err)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	err = descAuth.RegisterAuthV1HandlerFromEndpoint(ctx, mux, a.serviceProvider.GRPCConfig().Address(), opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:    a.serviceProvider.HTTPConfig().Address(),
		Handler: corsMiddleware.Handler(mux),
	}

	return nil
}

func (a *App) initSwaggerServer(_ context.Context, _ string) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json"))

	a.swaggerServer = &http.Server{
		Addr:    a.serviceProvider.SwaggerConfig().Address(),
		Handler: mux,
	}

	return nil
}

func (a *App) initPrometheusServer(_ context.Context, _ string) error {
	mux := http.NewServeMux()
	mux.Handle(a.serviceProvider.PromethueusConfig().Path(), promhttp.Handler())
	// mux.Handle("/metrics", promhttp.Handler())

	a.prometheusServer = &http.Server{
		Addr: a.serviceProvider.PromethueusConfig().Address(),
		// Addr:    "localhost:2112",
		Handler: mux,
	}

	return nil
}

func (a *App) runGRPCServer() error {
	logger.Info("GRPC server is running on", zap.String("address", a.serviceProvider.GRPCConfig().Address()))

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
	logger.Info("HTTP server is running on", zap.String("address", a.serviceProvider.HTTPConfig().Address()))

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runPrometheusServer() error {
	logger.Info("Prometheus server is running on", zap.String("address", a.serviceProvider.PromethueusConfig().Address()))

	err := a.prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	logger.Info("Swagger server is running on", zap.String("address", a.serviceProvider.SwaggerConfig().Address()))

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Serving swagger file", zap.String("path", path))

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("Open swagger file", zap.String("path", path))

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		logger.Info("Read swagger file", zap.String("path", path))

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("Write swagger file", zap.String("path", path))

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("Served swagger file", zap.String("path", path))
	}
}
