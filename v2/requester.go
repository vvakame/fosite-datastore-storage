package fdsstorage

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/ory/fosite"
	"go.mercari.io/datastore"
)

var _ fosite.Requester = (*DefaultRequester)(nil)
var _ fosite.AccessRequester = (*DefaultRequester)(nil)
var _ fosite.AuthorizeRequester = (*DefaultRequester)(nil)

var _ datastore.PropertyLoadSaver = (*DefaultRequester)(nil)

var _ ActiveStateModifier = (*DefaultRequester)(nil)
var _ ClientLoader = (*DefaultRequester)(nil)
var _ SessionRestorer = (*DefaultRequester)(nil)

// ActiveStateModifier provides an action to enable and disable for fosite.Requester.
type ActiveStateModifier interface {
	IsActive() bool
	SetActive(active bool)
}

// ClientLoader provides an action to set Client or get ClientID for fosite.Requester.
type ClientLoader interface {
	GetClientID() string
	SetClient(client fosite.Client)
}

// SessionRestorer provides an action to restore Session from saved Requester for fosite.Requester.
type SessionRestorer interface {
	RestoreSession(ctx context.Context, session fosite.Session) error
}

// DefaultRequester implements fosite.Request, fosite.AccessRequest and fosite.AuthorizeRequest.
type DefaultRequester struct {
	// for fosite.Request
	ID                string         `` // Not Datastore Key
	RequestedAt       time.Time      ``
	ClientID          string         `json:"-"`
	Client            fosite.Client  `datastore:"-"`
	RequestedScope    []string       ``
	GrantedScope      []string       ``
	EncodedForm       string         `json:"-"`
	Form              url.Values     `datastore:"-"`
	SessionJSON       string         `json:"-"`
	Session           fosite.Session `datastore:"-"`
	RequestedAudience []string       ``
	GrantedAudience   []string       ``
	// for fosite.AccessRequest
	GrantTypes       []string ``
	HandledGrantType []string ``
	// for fosite.AuthorizeRequest
	ResponseTypes        []string ``
	RedirectURI          string   ``
	State                string   ``
	HandledResponseTypes []string ``
	// others...
	Active    bool      ``
	UpdatedAt time.Time ``
	CreatedAt time.Time ``
}

// Load loads all of the provided properties into *DefaultRequester.
func (r *DefaultRequester) Load(ctx context.Context, ps []datastore.Property) error {
	err := datastore.LoadStruct(ctx, r, ps)
	if err != nil {
		return err
	}

	r.Form, err = url.ParseQuery(r.EncodedForm)
	if err != nil {
		return err
	}

	return nil
}

// Save saves all of *DefaultRequester's properties as a slice of Properties.
func (r *DefaultRequester) Save(ctx context.Context) ([]datastore.Property, error) {
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now()
	}
	r.UpdatedAt = time.Now()

	if r.Client != nil {
		r.ClientID = r.Client.GetID()
	} else {
		r.ClientID = ""
	}

	if r.Session != nil {
		b, err := json.Marshal(r.Session)
		if err != nil {
			return nil, err
		}
		r.SessionJSON = string(b)
	}

	r.EncodedForm = r.Form.Encode()

	return datastore.SaveStruct(ctx, r)
}

// IsActive returns this request is in active state.
func (r *DefaultRequester) IsActive() bool {
	return r.Active
}

// SetActive state to specified value.
func (r *DefaultRequester) SetActive(active bool) {
	r.Active = active
}

// GetClientID returns client ID.
func (r *DefaultRequester) GetClientID() string {
	return r.ClientID
}

// SetClient to specified value.
func (r *DefaultRequester) SetClient(client fosite.Client) {
	r.Client = client
}

// RestoreSession restores the passed Session to its original state and retains it as its own session.
func (r *DefaultRequester) RestoreSession(ctx context.Context, session fosite.Session) error {
	if r.SessionJSON != "" {
		err := json.Unmarshal([]byte(r.SessionJSON), session)
		if err != nil {
			return err
		}
		r.Session = session
	} else {
		r.Session = nil
	}

	return nil
}

// SetID sets the unique identifier.
func (r *DefaultRequester) SetID(id string) {
	r.ID = id
}

// GetID returns a unique identifier.
func (r *DefaultRequester) GetID() string {
	return r.ID
}

// GetRequestedAt returns the time the request was created.
func (r *DefaultRequester) GetRequestedAt() time.Time {
	return r.RequestedAt
}

// GetClient returns the requests client.
func (r *DefaultRequester) GetClient() fosite.Client {
	return r.Client
}

// GetRequestedScopes returns the request's scopes.
func (r *DefaultRequester) GetRequestedScopes() fosite.Arguments {
	return r.RequestedScope
}

// GetRequestedAudience returns the requested audiences for this request.
func (r *DefaultRequester) GetRequestedAudience() fosite.Arguments {
	return r.RequestedAudience
}

// SetRequestedScopes sets the request's scopes.
func (r *DefaultRequester) SetRequestedScopes(scopes fosite.Arguments) {
	r.RequestedScope = nil
	for _, scope := range scopes {
		r.AppendRequestedScope(scope)
	}
}

