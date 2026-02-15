package services

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type FirebaseAuth struct {
	ProjectID string
	keys      map[string]*rsa.PublicKey
	keysMu    sync.RWMutex
	keysExp   time.Time
}

type FirebaseClaims struct {
	jwt.RegisteredClaims
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	UserID        string `json:"user_id"`
}

func NewFirebaseAuth(projectID string) *FirebaseAuth {
	return &FirebaseAuth{ProjectID: projectID}
}

func (f *FirebaseAuth) VerifyIDToken(idToken string) (*FirebaseClaims, error) {
	keys, err := f.getPublicKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to get public keys: %w", err)
	}

	token, err := jwt.ParseWithClaims(idToken, &FirebaseClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("missing kid in token header")
		}
		key, ok := keys[kid]
		if !ok {
			// Refresh keys and retry
			refreshed, err := f.refreshPublicKeys()
			if err != nil {
				return nil, fmt.Errorf("failed to refresh keys: %w", err)
			}
			key, ok = refreshed[kid]
			if !ok {
				return nil, fmt.Errorf("unknown kid: %s", kid)
			}
		}
		return key, nil
	},
		jwt.WithIssuer(fmt.Sprintf("https://securetoken.google.com/%s", f.ProjectID)),
		jwt.WithAudience(f.ProjectID),
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*FirebaseClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	if claims.UserID == "" {
		return nil, fmt.Errorf("missing user_id in token")
	}

	return claims, nil
}

func (f *FirebaseAuth) getPublicKeys() (map[string]*rsa.PublicKey, error) {
	f.keysMu.RLock()
	if f.keys != nil && time.Now().Before(f.keysExp) {
		defer f.keysMu.RUnlock()
		return f.keys, nil
	}
	f.keysMu.RUnlock()

	return f.refreshPublicKeys()
}

func (f *FirebaseAuth) refreshPublicKeys() (map[string]*rsa.PublicKey, error) {
	f.keysMu.Lock()
	defer f.keysMu.Unlock()

	// Double-check after acquiring write lock
	if f.keys != nil && time.Now().Before(f.keysExp) {
		return f.keys, nil
	}

	resp, err := http.Get("https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var certs map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&certs); err != nil {
		return nil, err
	}

	keys := make(map[string]*rsa.PublicKey)
	for kid, certPEM := range certs {
		block, _ := pem.Decode([]byte(certPEM))
		if block == nil {
			continue
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			continue
		}
		rsaKey, ok := cert.PublicKey.(*rsa.PublicKey)
		if !ok {
			continue
		}
		keys[kid] = rsaKey
	}

	f.keys = keys
	f.keysExp = time.Now().Add(1 * time.Hour)

	return f.keys, nil
}
