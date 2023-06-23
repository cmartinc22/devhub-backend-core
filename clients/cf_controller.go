package clients

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/cmartinc22/devhub-backend-core/models"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/pedidosya/peya-go/logs"
)

var lock = &sync.Mutex{}

type CFControllerSpec interface {
	Validate(ctx context.Context, accessToken string) (*oidc.IDToken, error)
	ValidateCFToken(ctx context.Context, accessToken string, scope string) (*models.AuthResult, error)
}

type CFClient struct {
	enabled  bool
	config   *models.CFConfiguration
	cfConfig *oidc.Config
	verifier *oidc.IDTokenVerifier
}

var cFClientInstance *CFClient

const teamDomain = "https://deliveryhero.cloudflareaccess.com"

func getCFClient() *CFClient {
	if cFClientInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if cFClientInstance == nil {
			fmt.Println("Creating single instance now.")
			c := readCFConfig()
			cFClientInstance = &CFClient{
				enabled: c.Enabled,
				config:  c,
				cfConfig: &oidc.Config{
					ClientID: c.AUD,
				},
			}
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return cFClientInstance
}

/*
func NewCFClient() CFControllerSpec {
	c := readCFConfig()
	return &CFClient{
		enabled: c.Enabled,
		config:  c,
		cfConfig: &oidc.Config{
			ClientID: c.AUD,
		},
	}
}
*/

func (c *CFClient) Validate(ctx context.Context, accessToken string) (*oidc.IDToken, error) {
	if c.verifier == nil {
		keySet := oidc.NewRemoteKeySet(ctx, fmt.Sprintf("%s/cdn-cgi/access/certs", teamDomain))
		c.verifier = oidc.NewVerifier(teamDomain, keySet, c.cfConfig)
	}
	return c.verifier.Verify(ctx, accessToken)
}

func (c *CFClient) ValidateCFToken(ctx context.Context, accessToken string, scope string) (*models.AuthResult, error) {
	if !c.enabled {
		return &models.AuthResult{
			IsValid: true,
			Enabled: false,
		}, nil
	}

	if len(accessToken) == 0 {
		return &models.AuthResult{
			IsValid: false,
			Enabled: true,
		}, errors.New("authorization CF token should not be empty")
	}

	token, err := c.Validate(ctx, accessToken)
	if err != nil {
		return &models.AuthResult{
			IsValid: false,
			Enabled: true,
		}, err
	}

	identity := fmt.Sprintf("%s:%s", "subject", token.Subject)

	var claims struct {
		Email string `json:"email"`
	}
	if err := token.Claims(&claims); err != nil {
		logs.Errorf("[engine] couldn't get the claim email: %w", err)
	}

	if claims.Email != "" {
		identity = fmt.Sprintf("%s:%s", "user", claims.Email)
	}

	return &models.AuthResult{
		IsValid:  true,
		Enabled:  true,
		Identity: &identity,
	}, nil
}

// MIDDLEAWRE
func CFAuthCheckMiddleware(c CFClient, r *http.Request, scope string) (authResult *models.AuthResult, error *models.CustomError) {
	accessJWT := getCloudfareToken(r)
	result, err := c.ValidateCFToken(r.Context(), accessJWT, scope)
	if err != nil {
		return result, &models.CustomError{
			Code:     models.NotFound,
			Messages: []string{fmt.Sprintf("Resource not found: %s", err)},
		}
	}

	if !result.IsValid {
		return result, &models.CustomError{
			Code:     models.Forbidden,
			Messages: []string{"Forbidded"},
		}
	}

	return result, nil
}

func getCloudfareToken(r *http.Request) string {
	accessJWT := r.Header.Get("Cf-Access-Jwt-Assertion")
	if accessJWT != "" {
		return accessJWT
	}

	cfAuthorization, err := r.Cookie("CF_Authorization")
	if err == nil {
		accessJWT = cfAuthorization.Value
		if accessJWT != "" {
			return accessJWT
		}
	}

	return ""
}

func readCFConfig() *models.CFConfiguration {
	// Get Configuration FROM ENV
	cfg := &models.CFConfiguration{}

	// Anything that is not AUTH_ENABLED="true" will disable Authentication
	if authEnabled, ok := os.LookupEnv("AUTH_ENABLED"); ok && authEnabled != "true" {
		if cfEnabled, ok := os.LookupEnv("CF_ENABLED"); ok {
			cfg.Enabled = cfEnabled == "true"
		}
	} else {
		cfg.Enabled = false
	}

	if cfg.Enabled {
		if v, ok := os.LookupEnv("CLOUDFLARE_AUDIENCE"); ok {
			cfg.AUD = v
		} else {
			logs.Fatal("[main] missing CLOUDFLARE_AUDIENCE environment")
		}

		if v, ok := os.LookupEnv("SERVICE_PUBLIC_DOMAIN"); ok {
			cfg.ServicePublicDomain = v
		} else {
			logs.Fatal("[main] missing SERVICE_PUBLIC_DOMAIN environmet")
		}
	} else {
		logs.Info("[main] Cloudflare disabled")
	}

	return cfg
}
