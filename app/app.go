package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/alexchebotarsky/heatpump-api/client"
	"github.com/alexchebotarsky/heatpump-api/env"
	"github.com/alexchebotarsky/heatpump-api/server"
)

type App struct {
	Services []Service
	Clients  *Clients
}

func New(env *env.Config) (*App, error) {
	var app App
	var err error

	app.Clients, err = setupClients(env)
	if err != nil {
		return nil, fmt.Errorf("error setting up clients: %v", err)
	}

	app.Services, err = setupServices(env, app.Clients)
	if err != nil {
		return nil, fmt.Errorf("error setting up services: %v", err)
	}

	return &app, nil
}

func (app *App) Launch(ctx context.Context) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	errc := make(chan error, 1)

	for _, service := range app.Services {
		go service.Start(ctx, errc)
	}

	select {
	case <-ctx.Done():
		slog.Debug("Context is cancelled")
	case err := <-errc:
		slog.Error(fmt.Sprintf("Critical service error: %v", err))
	}

	var errors []error

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, service := range app.Services {
		err := service.Stop(ctx)
		if err != nil {
			errors = append(errors, fmt.Errorf("error stopping a service: %v", err))
		}
	}

	err := app.Clients.Close()
	if err != nil {
		errors = append(errors, fmt.Errorf("error closing app clients: %v", err))
	}

	if len(errors) > 0 {
		slog.Error(fmt.Sprintf("Error gracefully shutting down: %v", &client.ErrMultiple{Errs: errors}))
	} else {
		slog.Debug("App has been gracefully shut down")
	}
}

type Service interface {
	Start(context.Context, chan<- error)
	Stop(context.Context) error
}

func setupServices(env *env.Config, clients *Clients) ([]Service, error) {
	var services []Service

	server := server.New(env.Host, env.Port, server.Clients{})
	services = append(services, server)

	return services, nil
}

type Clients struct {
}

func setupClients(env *env.Config) (*Clients, error) {
	var c Clients
	var _ error

	return &c, nil
}

func (c *Clients) Close() error {
	var errors []error

	if len(errors) > 0 {
		return &client.ErrMultiple{Errs: errors}
	}

	return nil
}
