//+build wireinject

package idp

import (
	"context"

	"github.com/google/wire"
	"github.com/ory/fosite"
	"github.com/vvakame/fosite-datastore-storage/example/domains"
)

func InitializeProvider() (fosite.OAuth2Provider, error) {
	wire.Build(ProvideDatastore, ProvideStore, ProvideConfig, ProvideRSAPrivateKey, ProvideStrategy, ProvideOAuth2Provider)
	return nil, nil
}

func InitializeSession(ctx context.Context, user *domains.User) (Session, error) {
	wire.Build(ProvideSession)
	return nil, nil
}
