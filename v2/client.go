package fdsstorage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ory/fosite"
	"go.mercari.io/datastore"
	"gopkg.in/square/go-jose.v2"
)

var _ fosite.Client = (*DefaultClient)(nil)
var _ fosite.OpenIDConnectClient = (*DefaultClient)(nil)
var _ datastore.KeyLoader = (*DefaultClient)(nil)
var _ datastore.PropertyLoadSaver = (*DefaultClient)(nil)

// DefaultClient is a simple default implementation of the Client interface for Datastore.
// It's support fosite.Client and fosite.OpenIDConnectClient interface.
type DefaultClient struct {
	// for fosite.Client
	ID            string   `datastore:"-" boom:"id"`
	Secret        []byte   `datastore:",noindex"`
	RedirectURIs  []string ``
	GrantTypes    []string ``
	ResponseTypes []string ``
	Scopes        []string ``
	Audience      []string ``
	Public        bool     ``
	// for fosite.OpenIDConnectClient
	JSONWebKeysURI                string              ``
	JSONWebKeysJSON               string              `json:"-" datastore:",noindex"`
	JSONWebKeys                   *jose.JSONWebKeySet `datastore:"-"`
	TokenEndpointAuthMethod       string              ``
	RequestURIs                   []string            ``
	RequestObjectSigningAlgorithm string              ``
	// others...
	UpdatedAt time.Time ``
	CreatedAt time.Time ``
}

// LoadKey is restore Client ID from Datastore key.
func (cli *DefaultClient) LoadKey(ctx context.Context, key datastore.Key) error {
	cli.ID = key.Name()
	return nil
}

// Load loads all of the provided properties into *DefaultClient.
func (cli *DefaultClient) Load(ctx context.Context, ps []datastore.Property) error {
	err := datastore.LoadStruct(ctx, cli, ps)
	if err != nil {
		return err
	}

	if cli.JSONWebKeysJSON != "" {
		var jwks jose.JSONWebKeySet
		err = json.Unmarshal([]byte(cli.JSONWebKeysJSON), &jwks)
		if err != nil {
			return err
		}
		cli.JSONWebKeys = &jwks
	}

	return nil
}

// Save saves all of *DefaultClient's properties as a slice of Properties.
func (cli *DefaultClient) Save(ctx context.Context) ([]datastore.Property, error) {
	if cli.CreatedAt.IsZero() {
		cli.CreatedAt = time.Now()
	}
	cli.UpdatedAt = time.Now()

	if cli.JSONWebKeys != nil {
		b, err := json.Marshal(cli.JSONWebKeys)
		if err != nil {
			return nil, err
		}
		cli.JSONWebKeysJSON = string(b)
	} else {
		cli.JSONWebKeysJSON = ""
	}

	return datastore.SaveStruct(ctx, cli)
}

// GetID returns the client ID.
func (cli *DefaultClient) GetID() string {
	return cli.ID
}

// GetHashedSecret returns the hashed secret as it is stored in the store.
func (cli *DefaultClient) GetHashedSecret() []byte {
	return cli.Secret
}

// GetRedirectURIs returns the client's allowed redirect URIs.
func (cli *DefaultClient) GetRedirectURIs() []string {
	return cli.RedirectURIs
}

// GetGrantTypes returns the client's allowed grant types.
func (cli *DefaultClient) GetGrantTypes() fosite.Arguments {
	if len(cli.GrantTypes) == 0 {
		return fosite.Arguments{"authorization_code"}
	}
	return cli.GrantTypes
}

// GetResponseTypes returns the client's allowed response types.
// default is the "code".
func (cli *DefaultClient) GetResponseTypes() fosite.Arguments {
	if len(cli.ResponseTypes) == 0 {
		return fosite.Arguments{"code"}
	}
	return cli.ResponseTypes
}

// GetScopes returns the scopes this client is allowed to request.
func (cli *DefaultClient) GetScopes() fosite.Arguments {
	return cli.Scopes
}

// IsPublic returns true, if this client is marked as public.
func (cli *DefaultClient) IsPublic() bool {
	return cli.Public
}

// GetAudience returns the allowed audience(s) for this client.
func (cli *DefaultClient) GetAudience() fosite.Arguments {
	return cli.Audience
}

// GetRequestURIs is an array of request_uri values that are pre-registered by the RP for use at the OP. Servers MAY
// cache the contents of the files referenced by these URIs and not retrieve them at the time they are used in a request.
// OPs can require that request_uri values used be pre-registered with the require_request_uri_registration
// discovery parameter.
func (cli *DefaultClient) GetRequestURIs() []string {
	return cli.RequestURIs
}

// GetJSONWebKeys returns the JSON Web Key Set containing the public keys used by the client to authenticate.
func (cli *DefaultClient) GetJSONWebKeys() *jose.JSONWebKeySet {
	return nil
}

// GetJSONWebKeysURI returns the URL for lookup of JSON Web Key Set containing the
// public keys used by the client to authenticate.
func (cli *DefaultClient) GetJSONWebKeysURI() string {
	return cli.JSONWebKeysURI
}

// GetRequestObjectSigningAlgorithm returns JWS [JWS] alg algorithm [JWA] that MUST be used for signing Request Objects sent to the OP.
// All Request Objects from this Client MUST be rejected, if not signed with this algorithm.
func (cli *DefaultClient) GetRequestObjectSigningAlgorithm() string {
	return cli.RequestObjectSigningAlgorithm
}

// GetTokenEndpointAuthMethod returns requested Client Authentication method for the Token Endpoint. The options are client_secret_post,
// client_secret_basic, client_secret_jwt, private_key_jwt, and none.
func (cli *DefaultClient) GetTokenEndpointAuthMethod() string {
	return cli.TokenEndpointAuthMethod
}

// GetTokenEndpointAuthSigningAlgorithm returns JWS [JWS] alg algorithm [JWA] that MUST be used for signing the JWT [JWT] used to authenticate the
// Client at the Token Endpoint for the private_key_jwt and client_secret_jwt authentication methods.
func (cli *DefaultClient) GetTokenEndpointAuthSigningAlgorithm() string {
	return "RS256"
}
