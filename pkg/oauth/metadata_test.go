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

func TestMetadataOpenID_WellKnownURL(t *testing.T) {
	type test struct {
		name     string
		metadata oauth.MetadataOpenID
		wantURL  string
		wantErr  bool
	}

	tests := []test{
		{
			name: "valid issuer",
			metadata: oauth.MetadataOpenID{
				Issuer: "https://test.example.com",
			},
			wantURL: "https://test.example.com/.well-known/openid-configuration",
		},
		{
			name: "valid issuer with trailing slash",
			metadata: oauth.MetadataOpenID{
				Issuer: "https://test.example.com/",
			},
			wantURL: "https://test.example.com/.well-known/openid-configuration",
		},
		{
			name: "valid issuer with path",
			metadata: oauth.MetadataOpenID{
				Issuer: "https://test.example.com/some-issuer",
			},
			wantURL: "https://test.example.com/some-issuer/.well-known/openid-configuration",
		},
		{
			name: "valid issuer with path and trailing slash",
			metadata: oauth.MetadataOpenID{
				Issuer: "https://test.example.com/some-issuer/",
			},
			wantURL: "https://test.example.com/some-issuer/.well-known/openid-configuration",
		},
		{
			name: "invalid issuer, missing scheme",
			metadata: oauth.MetadataOpenID{
				Issuer: "test.example.com",
			},
			wantErr: true,
		},
		{
			name: "invalid issuer, missing scheme with path",
			metadata: oauth.MetadataOpenID{
				Issuer: "test.example.com/some-issuer",
			},
			wantErr: true,
		},
		{
			name: "invalid issuer, empty string",
			metadata: oauth.MetadataOpenID{
				Issuer: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL, err := tt.metadata.WellKnownURL()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantURL, gotURL)
			}
		})
	}
}

func TestMetadataOAuth_WellKnownURL(t *testing.T) {
	type test struct {
		name     string
		metadata oauth.MetadataOAuth
		wantURL  string
		wantErr  bool
	}

	tests := []test{
		{
			name: "valid issuer",
			metadata: oauth.MetadataOAuth{
				Issuer: "https://test.example.com",
			},
			wantURL: "https://test.example.com/.well-known/oauth-authorization-server",
		},
		{
			name: "valid issuer with trailing slash",
			metadata: oauth.MetadataOAuth{
				Issuer: "https://test.example.com/",
			},
			wantURL: "https://test.example.com/.well-known/oauth-authorization-server",
		},
		{
			name: "valid issuer with path",
			metadata: oauth.MetadataOAuth{
				Issuer: "https://test.example.com/some-issuer",
			},
			wantURL: "https://test.example.com/some-issuer/.well-known/oauth-authorization-server",
		},
		{
			name: "valid issuer with path and trailing slash",
			metadata: oauth.MetadataOAuth{
				Issuer: "https://test.example.com/some-issuer/",
			},
			wantURL: "https://test.example.com/some-issuer/.well-known/oauth-authorization-server",
		},
		{
			name: "invalid issuer, missing scheme",
			metadata: oauth.MetadataOAuth{
				Issuer: "test.example.com",
			},
			wantErr: true,
		},
		{
			name: "invalid issuer, missing scheme with path",
			metadata: oauth.MetadataOAuth{
				Issuer: "test.example.com/some-issuer",
			},
			wantErr: true,
		},
		{
			name: "invalid issuer, empty string",
			metadata: oauth.MetadataOAuth{
				Issuer: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL, err := tt.metadata.WellKnownURL()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantURL, gotURL)
			}
		})
	}
}

func TestMakeWellKnownURL(t *testing.T) {
	type test struct {
		name         string
		issuer       string
		wantURL      string
		wantErr      bool
		wellKnownURL string
	}

	tests := []test{
		{
			name:    "valid issuer",
			issuer:  "https://test.example.com",
			wantURL: "https://test.example.com/.well-known/openid-configuration",
		},
		{
			name:    "valid issuer with trailing slash",
			issuer:  "https://test.example.com/",
			wantURL: "https://test.example.com/.well-known/openid-configuration",
		},
		{
			name:    "valid issuer with path",
			issuer:  "https://test.example.com/some-issuer",
			wantURL: "https://test.example.com/some-issuer/.well-known/openid-configuration",
		},
		{
			name:    "valid issuer with path and trailing slash",
			issuer:  "https://test.example.com/some-issuer/",
			wantURL: "https://test.example.com/some-issuer/.well-known/openid-configuration",
		},
		{
			name:    "invalid issuer, missing scheme",
			issuer:  "test.example.com",
			wantErr: true,
		},
		{
			name:    "invalid issuer, missing scheme with path",
			issuer:  "test.example.com/some-issuer",
			wantErr: true,
		},
		{
			name:    "invalid issuer, empty string",
			issuer:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL, err := oauth.MakeWellKnownURL(tt.issuer, oauth.WellKnownOpenIDSuffix)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, gotURL)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantURL, gotURL)
			}
		})
	}
}
