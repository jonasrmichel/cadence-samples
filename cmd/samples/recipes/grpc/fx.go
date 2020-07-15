package main

import (
	"context"
	"errors"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

type EnvironmentLoader struct {
	fx.Out
}

func Load(spec interface{}) error {
	if err := envconfig.Process("", spec); err != nil {
		return err
	}

	return nil
}

type GreeterSettings struct {
	fx.In

	HostPort string `name:"HostPort"`
}

type GreeterSettingsLoader struct {
	fx.Out
	EnvironmentLoader

	HostPort string `name:"HostPort" envconfig:"HOST_PORT"`
}

func (l *GreeterSettingsLoader) Load() error {
	return Load(l)
}

var GreeterSettingsModule = fx.Provide(
	func() (l GreeterSettingsLoader, err error) {
		err = l.Load()

		if l.HostPort == "" {
			err = errors.New("HOST_PORT is undefined")
		}

		return
	},
)

type GreeterClient struct {
	fx.In

	Client pb.GreeterClient `name:"GreeterClient`
}

type GreeterClientLoader struct {
	fx.Out

	Client pb.GreeterClient `name:"GreeterClient`
}

func (l *GreeterClientLoader) Load(s GreeterSettings) (err error) {
	l.Client, err = NewGreeterClient(s.HostPort)
	return
}

func NewGreeterClient(hostPort string) (c pb.GreeterClient, err error) {
	if conn, e := grpc.Dial(hostPort, grpc.WithInsecure(), grpc.WithBlock()); e != nil {
		err = e
	} else {
		c = pb.NewGreeterClient(conn)
	}

	return
}

var GreeterClientModule = fx.Provide(
	func(s GreeterSettings) (l GreeterClientLoader, err error) {
		err = l.Load(s)
		return
	},
)

func MakeGreeterClient() (client pb.GreeterClient, err error) {
	var c GreeterClient
	app := fx.New(
		GreeterSettingsModule,
		GreeterClientModule,
		fx.Populate(&c),
	)

	err = app.Start(context.Background())
	defer app.Stop(context.Background())
	if err != nil {
		return
	}

	client = c.Client

	return
}
