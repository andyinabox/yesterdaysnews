package domain

import "context"

type Builder interface {
	Build(ctx context.Context, playlistIds ...string) error
	BuildVideoClips(context.Context) error
	BuildModel(context.Context) error
	BuildManifest(context.Context) error
	Promote(context.Context) error
}
