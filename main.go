package main

import (
	"github.com/djumanoff/amqp"
	users_lib "github.com/kirigaikabuto/recommendation-system-users-store"
	setdata_common "github.com/kirigaikabuto/setdata-common"
)

var (
	postgresUser         = "setdatauser"
	postgresPassword     = "123456789"
	postgresDatabaseName = "recommendation_system"
	postgresHost         = "localhost"
	postgresPort         = 5432
	postgresParams       = "sslmode=disable"
	amqpUrl              = "amqp://localhost:5672"
)

func main() {
	config := users_lib.PostgresConfig{
		Host:     postgresHost,
		Port:     postgresPort,
		User:     postgresUser,
		Password: postgresPassword,
		Database: postgresDatabaseName,
		Params:   postgresParams,
	}
	store, err := users_lib.NewPostgresUsersStore(config)
	if err != nil {
		panic(err)
		return
	}
	service := users_lib.NewUserService(store)
	commandHandler := setdata_common.NewCommandHandler(service)
	usersAmqpEndpoints := users_lib.NewUserAmqpEndpoints(commandHandler)
	rabbitConfig := amqp.Config{
		AMQPUrl:  amqpUrl,
		LogLevel: 5,
	}
	serverConfig := amqp.ServerConfig{
		ResponseX: "response",
		RequestX:  "request",
	}

	sess := amqp.NewSession(rabbitConfig)
	err = sess.Connect()
	if err != nil {
		panic(err)
		return
	}
	srv, err := sess.Server(serverConfig)
	if err != nil {
		panic(err)
		return
	}
	srv.Endpoint("users.create", usersAmqpEndpoints.MakeCreateUserAmqpEndpoint())
	srv.Endpoint("users.get", usersAmqpEndpoints.MakeGetUserAmqpEndpoint())
	srv.Endpoint("users.list", usersAmqpEndpoints.MakeListUserAmqpEndpoint())
	srv.Endpoint("users.update", usersAmqpEndpoints.MakeUpdateUserAmqpEndpoint())
	srv.Endpoint("users.delete", usersAmqpEndpoints.MakeDeleteUserAmqpEndpoint())
	srv.Endpoint("users.getByUsernameAndPassword", usersAmqpEndpoints.MakeGetUserByUsernameAndPasswordAmqpEndpoint())
	err = srv.Start()
	if err != nil {
		panic(err)
		return
	}
}
