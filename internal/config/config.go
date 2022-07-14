package config

import (
	"context"
	"github.com/go-funcards/envconfig"
	"github.com/go-funcards/grpc-pool"
	"github.com/go-funcards/jwt"
	"github.com/go-funcards/logger"
	"github.com/go-funcards/token"
	"google.golang.org/grpc"
	"sync"
)

type ServiceConfig struct {
	Addr string          `yaml:"address" env:"ADDRESS" env-required:"true"`
	Pool grpcpool.Config `yaml:"pool" env-prefix:"POOL_"`
}

func (cfg ServiceConfig) NewPool(ctx context.Context, opts ...grpc.DialOption) (*grpcpool.Pool, error) {
	return cfg.Pool.NewPool(ctx, func(ctx1 context.Context) (*grpc.ClientConn, error) {
		return grpc.DialContext(ctx1, cfg.Addr, opts...)
	})
}

type ServicesConfig struct {
	Authz    ServiceConfig `yaml:"authz" env-prefix:"AUTHZ_"`
	User     ServiceConfig `yaml:"user" env-prefix:"USER_"`
	Board    ServiceConfig `yaml:"board" env-prefix:"BOARD_"`
	Category ServiceConfig `yaml:"category" env-prefix:"CATEGORY_"`
	Tag      ServiceConfig `yaml:"tag" env-prefix:"TAG_"`
	Card     ServiceConfig `yaml:"card" env-prefix:"CARD_"`
}

type SwaggerConfig struct {
	Enable bool   `yaml:"enable" env:"ENABLE" env-default:"false"`
	Path   string `yaml:"path" env-default:"/swagger/*any"`
}

type RedisConfig struct {
	URI string `yaml:"uri" env:"URI"`
}

type ServerConfig struct {
	Addr string `yaml:"address" env:"ADDRESS" env-default:":80"`
}

type Config struct {
	Debug        bool           `yaml:"debug" env:"DEBUG_MODE" env-default:"false"`
	Log          logger.Config  `yaml:"log" env-prefix:"LOG_"`
	Server       ServerConfig   `yaml:"server" env-prefix:"SERVER_"`
	Redis        RedisConfig    `yaml:"redis" env-prefix:"REDIS_"`
	Services     ServicesConfig `yaml:"services" env-prefix:"SERVICE_"`
	Swagger      SwaggerConfig  `yaml:"swagger" env-prefix:"SWAGGER_"`
	RefreshToken token.Config   `yaml:"refresh_token" env-prefix:"REFRESH_TOKEN_"`
	JWT          struct {
		Signer   jwt.SignerConfig   `yaml:"signer" env-prefix:"SIGNER_"`
		Verifier jwt.VerifierConfig `yaml:"verifier" env-prefix:"VERIFIER_"`
	} `yaml:"jwt" env-prefix:"JWT_"`
}

var (
	cfg  Config
	once sync.Once
)

func GetConfig(path string) (Config, error) {
	var err error
	once.Do(func() {
		err = envconfig.ReadConfig(path, &cfg)
	})
	return cfg, err
}
