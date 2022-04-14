package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-logr/logr"
	gitopsv1alpha1 "github.com/weaveworks/cluster-controller/api/v1alpha1"
	mngr "github.com/weaveworks/weave-gitops/core/clustersmngr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ClusterStorage StorageType = "cluster"
	dataKey                    = "value"
)

func WithClusterCache(c client.Client) WithCacheFunc {
	return func(l logr.Logger, syncChan chan mngr.Cluster) (StorageType, Cache) {
		return ClusterStorage, newClusterStore(l, c, syncChan)
	}
}

type clusterStore struct {
	logger      logr.Logger
	client      client.Client
	clusterChan chan<- mngr.Cluster

	clusters []mngr.Cluster

	lock     sync.Mutex
	cancel   func()
	interval time.Duration
}

func newClusterStore(logger logr.Logger, c client.Client, clusterChan chan<- mngr.Cluster) *clusterStore {
	return &clusterStore{
		logger:      logger.WithName("cluster-cache"),
		client:      c,
		clusterChan: clusterChan,
		clusters:    []mngr.Cluster{},
		lock:        sync.Mutex{},
		cancel:      nil,
		interval:    DefaultPollIntervalSeconds,
	}
}

func (c *clusterStore) List() interface{} {
	return c.clusters
}

func (c *clusterStore) ForceRefresh() {
	c.update(context.Background())
}

func (c *clusterStore) Stop() {
	if c.cancel != nil {
		c.logger.Info("stopping cluster cache")

		c.cancel()
	}

	if c.clusterChan != nil {
		close(c.clusterChan)
	}
}

func (c *clusterStore) Start(ctx context.Context) {
	var newCtx context.Context

	newCtx, c.cancel = context.WithCancel(ctx)

	c.logger.Info("starting cluster cache")

	// Force load clusters on startup
	c.update(newCtx)

	go func() {
		ticker := time.NewTicker(c.interval * time.Second)

		// TODO why is this ticking so much?
		defer ticker.Stop()

		for {
			select {
			case <-newCtx.Done():
				break
			case <-ticker.C:
				continue
			}

			c.update(newCtx)
		}
	}()
}

func (c *clusterStore) update(ctx context.Context) {
	c.lock.Lock()
	defer c.lock.Unlock()

	clusters := []mngr.Cluster{}

	goClusters := &gitopsv1alpha1.GitopsClusterList{}

	if err := c.client.List(ctx, goClusters); err != nil && !meta.IsNoMatchError(err) {
		c.logger.Error(err, "cluster poll error")
	}

	for _, cluster := range goClusters.Items {
		var secretRef string

		if cluster.Spec.SecretRef != nil {
			secretRef = cluster.Spec.SecretRef.Name
		}

		if secretRef == "" && cluster.Spec.CAPIClusterRef != nil {
			secretRef = fmt.Sprintf("%s-kubeconfig", cluster.Spec.CAPIClusterRef.Name)
		}

		if secretRef == "" {
			continue
		}

		key := types.NamespacedName{
			Name:      secretRef,
			Namespace: cluster.Namespace,
		}

		var secret v1.Secret
		if err := c.client.Get(ctx, key, &secret); err != nil {
			c.logger.Error(err, "cluster secret poll error")

			continue
		}

		data, ok := secret.Data[dataKey]
		if !ok {
			continue
		}

		restCfg, err := clientcmd.RESTConfigFromKubeConfig([]byte(data))
		if err != nil {
			c.logger.Error(err, "config extract error")

			continue
		}

		cluster := mngr.Cluster{
			Name:        cluster.Name,
			Server:      restCfg.Host,
			BearerToken: restCfg.BearerToken,
			TLSConfig:   restCfg.TLSClientConfig,
		}

		if c.clusterChan != nil {
			c.clusterChan <- cluster
		}

		clusters = append(clusters, cluster)
	}

	c.clusters = clusters
}
