package controllers

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

type Controller struct {
	DB    *gorm.DB
	Cache *redis.Client
}
