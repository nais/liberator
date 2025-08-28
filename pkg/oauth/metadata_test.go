package oauth_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nais/liberator/pkg/oauth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	handlerSuccess = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := fmt.Sprintf(`{
  "issuer": "https://%[1]s",
  "jwks_uri": "https://%[1]s/jwks",
  "token_endpoint": "https://%[1]s/token"
}`, r.Host)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	})

	handlerServerError = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		_, _ = w.Write([]byte(`invalid`))
	})
)

func TestNewMetadataOpenID(t *testing.T) {
	server := httptest.NewTLSServer(handlerSuccess)
	defer server.Close()

	client := server.Client()
	wellKnownURL := server.URL + oauth.WellKnownOpenIDSuffix

	metadata, err := oauth.NewMetadataOpenID(t.Context(), wellKnownURL, oauth.WithHTTPClient(client))
	require.NoError(t, err)

	assert.NotEmpty(t, metadata.Issuer)
	assert.Equal(t, server.URL, metadata.Issuer)

	assert.NotEmpty(t, metadata.TokenEndpoint)
	assert.Equal(t, server.URL+"/token", metadata.TokenEndpoint)

	assert.NotEmpty(t, metadata.JwksURI)
	assert.Equal(t, server.URL+"/jwks", metadata.JwksURI)
}

func TestNewMetadataOpenID_InvalidResponse(t *testing.T) {
	server := httptest.NewTLSServer(handlerServerError)
	defer server.Close()

	client := server.Client()
	wellKnownUrl := server.URL + oauth.WellKnownOpenIDSuffix

	_, err := oauth.NewMetadataOpenID(t.Context(), wellKnownUrl, oauth.WithHTTPClient(client))
	assert.Error(t, err)
	assert.ErrorIs(t, err, oauth.ErrInvalidResponse)
}

func TestNewMetadataOAuth(t *testing.T) {
	server := httptest.NewTLSServer(handlerSuccess)
	defer server.Close()

	client := server.Client()
	wellKnownURL := server.URL + oauth.WellKnownOAuthSuffix

	metadata, err := oauth.NewMetadataOAuth(t.Context(), wellKnownURL, oauth.WithHTTPClient(client))
	require.NoError(t, err)

	assert.NotEmpty(t, metadata.Issuer)
	assert.Equal(t, server.URL, metadata.Issuer)

	assert.NotEmpty(t, metadata.TokenEndpoint)
	assert.Equal(t, server.URL+"/token", metadata.TokenEndpoint)

	assert.NotEmpty(t, metadata.JwksURI)
	assert.Equal(t, server.URL+"/jwks", metadata.JwksURI)
}

func TestNewMetadataOAuth_InvalidResponse(t *testing.T) {
	server := httptest.NewTLSServer(handlerServerError)
	defer server.Close()

	client := server.Client()
	wellKnownUrl := server.URL + oauth.WellKnownOAuthSuffix

	_, err := oauth.NewMetadataOAuth(t.Context(), wellKnownUrl, oauth.WithHTTPClient(client))
	assert.Error(t, err)
	assert.ErrorIs(t, err, oauth.ErrInvalidResponse)
}

