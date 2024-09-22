package actions

import (
	"github.com/gobuffalo/pop/v6"
	"net/http"
	"sync"

	"sharaphka/app/infrastructure/jwt"       // Подключаем JWT-сервис
	"sharaphka/app/interfaces/database"      // Репозиторий для пользователей
	authhttp "sharaphka/app/interfaces/http" // Подключаем обработчики аутентификации
	"sharaphka/app/usecases"                 // Используем сценарии использования
	"sharaphka/locales"
	"sharaphka/models"
	"sharaphka/public"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/v3/pop/popmw"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/middleware/forcessl"
	"github.com/gobuffalo/middleware/i18n"
	"github.com/gobuffalo/middleware/paramlogger"
	"github.com/unrolled/secure"
)

var ENV = envy.Get("GO_ENV", "development")

var (
	app     *buffalo.App
	appOnce sync.Once
	T       *i18n.Translator
)

// App инициализирует маршруты и middleware для Buffalo приложения
func App() *buffalo.App {
	appOnce.Do(func() {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_sharaphka_session",
		})

		// Принудительное использование SSL в production
		// app.Use(forceSSL())

		// Логирование параметров запросов
		app.Use(paramlogger.ParameterLogger)

		// Защита от CSRF атак
		//app.Use(csrf.New)

		// Каждое обращение заворачивается в транзакцию pop
		app.Use(popmw.Transaction(models.DB))

		// Настройка и использование переводов
		app.Use(translations())

		// Инициализация зависимостей для JWT
		env := envy.Get("GO_ENV", "development")
		db, err := pop.Connect(env) // Подключение к базе данных
		if err != nil {
			app.Stop(err)
		}

		userRepo := database.NewUserRepositoryDB(db)
		jwtService := jwt.NewJWTService()

		authUseCase := usecases.NewAuthUseCase(userRepo, jwtService)
		authHandler := authhttp.NewAuthHandler(authUseCase)

		// Добавление маршрутов для аутентификации
		app.POST("/login", authHandler.Login)

		// Пример защищённого маршрута
		protected := app.Group("/protected")
		protected.Use(authHandler.AuthMiddleware) // Middleware для проверки JWT
		protected.GET("/profile", authHandler.Profile)

		// Основные маршруты
		app.GET("/", HomeHandler)

		// Статические файлы
		app.ServeFiles("/", http.FS(public.FS()))
	})

	return app
}

// translations загружает файлы локалей и возвращает middleware для перевода
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(locales.FS(), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL принудительно перенаправляет запросы на HTTPS
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
