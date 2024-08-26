package app

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/BelyaevEI/microservices_auth/internal/config"
	descAuth "github.com/BelyaevEI/microservices_auth/pkg/auth_v1"
	desc "github.com/BelyaevEI/microservices_auth/pkg/user_v1"
	"github.com/BelyaevEI/platform_common/pkg/closer"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

var (
	configPath string
)

func init() {
	configPath = os.Getenv("CONFIG_PATH")
}

// App represents the app.
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

// NewApp creates a new app.
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// RunConsumerForCreateUser create user infinity
func (a *App) RunConsumerForCreateUser(ctx context.Context) error {

	err := a.serviceProvider.UserSaverConsumer(ctx).RunConsumer(ctx)
	if err != nil {
		log.Printf("failed to run consumer: %s", err.Error())
	}

	return nil
}

func gracefulShutdown(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-waitSignal():
		log.Println("terminating: via signal")
	}

	cancel()
	if wg != nil {
		wg.Wait()
	}
}

func waitSignal() chan os.Signal {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	return sigterm
}

// Run runs the app.
func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	ctx, cancel := context.WithCancel(ctx)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			log.Fatalf("failed to run grpc server: %v", err)
		}

	}()

	go func() {
		defer wg.Done()

		err := a.RunConsumerForCreateUser(ctx)
		if err != nil {
			log.Fatalf("failed to consume saver: %s", err.Error())
		}
	}()

	gracefulShutdown(ctx, cancel, wg)
	return a.runGRPCServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(configPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {

	creds, err := credentials.NewServerTLSFromFile("../service.pem", "../service.key")
	if err != nil {
		log.Fatalf("failed to load TLS keys: %v", err)
	}

	a.grpcServer = grpc.NewServer(grpc.Creds(creds))

	reflection.Register(a.grpcServer)

	desc.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserImpl(ctx))
	descAuth.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthImpl(ctx))
	// descAccess.RegisterAccessV1Server(a.grpcServer, a.serviceProvider.AccessImpl(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	log.Printf("GRPC server is running on %s", a.serviceProvider.GRPCConfig().Address())

	list, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}