func TestMetadataOpenID_Validate(t *testing.T) {
	type test struct {
		name         string
		wellKnownURL string
		metadata     oauth.MetadataOpenID
		wantErr      bool
	}

	tests := []test{
		{
			name:         "valid configuration",
			wellKnownURL: "https://test.example.com/.well-known/openid-configuration",
			metadata: oauth.MetadataOpenID{
				Issuer:        "https://test.example.com",
				TokenEndpoint: "https://test.example.com/token",
				JwksURI:       "https://test.example.com/jwks",
			},
		},
		{
			name:         "valid configuration, issuer with trailing slash",
			wellKnownURL: "https://test.example.com/.well-known/openid-configuration",
			metadata: oauth.MetadataOpenID{
				Issuer:        "https://test.example.com/",
				TokenEndpoint: "https://test.example.com/token",
				JwksURI:       "https://test.example.com/jwks",
			},
		},
		{
			name:         "valid configuration with path",
			wellKnownURL: "https://test.example.com/some-issuer/.well-known/openid-configuration",
			metadata: oauth.MetadataOpenID{
				Issuer:        "https://test.example.com/some-issuer",
				TokenEndpoint: "https://test.example.com/token",
				JwksURI:       "https://test.example.com/jwks",
			},
		},
		{
			name:         "valid configuration with path and trailing slash",
			wellKnownURL: "https://test.example.com/some-issuer/.well-known/openid-configuration",
			metadata: oauth.MetadataOpenID{
				Issuer:        "https://test.example.com/some-issuer/",
				TokenEndpoint: "https://test.example.com/token",
				JwksURI:       "https://test.example.com/jwks",
			},
		},
		{
			name:         "invalid configuration",
			wellKnownURL: "https://test.example.com/.well-known/openid-configuration",
			metadata: oauth.MetadataOpenID{
				Issuer:        "https://test.example.com/some-issuer",
				TokenEndpoint: "https://test.example.com/token",
				JwksURI:       "https://test.example.com/jwks",
			},
			wantErr: true,
		},
		{
			name:         "invalid configuration, issuer with trailing slash",
			wellKnownURL: "https://test.example.com/.well-known/openid-configuration",
			metadata: oauth.MetadataOpenID{
				Issuer:        "https://test.example.com/some-issuer/",
				TokenEndpoint: "https://test.example.com/token",
				JwksURI:       "https://test.example.com/jwks",
			},
			wantErr: true,
		},
		{
			name:         "invalid configuration, issuer with trailing slash",
			wellKnownURL: "https://test.example.com/.well-known/openid-configuration",
			metadata: oauth.MetadataOpenID{
				Issuer:        "https://test.example.com/some-issuer/",
				TokenEndpoint: "https://test.example.com/token",
				JwksURI:       "https://test.example.com/jwks",
			},
			wantErr: true,
		},
		{
			name:         "invalid configuration, missing issuer",
			wellKnownURL: "https://test.example.com/.well-known/openid-configuration",
			metadata: oauth.MetadataOpenID{
				Issuer:        "",
				TokenEndpoint: "https://test.example.com/token",
				JwksURI:       "https://test.example.com/jwks",
			},
			wantErr: true,
		},
		{
			name:         "invalid configuration, missing token_endpoint",
			wellKnownURL: "https://test.example.com/.well-known/openid-configuration",
			metadata: oauth.MetadataOpenID{
				Issuer:        "https://test.example.com",
				TokenEndpoint: "",
				JwksURI:       "https://test.example.com/jwks",
			},
			wantErr: true,
		},
		{
			name:         "invalid configuration, missing jwks_uri",
			wellKnownURL: "https://test.example.com/.well-known/openid-configuration",
			metadata: oauth.MetadataOpenID{
				Issuer:        "https://test.example.com",
				TokenEndpoint: "https://test.example.com/token",
				JwksURI:       "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.metadata.Validate(tt.wellKnownURL)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, oauth.ErrInvalidMetadata)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMetadataOAuth_Validate(t *testing.T) {
	type test struct {
		name         string
		wellKnownURL string
		metadata     oauth.MetadataOAuth
		wantErr      bool
	}

	tests := []test{
		{
			name:         "valid configuration",
			wellKnownURL: "https://test.example.com/.well-known/oauth-authorization-server",
			metadata: oauth.MetadataOAuth{
				Issuer:        "https://test.example.com",
				TokenEndpoint: "https://test.example.com/token",
				JwksURI:       "https://test.example.com/jwks",
			},
		},
		{
			name:         "valid configuration, issuer with trailing slash",
			wellKnownURL: "https://test.example.com/.well-known/oauth-authorization-server",
			metadata: oauth.MetadataOAuth{
				Issuer:        "https://test.example.com/",
				TokenEndpoint: "https://test.example.com/token",
				JwksURI:       "https://test.example.com/jwks",
			},
		},
		{
			name:         "valid configuration with path",
			wellKnownURL: "https://test.example.com/some-issuer/.well-known/oauth-authorization-server",
			metadata: oauth.MetadataOAuth{
				Issuer:        "https://test.example.com/some-issuer",
				TokenEndpoint: "https://test.example.com/token",
				JwksURI:       "https://test.example.com/jwks",
			},
		},
		{
			name:         "valid configuration with path and trailing slash",
			wellKnownURL: "https://test.example.com/some-issuer/.well-known/oauth-authorization-server",
			metadata: oauth.MetadataOAuth{
				Issuer:        "https://test.example.com/some-issuer/",
				TokenEndpoint: "https://test.example.com/token",
				JwksURI:       "https://test.example.com/jwks",
			},
		},
		{
			name:         "invalid configuration",
			wellKnownURL: "https://test.example.com/.well-known/oauth-authorization-server",
			metadata: oauth.MetadataOAuth{
				Issuer:        "https://test.example.com/some-issuer",
				TokenEndpoint: "https://test.example.com/token",
				JwksURI:       "https://test.example.com/jwks",
			},
			wantErr: true,
		},
		{
			name:         "invalid configuration, issuer with trailing slash",
			wellKnownURL: "https://test.example.com/.well-known/oauth-authorization-server",
			metadata: oauth.MetadataOAuth{
				Issuer:        "https://test.example.com/some-issuer/",
				TokenEndpoint: "https://test.example.com/token",
				JwksURI:       "https://test.example.com/jwks",
			},
			wantErr: true,
		},
		{
			name:         "invalid configuration, issuer with trailing slash",
			wellKnownURL: "https://test.example.com/.well-known/oauth-authorization-server",
			metadata: oauth.MetadataOAuth{
				Issuer:        "https://test.example.com/some-issuer/",
				TokenEndpoint: "https://test.example.com/token",
				JwksURI:       "https://test.example.com/jwks",
			},
			wantErr: true,
		},
		{
			name:         "invalid configuration, missing issuer",
			wellKnownURL: "https://test.example.com/.well-known/oauth-authorization-server",
			metadata: oauth.MetadataOAuth{
				Issuer:        "",
				TokenEndpoint: "https://test.example.com/token",
				JwksURI:       "https://test.example.com/jwks",
			},
			wantErr: true,
		},
		{
			name:         "invalid configuration, missing token_endpoint",
			wellKnownURL: "https://test.example.com/.well-known/oauth-authorization-server",
			metadata: oauth.MetadataOAuth{
				Issuer:        "https://test.example.com",
				TokenEndpoint: "",
				JwksURI:       "https://test.example.com/jwks",
			},
			wantErr: true,
		},
		{
			name:         "invalid configuration, missing jwks_uri",
			wellKnownURL: "https://test.example.com/.well-known/oauth-authorization-server",
			metadata: oauth.MetadataOAuth{
				Issuer:        "https://test.example.com",
				TokenEndpoint: "https://test.example.com/token",
				JwksURI:       "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.metadata.Validate(tt.wellKnownURL)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, oauth.ErrInvalidMetadata)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
