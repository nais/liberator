package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

const (
	WellKnownOAuthPath  = "/.well-known/oauth-authorization-server"
	WellKnownOpenIDPath = "/.well-known/openid-configuration"

	validationErrorTemplate = "well-known url does not match expected url derived from issuer: expected '%s', got '%s'"
)

type MetadataCommon struct {
	Issuer        string `json:"issuer"`
	JwksURI       string `json:"jwks_uri"`
	TokenEndpoint string `json:"token_endpoint"`
}

type MetadataOpenID struct {
	MetadataCommon
}

type MetadataOAuth struct {
	MetadataCommon
}

type metadata struct {
	client       *http.Client
	wellKnownUrl string
}

func Metadata(wellKnownUrl string) *metadata {
	return &metadata{
		client:       http.DefaultClient,
		wellKnownUrl: wellKnownUrl,
	}
}

func (in *metadata) WithHttpClient(client *http.Client) *metadata {
	in.client = client
	return in
}

// OpenID attempts to fetch, unmarshal and Validate an OpenID Connect provider metadata document.
func (in *metadata) OpenID(ctx context.Context) (*MetadataOpenID, error) {
	response, err := in.fetch(ctx)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	metadata, err := in.decodeOpenID(response)
	if err != nil {
		return nil, err
	}

	err = in.Validate(metadata.MetadataCommon, WellKnownOpenIDPath)
	if err != nil {
		return nil, fmt.Errorf("validating openid metadata: %w", err)
	}

	return metadata, nil
}

// OAuth attempts to fetch, unmarshal and Validate an OAuth Authorization Server metadata document.
func (in *metadata) OAuth(ctx context.Context) (*MetadataOAuth, error) {
	response, err := in.fetch(ctx)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	metadata, err := in.decodeOAuth(response)
	if err != nil {
		return nil, err
	}

	err = in.Validate(metadata.MetadataCommon, WellKnownOAuthPath)
	if err != nil {
		return nil, fmt.Errorf("validating oauth metadata: %w", err)
	}

	return metadata, nil
}

// Validate checks that a given issuer from the metadata document is equal to the well-known URL used to retrieve said document.
// See:
// - https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderConfigurationValidation
// - https://datatracker.ietf.org/doc/html/rfc8414#section-3.3
func (in *metadata) Validate(metadata MetadataCommon, expectedPath string) error {
	issuerUrl, err := url.Parse(metadata.Issuer)
	if err != nil {
		return fmt.Errorf("metadata contains an invalid issuer url: %w", err)
	}

	issuerUrl.Path = path.Join(issuerUrl.Path, expectedPath)
	expected := issuerUrl.String()

	if in.wellKnownUrl == expected {
		return nil
	}

	return fmt.Errorf(validationErrorTemplate, expected, in.wellKnownUrl)
}

func (in *metadata) fetch(ctx context.Context) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, in.wellKnownUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	response, err := in.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("fetching metadata from server: %w", err)
	}

	return response, nil
}

func (in *metadata) decodeOpenID(response *http.Response) (*MetadataOpenID, error) {
	var metadata MetadataOpenID
	if err := json.NewDecoder(response.Body).Decode(&metadata); err != nil {
		return nil, fmt.Errorf("decoding metadata: %w", err)
	}

	return &metadata, nil
}

func (in *metadata) decodeOAuth(response *http.Response) (*MetadataOAuth, error) {
	var metadata MetadataOAuth
	if err := json.NewDecoder(response.Body).Decode(&metadata); err != nil {
		return nil, fmt.Errorf("decoding metadata: %w", err)
	}

	return &metadata, nil
}
