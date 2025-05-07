package repository

import "github.com/google/wire"

var RepositoryProvider = wire.NewSet(
	new(UserRepositoryImpl),
	wire.Bind(new(UserRepository), new(*UserRepositoryImpl)),
)
