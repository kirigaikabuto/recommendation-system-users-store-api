package main

import (
	"github.com/djumanoff/amqp"
	users_lib "github.com/kirigaikabuto/recommendation-system-users-store"
	setdata_common "github.com/kirigaikabuto/setdata-common"
)

var (

)
func main() {
	config := users_lib.PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "setdatauser",
		Password: "123456789",
		Database: "recommendation_system",
		Params:   "sslmode=disable",
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
		Host:     "localhost",
		Port:     5672,
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
