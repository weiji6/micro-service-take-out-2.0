//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
	"user/internal/handler"
	"user/internal/repository"
	"user/internal/service"
)

var ProviderSet = wire.NewSet(
	handler.HandlerProvider,
	service.ServiceProvider,
	repository.RepositoryProvider,
)