// SetRequestedAudience sets the requested audience.
func (r *DefaultRequester) SetRequestedAudience(audience fosite.Arguments) {
	r.RequestedAudience = nil
	for _, a := range audience {
		r.AppendRequestedAudience(a)
	}
}

// AppendRequestedScope appends a scope to the request.
func (r *DefaultRequester) AppendRequestedScope(scope string) {
	for _, old := range r.RequestedScope {
		if scope == old {
			return
		}
	}
	r.RequestedScope = append(r.RequestedScope, scope)
}

// AppendRequestedAudience appends a audience to the request.
func (r *DefaultRequester) AppendRequestedAudience(s string) {
	for _, old := range r.RequestedAudience {
		if s == old {
			return
		}
	}
	r.RequestedAudience = append(r.RequestedAudience, s)
}

// GetGrantedScopes returns all granted scopes.
func (r *DefaultRequester) GetGrantedScopes() fosite.Arguments {
	return r.GrantedScope
}

// GetGrantedAudience returns all granted scopes.
func (r *DefaultRequester) GetGrantedAudience() fosite.Arguments {
	return r.GrantedAudience
}

// GrantScope marks a request's scope as granted.
func (r *DefaultRequester) GrantScope(scope string) {
	for _, old := range r.GrantedScope {
		if scope == old {
			return
		}
	}
	r.GrantedScope = append(r.GrantedScope, scope)
}

// GrantAudience marks a request's audience as granted.
func (r *DefaultRequester) GrantAudience(audience string) {
	for _, old := range r.GrantedAudience {
		if audience == old {
			return
		}
	}
	r.GrantedAudience = append(r.GrantedAudience, audience)
}

// GetSession returns a pointer to the request's session or nil if none is set.
func (r *DefaultRequester) GetSession() (session fosite.Session) {
	return r.Session
}

// SetSession sets the request's session pointer.
func (r *DefaultRequester) SetSession(session fosite.Session) {
	r.Session = session
}

// GetRequestForm returns the request's form input.
func (r *DefaultRequester) GetRequestForm() url.Values {
	return r.Form
}

// Merge merges the argument into the method receiver.
func (r *DefaultRequester) Merge(requester fosite.Requester) {
	for _, scope := range requester.GetRequestedScopes() {
		r.AppendRequestedScope(scope)
	}
	for _, scope := range requester.GetGrantedScopes() {
		r.GrantScope(scope)
	}

	for _, aud := range requester.GetRequestedAudience() {
		r.AppendRequestedAudience(aud)
	}
	for _, aud := range requester.GetGrantedAudience() {
		r.GrantAudience(aud)
	}

	r.RequestedAt = requester.GetRequestedAt()
	r.Client = requester.GetClient()
	r.Session = requester.GetSession()

	for k, v := range requester.GetRequestForm() {
		r.Form[k] = v
	}
}

// Sanitize returns a sanitized clone of the request which can be used for storage.
func (r *DefaultRequester) Sanitize(allowedParameters []string) fosite.Requester {
	n := &DefaultRequester{}
	allowed := make(map[string]bool)
	for _, v := range allowedParameters {
		allowed[v] = true
	}

	*n = *r
	n.ID = r.GetID()
	n.Form = url.Values{}
	for k := range r.Form {
		if _, ok := allowed[k]; ok {
			n.Form.Add(k, r.Form.Get(k))
		}
	}

	return n
}

// GetGrantTypes returns the requests grant type.
func (r *DefaultRequester) GetGrantTypes() (grantTypes fosite.Arguments) {
	return r.GrantTypes
}

// GetResponseTypes returns the requested response types
func (r *DefaultRequester) GetResponseTypes() (responseTypes fosite.Arguments) {
	return r.ResponseTypes
}

// SetResponseTypeHandled marks a response_type (e.g. token or code) as handled indicating that the response type
// is supported.
func (r *DefaultRequester) SetResponseTypeHandled(responseType string) {
	r.HandledResponseTypes = append(r.HandledResponseTypes, responseType)
}

// DidHandleAllResponseTypes returns if all requested response types have been handled correctly
func (r *DefaultRequester) DidHandleAllResponseTypes() bool {
	for _, rt := range r.ResponseTypes {
		for _, handle := range r.HandledResponseTypes {
			if rt == handle {
				return false
			}
		}
	}

	return len(r.ResponseTypes) > 0

}

// GetRedirectURI returns the requested redirect URI
func (r *DefaultRequester) GetRedirectURI() *url.URL {
	if r.RedirectURI == "" {
		return nil
	}
	redirectURL, err := url.Parse(r.RedirectURI)
	if err != nil {
		return nil
	}
	return redirectURL
}

// IsRedirectURIValid returns false if the redirect is not rfc-conform (i.e. missing client, not on white list,
// or malformed)
func (r *DefaultRequester) IsRedirectURIValid() (isValid bool) {
	if r.GetRedirectURI() == nil {
		return false
	}

	raw := r.GetRedirectURI().String()
	if r.GetClient() == nil {
		return false
	}

	redirectURI, err := fosite.MatchRedirectURIWithClientRedirectURIs(raw, r.GetClient())
	if err != nil {
		return false
	}
	return fosite.IsValidRedirectURI(redirectURI)

}

// GetState returns the request's state.
func (r *DefaultRequester) GetState() (state string) {
	return r.State
}
