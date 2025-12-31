package auth

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	//"golang.org/x/oauth2"

	"github.com/snpavlov/app_aircraft/internal/conf"
	"github.com/snpavlov/app_aircraft/internal/util"
)

// Определяем интерфейс репозитория IAppAuthService
type IAppAuthService interface {
	GetAppUser(ctx context.Context, token string) (AppUser, error)
	GetJwtToken(authHeader string) (string, error)
}

type AppAuthService struct {
	Provider *oidc.Provider
	Verifier *oidc.IDTokenVerifier
}

func (service AppAuthService) NewAppAuthService(config conf.IConfiguration) (IAppAuthService, error) {
	ctx := context.Background()

	authority, err := config.GetAuthorityAddress()
	if err != nil {
		log.Fatalf("Не удалось получить адрес OIDC из конфигурации: %v", err)
	}

	service.Provider, err = oidc.NewProvider(ctx, authority)
	if err != nil {
		log.Fatalf("Не удалось создать OIDC провайдер: %v", err)
	}

	// Создаем верификатор токенов.
	// Указываем ожидаемого клиента (audience). Здесь это ID клиента в Keycloak.
	service.Verifier = service.Provider.Verifier(&oidc.Config{
		SkipClientIDCheck: true, // ID вашего клиента в Keycloak
	})

	return service, nil
}

func (service AppAuthService) GetAppUser(ctx context.Context, token string) (AppUser, error) {
	var appuser = AppUser{IsAuthenticated: false}

    // Верификация токена
    oauth2Token, err := service.Verifier.Verify(ctx, token)
    if err != nil {
        return appuser, fmt.Errorf("не валидный OIDC токен: %v", err)
    }

    var jclaims = JsonClaims{}
    if err := oauth2Token.Claims(&jclaims); err != nil {
        log.Fatalf("Не удалось получить клеймы из OIDC токена: %v", err)
        return appuser, fmt.Errorf("не удалось получить клеймы из OIDC токена: %v", err)
    }   

	// Собрать роли рилма в массив
	realmRoles := util.Map(jclaims.RealmAccess.Roles, func(p string) string {
		return p
	})

	// Собрать роли ресурсов в массив
	resourceRoles := util.MapMany2[string, struct {Roles []string `json:"roles"`}, string](
		jclaims.ResourceAccess, 
		func(p struct {Roles []string `json:"roles"`}) []string {
			return p.Roles
		},
		func(p string) string {
			return p
		})

	roles := util.DistictSet(util.Concat(realmRoles, resourceRoles),
			func(p string) string {
				return p
			})

	appuser = AppUser{
		IsAuthenticated: true,
		UserId: jclaims.UserId,
		UserName: jclaims.Username,
		FirstName: jclaims.FirstName,
		LastName: jclaims.LastName,
		Email: jclaims.Email,
		Roles: roles,
	}

	return appuser, nil
}

func (service AppAuthService) GetJwtToken(authHeader string) (string, error) {

	// Проверяем формат: Bearer <token>
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
	    return "", fmt.Errorf("Неверный формат заголовка. Ожидается: Bearer <token>")
	}

	tokenString := parts[1]

	return tokenString, nil
}