package auth

type AppUser struct {
    UserId string
    UserName string
    Email string
    FirstName string
    LastName string
    Roles map[string]struct{}
    IsAuthenticated bool
}

func(user AppUser) HasAnyRole(roles ...string) bool {
    for _, role := range roles {
        _, exists := user.Roles[role]
        if exists {
            return true
        }
    }
    return false
}

func(user AppUser) HasAllRole(roles ...string) bool {
    for _, role := range roles {
        _, exists := user.Roles[role]
        if !exists {
            return false
        }
    }
    return true
}

type JsonClaims struct {
    Email     string `json:"email"`
    UserId    string `json:"sub"`
    Username  string `json:"preferred_username"`
    FirstName string `json:"given_name"`
    LastName  string `json:"family_name"`
    ClientId  string `json:"azp"`
    RealmAccess struct {
        Roles []string `json:"roles"`
    } `json:"realm_access"`
    ResourceAccess map[string]struct {
        Roles []string `json:"roles"`
    } `json:"resource_access"`
}

const (
	AppRole_Reader = "probe-reader"
    AppRole_Contrib = "probe-contrib"
	AppRole_Owner = "probe-owner"
    AppRole_Admim = "probe-admin"    
)