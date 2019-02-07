package factory

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
	"os"
	"sumwhere/middlewares"
)

func DB(ctx context.Context) xorm.Interface {
	v := ctx.Value(middlewares.ContextDBName)
	if v == nil {
		panic("DB is not exist")
	}
	if db, ok := v.(*xorm.Session); ok {
		return db
	}
	if db, ok := v.(*xorm.Engine); ok {
		return db
	}
	panic("DB is not exist")
}

func Logger(ctx context.Context) *logrus.Entry {
	v := ctx.Value(middlewares.ContextLoggerName)
	if v == nil {
		return logrus.WithFields(logrus.Fields{})
	}
	if logger, ok := v.(*logrus.Entry); ok {
		return logger
	}
	return logrus.WithFields(logrus.Fields{})
}

func Redis(ctx context.Context, name string) *redis.Client {

	if os.Getenv("RELEASE_SYSTEM") != "kubernetes" {
		name = middlewares.ContextSetRedisName
	}

	v := ctx.Value(name)

	if v == nil {
		panic("Redis is not exist")
	}

	if redis, ok := v.(*redis.Client); ok {
		return redis
	}
	panic("Redis is not exist")
}
