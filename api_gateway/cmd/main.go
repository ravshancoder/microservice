package main

import (
	"fmt"

	"github.com/casbin/casbin/util"
	"github.com/casbin/casbin/v2"
	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	r "github.com/gomodule/redigo/redis"
	"github.com/microservice/api_gateway/api"
	"github.com/microservice/api_gateway/config"
	"github.com/microservice/api_gateway/pkg/logger"
	"github.com/microservice/api_gateway/services"
	"github.com/microservice/api_gateway/storage/redis"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "api_gateway")

	serviceManager, err := services.NewServiceManager(&cfg)

	log.Info("main: sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.String("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase),
	)

	psqlString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
	)

	a, err := gormadapter.NewAdapter("postgres", psqlString, true)
	if err != nil {
		log.Error("new  adapter error", logger.Error(err))
		return
	}
	
	fmt.Println(a, cfg.AuthConfigPath)
	
	casbinEnforcer, err := casbin.NewEnforcer(cfg.AuthConfigPath, a)
	if err != nil {
		log.Error("new enforcer error", logger.Error(err))
		return
	}
	err = casbinEnforcer.LoadPolicy()
	if err != nil {
		log.Error("casbin error load policy", logger.Error(err))
		return
	}

	casbinEnforcer.GetRoleManager().(*defaultrolemanager.RoleManager).AddMatchingFunc("keyMatch", util.KeyMatch)
	casbinEnforcer.GetRoleManager().(*defaultrolemanager.RoleManager).AddMatchingFunc("keyMatch3", util.KeyMatch3)

	pool := &r.Pool{
		MaxIdle: 80,
		MaxActive: 12000,
		Dial: func() (r.Conn, error) {
			c, err := r.Dial("tcp", fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort))
			if err != nil{
				panic(err.Error())
			}
			return c, err
		},
	}

	server := api.New(api.Option{
		Conf:           cfg,
		Logger:         log,
		ServiceManager: serviceManager,
		RedisRepo:      redis.NewRedisRepo(pool),
		CasbinEnforcer: casbinEnforcer,
	})

	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("failed to run HTTP server: ", logger.Error(err))
		panic(err)
	}

}
