package fetcher

import (
	"context"
	"errors"

	"github.com/weaveworks/weave-gitops/core/cache"
	mngr "github.com/weaveworks/weave-gitops/core/clustersmngr"
	"k8s.io/client-go/rest"
)

type multiClusterFetcher struct {
	cfg   *rest.Config
	cache cache.Container
}

func NewMultiClusterFetcher(config *rest.Config, cache cache.Container) (mngr.ClusterFetcher, error) {
	return multiClusterFetcher{
		cfg:   config,
		cache: cache,
	}, nil
}

func (f multiClusterFetcher) Fetch(ctx context.Context) ([]mngr.Cluster, error) {
	clusters := []mngr.Cluster{f.self()}

	res, err := f.cache.List(cache.ClusterStorage)
	if err != nil {
		return nil, err
	}

	c, ok := res.([]mngr.Cluster)
	if !ok {
		return nil, errors.New("could not convert objects to []mngr.Cluster")
	}

	return append(clusters, c...), nil
}

func (f *multiClusterFetcher) self() mngr.Cluster {
	return mngr.Cluster{
		Name:        mngr.DefaultCluster,
		Server:      f.cfg.Host,
		BearerToken: f.cfg.BearerToken,
		TLSConfig:   f.cfg.TLSClientConfig,
	}
}
