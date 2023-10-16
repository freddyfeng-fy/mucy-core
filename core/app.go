package core

import (
	"github.com/go-redis/redis/v8"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Application struct {
	Log   *zap.Logger
	DB    *gorm.DB
	Redis *redis.Client
	Nacos *naming_client.INamingClient
}

var App = new(Application)
