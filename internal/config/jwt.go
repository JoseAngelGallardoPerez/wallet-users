package config

import (
	"github.com/Confialink/wallet-pkg-env_config"

	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"sort"
	"strings"
)

type JwtConfiguration struct {
	SecretKeyPath string
	PublicKeyPath string
	jwt.SigningMethod
	Secret string
}

// Init initializes environment variables
func (c *JwtConfiguration) Init() error {
	signingMethods := map[string]jwt.SigningMethod{
		"ES256": jwt.SigningMethodES256,
		"ES384": jwt.SigningMethodES384,
		"ES512": jwt.SigningMethodES512,
		"HS256": jwt.SigningMethodHS256,
		"HS384": jwt.SigningMethodHS384,
		"HS512": jwt.SigningMethodHS512,
	}
	availableMethods := make([]string, 0, len(signingMethods))
	for k := range signingMethods {
		availableMethods = append(availableMethods, k)
		sort.Strings(availableMethods)
	}

	requestedSigningMethod := os.Getenv("VELMIE_WALLET_USERS_JWT_SIGNING_METHOD")
	if requestedSigningMethod == "" {
		return errors.New("VELMIE_WALLET_USERS_JWT_SIGNING_METHOD environment variable is required and cannot be empty")
	}
	if signingMethod, ok := signingMethods[requestedSigningMethod]; !ok {
		format := "requested signing method \"%s\" is not supported, set VELMIE_WALLET_USERS_JWT_SIGNING_METHOD environment " +
			"variable with one of the following values: %s"
		return fmt.Errorf(format, requestedSigningMethod, strings.Join(availableMethods, ", "))
	} else {
		c.SigningMethod = signingMethod
	}

	switch c.SigningMethod.(type) {
	case *jwt.SigningMethodHMAC:
		secret := os.Getenv("VELMIE_WALLET_USERS_JWT_SECRET")
		if secret == "" {
			return errors.New("VELMIE_WALLET_USERS_JWT_SECRET environment variable is required and cannot be empty " +
				"if signing method is HMAC")
		}
		c.Secret = secret
	case *jwt.SigningMethodECDSA:
		secretEnv := "VELMIE_WALLET_USERS_SECRET_KEY_PATH"
		publicEnv := "VELMIE_WALLET_USERS_PUBLIC_KEY_PATH"
		c.SecretKeyPath = env_config.Env(secretEnv, "./keys/jwt.pem")
		c.PublicKeyPath = env_config.Env(publicEnv, "./keys/jwt.pub")
		for k, path := range map[string]string{secretEnv: c.SecretKeyPath, publicEnv: c.PublicKeyPath} {
			info, err := os.Stat(path)
			if os.IsNotExist(err) {
				return fmt.Errorf("file is not exist %s, set correct file path with %s environment variable", path, k)
			}
			if os.IsPermission(err) {
				return fmt.Errorf("failed to read %s, permission error (%s)", path, k)
			}
			if info.IsDir() {
				return fmt.Errorf("%s is a dirctory (%s) expected key file", path, k)
			}
		}
	}

	return nil
}
