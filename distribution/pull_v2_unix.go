// +build !windows

package distribution

import (
	"github.com/containers/image/signature"
	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
)

func (ld *v2LayerDescriptor) open(ctx context.Context) (distribution.ReadSeekCloser, error) {
	blobs := ld.repo.Blobs(ctx)
	return blobs.Open(ctx, ld.digest)
}

func configurePolicyContext() (*signature.PolicyContext, error) {
	defaultPolicy, err := signature.DefaultPolicy(nil)
	if err != nil {
		return nil, err
	}
	pc, err := signature.NewPolicyContext(defaultPolicy)
	if err != nil {
		return nil, err
	}
	return pc, nil
}
