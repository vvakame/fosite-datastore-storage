package fdsstorage

import (
	"errors"
)

var errUnsupportedRequesterType = errors.New("requester type must be *fosite.Request or *fosite.AccessRequest or *fosite.AuthorizeRequest or datastore.PropertyLoadSaver")
var errRequesterNeedsActiveStateModifier = errors.New("requester is not implement ActiveStateModifier")
var errRequesterNeedsClientLoader = errors.New("requester is not implement ClientLoader")
var errUnsupportedClientType = errors.New("client type must be *fosite.DefaultClient or *fosite.DefaultOpenIDConnectClient or datastore.PropertyLoadSaver")

var errInvalidTxContext = errors.New("context doesn't in tx context")
