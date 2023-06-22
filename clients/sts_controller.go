package clients

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/deliveryhero/devhub-backend-core/models"
	sts "github.com/deliveryhero/pd-sts-go-sdk"
	"github.com/pedidosya/peya-go/logs"
)

type STSControllerSpec interface {
	Validate(ctx context.Context, accessToken string, scope string) (*models.AuthResult, error)
}

type STSClient struct {
	introspectionService introspectionService
	Config               *models.STSConfiguration
}

// NewClient creates an STS client that contains the introspection service.
func NewSTSClient() STSControllerSpec {
	c := readSTSConfig()
	if c.Enabled {
		tokenService, err := newTokenService(c)
		if err != nil {
			logs.Errorf("Could not initialise STS Client: %v", err.Error())
			panic(err)
		}

		introspectionService := sts.NewIntrospectionService(c.URL, tokenService)

		return &STSClient{
			introspectionService: introspectionService,
			Config:               c,
		}
	}
	return nil
}

// newTokenService creates a new STS token service.
func newTokenService(cfg *models.STSConfiguration) (sts.TokenService, error) {
	var privateKey *rsa.PrivateKey
	privateKeyString, err := base64.StdEncoding.DecodeString(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 string. %v", err.Error())
	}

	block, _ := pem.Decode(privateKeyString)
	if block == nil {
		return nil, errors.New("failed to decode pem for STS PrivateKey")
	}

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		privateKey = key
	} else if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey:
			privateKey = key
		default:
			return nil, errors.New("crypto/tls: found unknown private key type in PKCS#8 wrapping")
		}
	} else {
		return nil, errors.New("crypto/tls: failed to parse private key")
	}
	if err != nil {
		return nil, fmt.Errorf("x509.ParsePKCS1PrivateKey: %v", err)
	}

	return sts.NewTokenService(cfg.URL, cfg.ClientID, cfg.KeyID, privateKey), nil
}

// Validate validates customer token using STS.
func (c *STSClient) Validate(ctx context.Context, accessToken string, scope string) (*models.AuthResult, error) {
	token, err := c.introspectionService.GetMetadata(ctx, accessToken, scope)
	if err != nil {
		logs.Errorf("[sts_controller] failed to get token metadata: %s", err)
		return &models.AuthResult{
			IsValid: false,
			Enabled: true,
		}, err
	}

	if token == nil {
		logs.Errorf("[sts_controller] invalid token")
		return &models.AuthResult{
			IsValid: false,
			Enabled: true,
		}, nil
	}

	identity := fmt.Sprintf("%s:%s", "sts", token.Subject)

	return &models.AuthResult{
		IsValid:  true,
		Enabled:  true,
		Identity: &identity,
	}, nil
}

// introspectionService is an extracted interface from sts sdk IntrospectionService.
type introspectionService interface {
	GetMetadata(ctx context.Context, accessToken string, scopes ...string) (*sts.Token, error)
	IsValidToken(ctx context.Context, accessToken string, scopes ...string) (bool, error)
	WithTimeout(duration time.Duration) sts.IntrospectionService
}

// MIDDLEWARE
func STSAuthCheckMiddleware(sts STSClient, r *http.Request, scope string) (authResult *models.AuthResult, error *models.CustomError) {
	accessToken := getSTSToken(r)
	result, err := sts.validateSTSToken(r.Context(), accessToken, scope)
	if err != nil {
		return result, &models.CustomError{
			Code:     models.NotFound,
			Messages: []string{fmt.Sprintf("Resource not found: %s", err)},
		}
	}

	if !result.IsValid {
		return &models.AuthResult{
				IsValid: false,
				Enabled: true,
			}, &models.CustomError{
				Code:     models.Forbidden,
				Messages: []string{"Forbidded"},
			}
	}

	return result, nil
}

func (c *STSClient) validateSTSToken(ctx context.Context, accessToken string, scope string) (*models.AuthResult, error) {
	if !c.Config.Enabled {
		return &models.AuthResult{
			IsValid: true,
			Enabled: false,
		}, nil
	}

	if accessToken == "" {
		return &models.AuthResult{
			IsValid: false,
			Enabled: true,
		}, errors.New("authorization STS token should not be empty")
	}

	stsToken, err := c.Validate(ctx, accessToken, scope)
	if err != nil {
		return &models.AuthResult{
			IsValid: false,
			Enabled: true,
		}, err
	}

	return stsToken, nil
}

func getSTSToken(r *http.Request) string {
	return r.Header.Get("Authorization")
}

func readSTSConfig() (*models.STSConfiguration) {
	// Get Configuration FROM ENV
	cfg := &models.STSConfiguration{}

	// Anything that is not AUTH_ENABLED="true" will disable Authentication
	if authEnabled, ok := os.LookupEnv("AUTH_ENABLED"); ok && authEnabled != "true" {
		if stsEnabled, ok := os.LookupEnv("STS_ENABLED"); ok {
			cfg.Enabled = stsEnabled == "true"
		}
	} else {
		cfg.Enabled = false
	}

	if cfg.Enabled {
		v, ok := os.LookupEnv("STS_URL")
		if ok {
			cfg.URL = v
		} else {
			logs.Fatal("[main] missing STS_URL environment or sts.url setting")
		}

		v, ok = os.LookupEnv("STS_CLIENT_ID")
		if ok {
			cfg.ClientID = v
		} else {
			logs.Fatal("[main] missing STS_CLIENT_ID environment or sts.client_id setting")
		}

		v, ok = os.LookupEnv("STS_KEY_ID")
		if ok {
			cfg.KeyID = v
		} else {
			logs.Fatal("[main] missing STS_KEY_ID environment or sts.key_id setting")
		}

		v, ok = os.LookupEnv("STS_TIMEOUT")
		if ok {
			i, err := strconv.Atoi(v)
			if err == nil {
				cfg.Timeout = i
			} else {
				logs.Debug("Using default timeout for STS")
				cfg.Timeout = 1500
			}
		} else {
			logs.Fatal("[main] missing STS_TIMEOUT environment or sts.timeout setting")
		}

		v, ok = os.LookupEnv("STS_PRIVATE_KEY")
		if ok {
			cfg.PrivateKey = v
		} else {
			logs.Fatal("[main] missing STS_PRIVATE_KEY environment or sts.private_key setting")
		}
	} else {
		logs.Info("[main] STS disabled")
	}

	return cfg
}