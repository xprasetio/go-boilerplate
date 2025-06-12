package container

import (
	"boilerplate/config"
	"boilerplate/internal/category"
	"boilerplate/internal/user"
	"boilerplate/pkg/database"
	"boilerplate/pkg/logger"
	"boilerplate/pkg/middleware"
	"boilerplate/pkg/redis"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

func NewContainer() (di.Container, error) {
	builder, err := di.NewBuilder()
	if err != nil {
		return di.Container{}, err
	}

	defs := []di.Def{
		{
			Name: ConfigDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				return config.LoadConfig()
			},
		},
		{
			Name: LoggerDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				return logger.NewLogger(), nil
			},
		},
		{
			Name: DBDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(ConfigDefName).(config.Config)
				return database.InitDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
			},
		},
		{
			Name: RedisClientDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(ConfigDefName).(config.Config)
				return redis.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB), nil
			},
		},
		{
			Name: UserServiceDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(ConfigDefName).(config.Config)
				db := ctn.Get(DBDefName).(*gorm.DB)
				logger := ctn.Get(LoggerDefName).(logger.Logger)
				redisClient := ctn.Get(RedisClientDefName).(*redis.RedisClient)
				return user.NewUserService(db, cfg.JWTSecret, logger, redisClient), nil
			},
		},
		{
			Name: CategoryServiceDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				db := ctn.Get(DBDefName).(*gorm.DB)
				return category.NewCategoryService(db), nil
			},
		},
		{
			Name: UserHandlerDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				userService := ctn.Get(UserServiceDefName).(user.UserServiceInterface)
				return user.NewUserHandler(userService), nil
			},
		},
		{
			Name: CategoryHandlerDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				categoryService := ctn.Get(CategoryServiceDefName).(category.CategoryServiceInterface)
				return category.NewCategoryHandler(categoryService), nil
			},
		},
		{
			Name: AuthMiddlewareDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(ConfigDefName).(config.Config)
				userService := ctn.Get(UserServiceDefName).(user.UserServiceInterface)
				return middleware.AuthMiddleware(userService, cfg.JWTSecret), nil
			},
		},
		{
			Name: AdminAuthMiddlewareDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				return middleware.AdminMiddleware(), nil
			},
		},
		{
			Name: ValidatorDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				return validator.New(), nil
			},
		},
		{
			Name: EchoDefName,
			Build: func(ctn di.Container) (interface{}, error) {
				e := echo.New()
				validate := ctn.Get(ValidatorDefName).(*validator.Validate)
				e.Validator = &CustomValidator{validator: validate}
				return e, nil
			},
		},
	}

	if err := builder.Add(defs...); err != nil {
		return di.Container{}, err
	}

	return builder.Build(), nil
}

// CustomValidator adalah custom validator untuk Echo
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
