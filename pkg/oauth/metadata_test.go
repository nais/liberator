package oauth_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nais/liberator/pkg/oauth"
)

const (
	metadataResponseTemplate = `{
  "issuer": "https://%[1]s",
  "jwks_uri": "https://%[1]s/jwks",
  "token_endpoint": "https://%[1]s/token"
}`
)

func TestMetadataFetcher_FetchOpenID(t *testing.T) {
	server := server()
	defer server.Close()

	client := server.Client()
	wellKnownUrl := server.URL + oauth.WellKnownOpenIDPath

	metadata, err := oauth.Metadata(wellKnownUrl).
		WithHttpClient(client).
		OpenID(context.Background())

	assert.NoError(t, err)

	assert.NotEmpty(t, metadata.Issuer)
	assert.Equal(t, server.URL, metadata.Issuer)

	assert.NotEmpty(t, metadata.TokenEndpoint)
	assert.Equal(t, server.URL+"/token", metadata.TokenEndpoint)

	assert.NotEmpty(t, metadata.JwksURI)
	assert.Equal(t, server.URL+"/jwks", metadata.JwksURI)
}

func TestMetadataFetcher_FetchOAuth(t *testing.T) {
	server := server()
	defer server.Close()

	client := server.Client()
	wellKnownUrl := server.URL + oauth.WellKnownOAuthPath

	metadata, err := oauth.Metadata(wellKnownUrl).
		WithHttpClient(client).
		OAuth(context.Background())

	assert.NoError(t, err)

	assert.NotEmpty(t, metadata.Issuer)
	assert.Equal(t, server.URL, metadata.Issuer)

	assert.NotEmpty(t, metadata.TokenEndpoint)
	assert.Equal(t, server.URL+"/token", metadata.TokenEndpoint)

	assert.NotEmpty(t, metadata.JwksURI)
	assert.Equal(t, server.URL+"/jwks", metadata.JwksURI)
}

func TestMetadataOpenID_Validate(t *testing.T) {
	type test struct {
		name         string
		wellKnownUrl string
		issuer       string
		wantErr      bool
	}

	tests := []test{
		{
			name:         "valid configuration",
			wellKnownUrl: "https://test.example.com/.well-known/openid-configuration",
			issuer:       "https://test.example.com",
		},
		{
			name:         "valid configuration, issuer with trailing slash",
			wellKnownUrl: "https://test.example.com/.well-known/openid-configuration",
			issuer:       "https://test.example.com/",
		},
		{
			name:         "valid configuration with path",
			wellKnownUrl: "https://test.example.com/some-issuer/.well-known/openid-configuration",
			issuer:       "https://test.example.com/some-issuer",
		},
		{
			name:         "valid configuration with path and trailing slash",
			wellKnownUrl: "https://test.example.com/some-issuer/.well-known/openid-configuration",
			issuer:       "https://test.example.com/some-issuer/",
		},
		{
			name:         "invalid configuration",
			wellKnownUrl: "https://test.example.com/.well-known/openid-configuration",
			issuer:       "https://test.example.com/some-issuer",
			wantErr:      true,
		},
		{
			name:         "invalid configuration, issuer with trailing slash",
			wellKnownUrl: "https://test.example.com/.well-known/openid-configuration",
			issuer:       "https://test.example.com/some-issuer/",
			wantErr:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			metadata := oauth.MetadataOpenID{
				MetadataCommon: oauth.MetadataCommon{
					Issuer: test.issuer,
				},
			}

			err := oauth.Metadata(test.wellKnownUrl).
				Validate(metadata.MetadataCommon, oauth.WellKnownOpenIDPath)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMetadataOAuth_Validate(t *testing.T) {
	type test struct {
		name         string
		wellKnownUrl string
		issuer       string
		wantErr      bool
	}

	tests := []test{
		{
			name:         "valid configuration",
			wellKnownUrl: "https://test.example.com/.well-known/oauth-authorization-server",
			issuer:       "https://test.example.com",
		},
		{
			name:         "valid configuration, issuer with trailing slash",
			wellKnownUrl: "https://test.example.com/.well-known/oauth-authorization-server",
			issuer:       "https://test.example.com/",
		},
		{
			name:         "valid configuration with path",
			wellKnownUrl: "https://test.example.com/some-issuer/.well-known/oauth-authorization-server",
			issuer:       "https://test.example.com/some-issuer",
		},
		{
			name:         "valid configuration with path and trailing slash",
			wellKnownUrl: "https://test.example.com/some-issuer/.well-known/oauth-authorization-server",
			issuer:       "https://test.example.com/some-issuer/",
		},
		{
			name:         "invalid configuration",
			wellKnownUrl: "https://test.example.com/.well-known/oauth-authorization-server",
			issuer:       "https://test.example.com/some-issuer",
			wantErr:      true,
		},
		{
			name:         "invalid configuration, issuer with trailing slash",
			wellKnownUrl: "https://test.example.com/.well-known/oauth-authorization-server",
			issuer:       "https://test.example.com/some-issuer/",
			wantErr:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			metadata := oauth.MetadataOAuth{
				MetadataCommon: oauth.MetadataCommon{
					Issuer: test.issuer,
				},
			}
			err := oauth.Metadata(test.wellKnownUrl).
				Validate(metadata.MetadataCommon, oauth.WellKnownOAuthPath)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func server() *httptest.Server {
	var handler http.HandlerFunc
	handler = func(w http.ResponseWriter, r *http.Request) {
		response := fmt.Sprintf(metadataResponseTemplate, r.Host)
		_, _ = w.Write([]byte(response))
	}

	return httptest.NewTLSServer(handler)
}
