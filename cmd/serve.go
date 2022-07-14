package cmd

import (
	"context"
	"fmt"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-funcards/funapi/docs"
	v1AuthzService "github.com/go-funcards/funapi/internal/client/authz_service/v1"
	v1BoardService "github.com/go-funcards/funapi/internal/client/board_service/v1"
	v1CardService "github.com/go-funcards/funapi/internal/client/card_service/v1"
	v1CategoryService "github.com/go-funcards/funapi/internal/client/category_service/v1"
	v1TagService "github.com/go-funcards/funapi/internal/client/tag_service/v1"
	v1UserService "github.com/go-funcards/funapi/internal/client/user_service/v1"
	"github.com/go-funcards/funapi/internal/config"
	"github.com/go-funcards/funapi/internal/gin/httputil"
	"github.com/go-funcards/funapi/internal/gin/middleware"
	"github.com/go-funcards/funapi/internal/handlers"
	"github.com/go-funcards/funapi/internal/handlers/v1/boards"
	"github.com/go-funcards/funapi/internal/handlers/v1/cards"
	"github.com/go-funcards/funapi/internal/handlers/v1/categories"
	"github.com/go-funcards/funapi/internal/handlers/v1/members"
	"github.com/go-funcards/funapi/internal/handlers/v1/session"
	"github.com/go-funcards/funapi/internal/handlers/v1/tags"
	"github.com/go-funcards/funapi/internal/handlers/v1/users"
	"github.com/go-funcards/graceful"
	"github.com/go-funcards/token-redis"
	"github.com/go-funcards/validate"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/spf13/cobra"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve Rest API Gateway",
	Long:  "Serve Rest API Gateway",
	RunE:  executeServeCommand,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	if v, ok := validate.Default.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" || name == "" {
				if name = fld.Tag.Get("uri"); len(name) > 0 {
					return name
				}
				if name = fld.Tag.Get("ctx"); len(name) > 0 {
					return name
				}
				if name = fld.Tag.Get("form"); len(name) > 0 {
					return name
				}
				return ""
			}
			return name
		})
	}
}

// @title FunCards API
// @version 0.0.1-alpha
// @description REST API for FunCards App

// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func executeServeCommand(cmd *cobra.Command, _ []string) error {
	cfg, err := config.GetConfig(globalFlags.ConfigFile)
	if err != nil {
		return err
	}

	logger, err := cfg.Log.BuildLogger(cfg.Debug)
	if err != nil {
		return err
	}
	logger.Debug("config and logger initialized")

	logger.Debug("parsing redis url")
	opt, err := redis.ParseURL(cfg.Redis.URI)
	if err != nil {
		return err
	}

	logger.Debug("initializing redis client")
	rdb := redis.NewClient(opt)
	defer rdb.Close()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(grpc_zap.UnaryClientInterceptor(logger)),
	}

	authzPool, err := cfg.Services.Authz.NewPool(cmd.Context(), opts...)
	if err != nil {
		return err
	}
	defer authzPool.Close()

	userPool, err := cfg.Services.User.NewPool(cmd.Context(), opts...)
	if err != nil {
		return err
	}
	defer userPool.Close()

	boardPool, err := cfg.Services.Board.NewPool(cmd.Context(), opts...)
	if err != nil {
		return err
	}
	defer boardPool.Close()

	tagPool, err := cfg.Services.Tag.NewPool(cmd.Context(), opts...)
	if err != nil {
		return err
	}
	defer tagPool.Close()

	categoryPool, err := cfg.Services.Category.NewPool(cmd.Context(), opts...)
	if err != nil {
		return err
	}
	defer categoryPool.Close()

	cardPool, err := cfg.Services.Card.NewPool(cmd.Context(), opts...)
	if err != nil {
		return err
	}
	defer cardPool.Close()

	generator, err := cfg.JWT.Signer.Generator()
	if err != nil {
		return err
	}

	logger.Debug("initializing token verifier")
	verifier, err := cfg.JWT.Verifier.Verifier()
	if err != nil {
		return err
	}

	logger.Debug("initializing token service")
	tokenService := tokenredis.New(cfg.RefreshToken, generator, rdb)

	authorize := middleware.Authorize(verifier, cfg.RefreshToken.TokenType)

	subjectService := &v1AuthzService.SubjectService{Pool: authzPool}
	checkerService := &v1AuthzService.CheckerService{Pool: authzPool}
	userService := &v1UserService.UserService{Pool: userPool}
	boardService := &v1BoardService.BoardService{Pool: boardPool}
	tagService := &v1TagService.TagService{Pool: tagPool}
	categoryService := &v1CategoryService.CategoryService{Pool: categoryPool}
	cardService := &v1CardService.CardService{Pool: cardPool}

	sessionHandler := &session.Handler{
		UserService:    userService,
		SubjectService: subjectService,
		TokenService:   tokenService,
		Log:            logger,
	}

	userHandler := &users.Handler{
		UserService: userService,
		IsGranted:   httputil.IsGranted(checkerService, "USER"),
		Log:         logger,
	}

	tagHandler := &tags.Handler{
		BaseBoard: &handlers.BaseBoard{
			BoardService: boardService,
			IsGrantedFn:  httputil.IsGranted(checkerService, "TAG"),
		},
		TagService: tagService,
		Log:        logger,
	}

	categoryHandler := &categories.Handler{
		BaseBoard: &handlers.BaseBoard{
			BoardService: boardService,
			IsGrantedFn:  httputil.IsGranted(checkerService, "CATEGORY"),
		},
		CategoryService: categoryService,
		Log:             logger,
	}

	cardHandler := &cards.Handler{
		BaseBoard: &handlers.BaseBoard{
			BoardService: boardService,
			IsGrantedFn:  httputil.IsGranted(checkerService, "CARD"),
		},
		CardService: cardService,
		Log:         logger,
	}

	memberHandler := &members.Handler{
		BaseBoard: &handlers.BaseBoard{
			BoardService: boardService,
			IsGrantedFn:  httputil.IsGranted(checkerService, "BOARD"),
		},
		SubjectService: subjectService,
		Log:            logger,
	}

	boardHandler := &boards.Handler{
		BaseBoard: &handlers.BaseBoard{
			BoardService: boardService,
			IsGrantedFn:  memberHandler.IsGrantedFn,
		},
		MemberHandler:  memberHandler,
		SubjectService: subjectService,
		Log:            logger,
	}

	r := router(cfg.Debug, logger)
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			sessionHandler.Register(v1)

			authorized := v1.Group("", authorize)
			{
				userHandler.Register(authorized)
				boardHandler.Register(authorized)
				tagHandler.Register(authorized)
				categoryHandler.Register(authorized)
				cardHandler.Register(authorized)
			}
		}
	}

	if cfg.Swagger.Enable {
		docs.SwaggerInfo.BasePath = "/api/v1"
		r.GET(cfg.Swagger.Path, ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.PersistAuthorization(true)))
	}

	return serve(r, cfg, logger)
}

func router(debug bool, logger *zap.Logger) *gin.Engine {
	logger.Debug("initializing gin")

	mode := gin.ReleaseMode
	if debug {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)
	binding.Validator = nil

	r := gin.New()
	r.HandleMethodNotAllowed = true
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	r.Use(middleware.APIError())

	return r
}

func serve(router http.Handler, cfg config.Config, logger *zap.Logger) error {
	logger.Info("bind application to host and port", zap.String("addr", cfg.Server.Addr))
	lst, err := net.Listen("tcp", cfg.Server.Addr)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info("application initialized and started")

	go func() {
		if err := srv.Serve(lst); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", zap.Error(err))
		}
	}()

	graceful.Exit(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err = srv.Shutdown(ctx); err != nil {
			logger.Fatal("server shutdown", zap.Error(err))
		}
		select {
		case <-ctx.Done():
			fmt.Println("Timeout of 5 seconds.")
		}
	})

	return nil
}
