package openstackclient

import "context"

type IdentityClient interface {
	AuthToken(context.Context) (string, error)
}
