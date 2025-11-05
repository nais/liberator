package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"path"
)

type WellKnownSuffix = string

const (
	WellKnownOAuthSuffix  WellKnownSuffix = "/.well-known/oauth-authorization-server"
	WellKnownOpenIDSuffix WellKnownSuffix = "/.well-known/openid-configuration"
)

var (
	ErrInvalidResponse = errors.New("invalid response from server")
	ErrInvalidMetadata = errors.New("invalid metadata")
)

type MetadataFetchOption func(m *metadataFetcher)

func WithHTTPClient(client *http.Client) MetadataFetchOption {
	return func(m *metadataFetcher) {
		m.client = client
	}
}

// NewMetadataOpenID attempts to fetch, unmarshal and validate an OpenID Connect provider metadata document.
func NewMetadataOpenID(ctx context.Context, wellKnownURL string, opts ...MetadataFetchOption) (*MetadataOpenID, error) {
	fetcher := &metadataFetcher{
		client: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(fetcher)
	}

	var metadata MetadataOpenID
	if err := fetcher.get(ctx, wellKnownURL, &metadata); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidResponse, err)
	}

	if err := metadata.Validate(wellKnownURL); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// NewMetadataOAuth attempts to fetch, unmarshal and validate an OAuth Authorization Server metadata document.
func NewMetadataOAuth(ctx context.Context, wellKnownURL string, opts ...MetadataFetchOption) (*MetadataOAuth, error) {
	fetcher := &metadataFetcher{
		client: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(fetcher)
	}

	var metadata MetadataOAuth
	if err := fetcher.get(ctx, wellKnownURL, &metadata); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidResponse, err)
	}

	if err := metadata.Validate(wellKnownURL); err != nil {
		return nil, err
	}

	return &metadata, nil
}

type MetadataOpenID struct {
	Issuer        string `json:"issuer"`
	JwksURI       string `json:"jwks_uri"`
	TokenEndpoint string `json:"token_endpoint"`
}

func (m MetadataOpenID) Validate(wellKnownURL string) error {
	if m.JwksURI == "" {
		return fmt.Errorf("%w: metadata does not contain jwks_uri", ErrInvalidMetadata)
	}
	if m.TokenEndpoint == "" {
		return fmt.Errorf("%w: metadata does not contain token_endpoint", ErrInvalidMetadata)
	}
	if err := validateIssuer(m.Issuer, WellKnownOpenIDSuffix, wellKnownURL); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidMetadata, err)
	}

	return nil
}

func (m MetadataOpenID) WellKnownURL() (string, error) {
	return MakeWellKnownURL(m.Issuer, WellKnownOpenIDSuffix)
}

type MetadataOAuth struct {
	Issuer        string `json:"issuer"`
	JwksURI       string `json:"jwks_uri"`
	TokenEndpoint string `json:"token_endpoint"`
}

func (m MetadataOAuth) Validate(wellKnownURL string) error {
	if m.JwksURI == "" {
		return fmt.Errorf("%w: metadata does not contain jwks_uri", ErrInvalidMetadata)
	}
	if m.TokenEndpoint == "" {
		return fmt.Errorf("%w: metadata does not contain token_endpoint", ErrInvalidMetadata)
	}
	if err := validateIssuer(m.Issuer, WellKnownOAuthSuffix, wellKnownURL); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidMetadata, err)
	}

	return nil
}

func (m MetadataOAuth) WellKnownURL() (string, error) {
	return MakeWellKnownURL(m.Issuer, WellKnownOAuthSuffix)
}

type metadataFetcher struct {
	client *http.Client
}

func (f *metadataFetcher) get(ctx context.Context, wellKnownURL string, v any) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, wellKnownURL, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	response, err := f.client.Do(request)
	if err != nil {
		return fmt.Errorf("fetching metadata from server: %w", err)
	}
	defer response.Body.Close()

	// https://datatracker.ietf.org/doc/html/rfc8414#section-3.2
	// > A successful response MUST use the 200 OK HTTP
	//   status code and return a JSON object using the "application/json"
	//   content type
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from server: %d", response.StatusCode)
	}

	contentType, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil {
		return fmt.Errorf("parsing content-type from server: %w", err)
	}
	if contentType != "application/json" {
		return fmt.Errorf("unexpected content-type from server (want: %q, got: %q)", "application/json", response.Header.Get("Content-Type"))
	}
	if err := json.NewDecoder(response.Body).Decode(&v); err != nil {
		return fmt.Errorf("decoding metadata: %w", err)
	}

	return nil
}

// validateIssuer checks that a given issuer from the metadata document is equal to the well-known URL used to retrieve said document.
// See:
// - https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderConfigurationValidation
// - https://datatracker.ietf.org/doc/html/rfc8414#section-3.3
func validateIssuer(issuer string, suffix WellKnownSuffix, wellKnownURL string) error {
	if wellKnownURL == "" {
		return fmt.Errorf("well-known URL is empty")
	}

	expected, err := MakeWellKnownURL(issuer, suffix)
	if err != nil {
		return err
	}

	if wellKnownURL != expected {
		return fmt.Errorf("well-known URL does not match expected URL derived from issuer: expected '%s', got '%s'", expected, wellKnownURL)
	}

	return nil
}

// MakeWellKnownURL constructs the well-known URL by appending the given suffix to the issuer URL.
func MakeWellKnownURL(issuerURL string, suffix WellKnownSuffix) (string, error) {
	if issuerURL == "" {
		return "", fmt.Errorf("issuer is empty")
	}
	u, err := url.Parse(issuerURL)
	if err != nil {
		return "", fmt.Errorf("issuer is not a valid URL: %w", err)
	}
	if u.Scheme == "" {
		return "", fmt.Errorf("issuer URL does not contain a scheme")
	}

	u.Path = path.Join(u.Path, suffix)
	return u.String(), nil
}
