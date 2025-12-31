package auth

import (
	"fmt"
	"log"
	"testing"
	"context"

	"github.com/snpavlov/app_aircraft/internal/conf"
)

func HelperTest_CreateAppauth() (IAppAuthService, error) {
    
    // создать экземпляр конфигурации и загрузить данные
    config, err := conf.Configuration{}.New().LoadConfiguration("./../.."); 

    if err != nil {
		log.Fatalf("Не удалось загрузить конфигурацию: %v", err)
		return nil, err		
    }
	
	service, err := AppAuthService{}.NewAppAuthService(config)

	if err != nil {
		log.Fatalf("Ошибка инициализации сервиса 'AppAuthService': %v", err)
		return nil, err
	}

     return service, nil;
}

// TestService_AppauthCreate
func TestService_AppauthCreate(t *testing.T) {

	authsrv, err := HelperTest_CreateAppauth()

    if err != nil {
        t.Errorf("Не удалось создать сервис аутентификации: %v", err)
    }

	if (authsrv != nil){
		fmt.Printf("t: %v\n", t)
	}
	
}

// TestGetExistsByCode
func TestService_GetAppUser(t *testing.T) {
	tests := []struct {
		name   string
		token   string
		expected bool
	}{
		{"GetAppUser_Success", token1, true},
		{"GetAppUser_Failed", "XXX", false},
	}

	ctx := context.Background()
	authsrv, err := HelperTest_CreateAppauth()

    if err != nil {
        t.Errorf("Не удалось создать сервис аутентификации: %v", err)
    }

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			
			appuser, err := authsrv.GetAppUser(ctx, tt.token)
			if err != nil {
				t.Errorf("Ошибка проверки токена 'GetAppUser': %v", err)
			}

			exists := false;
			if appuser.IsAuthenticated {
				exists = true
			}

			if exists != tt.expected {
				t.Errorf("Ошибка в тесте %s", tt.name)
			}
		})
	}	

}

var (
	token1 = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJ2RzhSalk1bjZnWmZ3cGFrSzE3Y2ZLbTllVW5GeHY1eC0yUG01WUs1TG9VIn0.eyJleHAiOjE3Njc0MDU1NjMsImlhdCI6MTc2NzM5MTE2MywianRpIjoib25ydHJvOjQ0MzI5MWIyLTNjN2QtNzhmOC1kYzE4LTBhZTBhNzQyZGI2YiIsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6ODA4Mi9yZWFsbXMvcHJvYmUtYXBwIiwiYXVkIjoiYWNjb3VudCIsInN1YiI6ImUzMThhYTVmLTQ5OGYtNDJjMi1iMjhkLTljY2E4NGNlOGZhZiIsInR5cCI6IkJlYXJlciIsImF6cCI6InByb2JlLWFwcC1jbGllbnQiLCJzaWQiOiIwYzFkYTJkMi02ZGM0LTRkOTgtOGQ1Mi1hODVmYjUyYmI5YTIiLCJhY3IiOiIxIiwiYWxsb3dlZC1vcmlnaW5zIjpbIi8qIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJkZWZhdWx0LXJvbGVzLXByb2JlLWFwcCIsIm9mZmxpbmVfYWNjZXNzIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJzY29wZSI6InByb2ZpbGUgZW1haWwiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwibmFtZSI6IlRlc3QgVXNlcjAxIiwicHJlZmVycmVkX3VzZXJuYW1lIjoidGVzdHVzZXIwMSIsImdpdmVuX25hbWUiOiJUZXN0IiwiZmFtaWx5X25hbWUiOiJVc2VyMDEiLCJlbWFpbCI6InRlc3R1c2VyMDFAcHJvYmVhcHAuZGV2Lm9yZyJ9.fkpaA2IpSzyvqf6fzMk5qDI3BiwILDyZorJV09_zGkOs6pz9ZE153tci4LYJOoUFBPp6vf-B2-dHcx8NNhpeArQIH17lA0YYkS7OXHgKg6WNdLRnOgdUf_yA8-v3v6qfIV5nDEeaoz4Y-EqUURxQKCbZla45sRLxi2ymK5pEUmPqlVDmHcqGGU2qONAJi6Q4H1Jz_7z7d0Qd4ZYpk1liWwwTFzsTdZ-vCzybpqMaHGVAXLGuzC9NkJmjk0vxkOY8nnK8yk0PubO8xzXDgRMXLiNL3mwoPZQeMJN8244odJ6y7QyOG6fnWIJNNqJz8bVTSNz_zmG3wJHcB4bM4EyQNg")
