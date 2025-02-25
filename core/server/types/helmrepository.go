package types

import (
	"github.com/fluxcd/source-controller/api/v1beta1"
	pb "github.com/weaveworks/weave-gitops/pkg/api/core"
)

func HelmRepositoryToProto(helmRepository *v1beta1.HelmRepository, clusterName string) *pb.HelmRepository {
	return &pb.HelmRepository{
		Name:          helmRepository.Name,
		Namespace:     helmRepository.Namespace,
		Url:           helmRepository.Spec.URL,
		Interval:      durationToInterval(helmRepository.Spec.Interval),
		Conditions:    mapConditions(helmRepository.Status.Conditions),
		Suspended:     helmRepository.Spec.Suspend,
		LastUpdatedAt: lastUpdatedAt(helmRepository),
		ClusterName:   clusterName,
	}
}
