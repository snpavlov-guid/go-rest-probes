package server

import (
    "net/http"
	"github.com/gin-gonic/gin"

	"github.com/snpavlov/app_aircraft/internal/auth"
)


// Middleware middleware для проверки JWT токена
func (server AppServer) JwtMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        
		// Извлекаем токен из заголовка Authorization
        authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			 c.Set(server.AppUserKey, auth.AppUser{IsAuthenticated: false})
			 c.Next()
			 return 
		}

		token, err := server.authService.GetJwtToken(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
        		"error": err.Error(),
             })
		}

		appuser, err := server.authService.GetAppUser(c, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
        		"error": err.Error(),
             })
		}

		c.Set(server.AppUserKey, appuser)
        
        c.Next()
    }
}

// Middleware middleware для проверки авторизации пользователя
func (server AppServer) Authorize(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {

        user, exists := c.Get(server.AppUserKey)
        if !exists {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                    "error": "Пользователь не авторизован",
                }) 
            return
        }

        // Проверка типа с помощью type assertion
        appuser, ok := user.(auth.AppUser)
        if !ok {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                    "error": "Данные пользователя не соотвутствуют ожидаемому формату",
                }) 
            return
        }

        if !appuser.IsAuthenticated {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                    "error": "Пользователь не авторизован",
                }) 
            return		
        }

        if len(roles) == 0 {
            c.Next()
            return
        }

        if !appuser.HasAnyRole(roles...) {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                    "error": "Пользователь не авторизован",
                }) 
            return		       
        }

        c.Next()
    }
}

