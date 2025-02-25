package clustersmngr_test

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/rand"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta2"
	"github.com/fluxcd/pkg/apis/meta"
	"github.com/weaveworks/weave-gitops/core/clustersmngr"
	"github.com/weaveworks/weave-gitops/pkg/server/auth"
)

func TestClientGet(t *testing.T) {
	g := NewGomegaWithT(t)
	ns := createNamespace(g)

	clusterName := "mycluster"

	appName := "myapp" + rand.String(5)

	clientsPool := clustersmngr.NewClustersClientsPool()

	err := clientsPool.Add(
		clustersmngr.ClientConfigWithUser(&auth.UserPrincipal{}),
		clustersmngr.Cluster{
			Name:      clusterName,
			Server:    k8sEnv.Rest.Host,
			TLSConfig: k8sEnv.Rest.TLSClientConfig,
		},
	)

	g.Expect(err).To(BeNil())

	clustersClient := clustersmngr.NewClient(clientsPool)

	kust := &kustomizev1.Kustomization{
		ObjectMeta: v1.ObjectMeta{
			Name:      appName,
			Namespace: ns.Name,
		},
		Spec: kustomizev1.KustomizationSpec{
			SourceRef: kustomizev1.CrossNamespaceSourceReference{
				Kind: "GitRepository",
			},
		},
	}
	ctx := context.Background()
	g.Expect(k8sEnv.Client.Create(ctx, kust)).To(Succeed())

	k := &kustomizev1.Kustomization{}

	g.Expect(clustersClient.Get(ctx, clusterName, types.NamespacedName{Name: appName, Namespace: ns.Name}, k)).To(Succeed())
	g.Expect(k.Name).To(Equal(appName))
}

func TestClientClusteredList(t *testing.T) {
	g := NewGomegaWithT(t)
	ns := createNamespace(g)

	clusterName := "mycluster"
	appName := "myapp" + rand.String(5)

	clientsPool := clustersmngr.NewClustersClientsPool()

	err := clientsPool.Add(
		clustersmngr.ClientConfigWithUser(&auth.UserPrincipal{}),
		clustersmngr.Cluster{
			Name:      clusterName,
			Server:    k8sEnv.Rest.Host,
			TLSConfig: k8sEnv.Rest.TLSClientConfig,
		},
	)
	g.Expect(err).To(BeNil())

	clustersClient := clustersmngr.NewClient(clientsPool)

	kust := &kustomizev1.Kustomization{
		ObjectMeta: v1.ObjectMeta{
			Name:      appName,
			Namespace: ns.Name,
		},
		Spec: kustomizev1.KustomizationSpec{
			SourceRef: kustomizev1.CrossNamespaceSourceReference{
				Kind: "GitRepository",
			},
		},
	}
	ctx := context.Background()
	g.Expect(k8sEnv.Client.Create(ctx, kust)).To(Succeed())

	cklist := clustersmngr.NewClusteredList(func() client.ObjectList {
		return &kustomizev1.KustomizationList{}
	})

	g.Expect(clustersClient.ClusteredList(ctx, cklist, client.InNamespace(ns.Name))).To(Succeed())

	klist := cklist.Lists()[clusterName].(*kustomizev1.KustomizationList)

	g.Expect(klist.Items).To(HaveLen(1))
	g.Expect(klist.Items[0].Name).To(Equal(appName))

	gitRepo := &sourcev1.GitRepository{
		ObjectMeta: v1.ObjectMeta{
			Name:      appName,
			Namespace: ns.Name,
		},
		Spec: sourcev1.GitRepositorySpec{
			URL: "https://example.com/repo",
			SecretRef: &meta.LocalObjectReference{
				Name: "somesecret",
			},
		},
	}

	g.Expect(k8sEnv.Client.Create(ctx, gitRepo)).To(Succeed())

	cgrlist := clustersmngr.NewClusteredList(func() client.ObjectList {
		return &sourcev1.GitRepositoryList{}
	})

	g.Expect(clustersClient.ClusteredList(ctx, cgrlist)).To(Succeed())

	glist := cgrlist.Lists()[clusterName].(*sourcev1.GitRepositoryList)
	g.Expect(glist.Items).To(HaveLen(1))
	g.Expect(glist.Items[0].Name).To(Equal(appName))
}

func TestClientList(t *testing.T) {
	g := NewGomegaWithT(t)
	ns := createNamespace(g)

	clusterName := "mycluster"
	appName := "myapp" + rand.String(5)

	clientsPool := clustersmngr.NewClustersClientsPool()

	err := clientsPool.Add(
		clustersmngr.ClientConfigWithUser(&auth.UserPrincipal{}),
		clustersmngr.Cluster{
			Name:      clusterName,
			Server:    k8sEnv.Rest.Host,
			TLSConfig: k8sEnv.Rest.TLSClientConfig,
		},
	)
	g.Expect(err).To(BeNil())

	clustersClient := clustersmngr.NewClient(clientsPool)

	kust := &kustomizev1.Kustomization{
		ObjectMeta: v1.ObjectMeta{
			Name:      appName,
			Namespace: ns.Name,
		},
		Spec: kustomizev1.KustomizationSpec{
			SourceRef: kustomizev1.CrossNamespaceSourceReference{
				Kind: "GitRepository",
			},
		},
	}
	ctx := context.Background()
	g.Expect(k8sEnv.Client.Create(ctx, kust)).To(Succeed())

	list := &kustomizev1.KustomizationList{}

	g.Expect(clustersClient.List(ctx, clusterName, list, client.InNamespace(ns.Name))).To(Succeed())

	g.Expect(list.Items).To(HaveLen(1))
	g.Expect(list.Items[0].Name).To(Equal(appName))
}

func TestClientCreate(t *testing.T) {
	g := NewGomegaWithT(t)
	ns := createNamespace(g)

	clusterName := "mycluster"
	appName := "myapp" + rand.String(5)

	clientsPool := clustersmngr.NewClustersClientsPool()

	err := clientsPool.Add(
		clustersmngr.ClientConfigWithUser(&auth.UserPrincipal{}),
		clustersmngr.Cluster{
			Name:      clusterName,
			Server:    k8sEnv.Rest.Host,
			TLSConfig: k8sEnv.Rest.TLSClientConfig,
		},
	)
	g.Expect(err).To(BeNil())

	clustersClient := clustersmngr.NewClient(clientsPool)

	kust := &kustomizev1.Kustomization{
		ObjectMeta: v1.ObjectMeta{
			Name:      appName,
			Namespace: ns.Name,
		},
		Spec: kustomizev1.KustomizationSpec{
			SourceRef: kustomizev1.CrossNamespaceSourceReference{
				Kind: "GitRepository",
			},
		},
	}
	ctx := context.Background()

	g.Expect(clustersClient.Create(ctx, clusterName, kust)).To(Succeed())

	k := &kustomizev1.Kustomization{}
	g.Expect(clustersClient.Get(ctx, clusterName, types.NamespacedName{Name: appName, Namespace: ns.Name}, k)).To(Succeed())
	g.Expect(k.Name).To(Equal(appName))
}

func TestClientDelete(t *testing.T) {
	g := NewGomegaWithT(t)
	ns := createNamespace(g)

	clusterName := "mycluster"
	appName := "myapp" + rand.String(5)

	clientsPool := clustersmngr.NewClustersClientsPool()

	err := clientsPool.Add(
		clustersmngr.ClientConfigWithUser(&auth.UserPrincipal{}),
		clustersmngr.Cluster{
			Name:      clusterName,
			Server:    k8sEnv.Rest.Host,
			TLSConfig: k8sEnv.Rest.TLSClientConfig,
		},
	)
	g.Expect(err).To(BeNil())

	clustersClient := clustersmngr.NewClient(clientsPool)

	kust := &kustomizev1.Kustomization{
		ObjectMeta: v1.ObjectMeta{
			Name:      appName,
			Namespace: ns.Name,
		},
		Spec: kustomizev1.KustomizationSpec{
			SourceRef: kustomizev1.CrossNamespaceSourceReference{
				Kind: "GitRepository",
			},
		},
	}
	ctx := context.Background()

	g.Expect(k8sEnv.Client.Create(ctx, kust)).To(Succeed())

	g.Expect(clustersClient.Delete(ctx, clusterName, kust)).To(Succeed())
}

func TestClientUpdate(t *testing.T) {
	g := NewGomegaWithT(t)
	ns := createNamespace(g)

	clusterName := "mycluster"
	appName := "myapp" + rand.String(5)

	clientsPool := clustersmngr.NewClustersClientsPool()

	err := clientsPool.Add(
		clustersmngr.ClientConfigWithUser(&auth.UserPrincipal{}),
		clustersmngr.Cluster{
			Name:      clusterName,
			Server:    k8sEnv.Rest.Host,
			TLSConfig: k8sEnv.Rest.TLSClientConfig,
		},
	)
	g.Expect(err).To(BeNil())

	clustersClient := clustersmngr.NewClient(clientsPool)

	kust := &kustomizev1.Kustomization{
		ObjectMeta: v1.ObjectMeta{
			Name:      appName,
			Namespace: ns.Name,
		},
		Spec: kustomizev1.KustomizationSpec{
			Path: "/foo",
			SourceRef: kustomizev1.CrossNamespaceSourceReference{
				Kind: "GitRepository",
			},
		},
	}
	ctx := context.Background()
	g.Expect(k8sEnv.Client.Create(ctx, kust)).To(Succeed())

	kust.Spec.Path = "/bar"
	g.Expect(clustersClient.Update(ctx, clusterName, kust)).To(Succeed())

	k := &kustomizev1.Kustomization{}
	g.Expect(k8sEnv.Client.Get(ctx, types.NamespacedName{Name: appName, Namespace: ns.Name}, k)).To(Succeed())
	g.Expect(k.Spec.Path).To(Equal("/bar"))
}

func TestClientPatch(t *testing.T) {
	g := NewGomegaWithT(t)
	ns := createNamespace(g)

	clusterName := "mycluster"
	appName := "myapp" + rand.String(5)

	clientsPool := clustersmngr.NewClustersClientsPool()

	err := clientsPool.Add(
		clustersmngr.ClientConfigWithUser(&auth.UserPrincipal{}),
		clustersmngr.Cluster{
			Name:      clusterName,
			Server:    k8sEnv.Rest.Host,
			TLSConfig: k8sEnv.Rest.TLSClientConfig,
		},
	)
	g.Expect(err).To(BeNil())

	clustersClient := clustersmngr.NewClient(clientsPool)

	kust := &kustomizev1.Kustomization{
		TypeMeta: metav1.TypeMeta{
			Kind:       kustomizev1.KustomizationKind,
			APIVersion: kustomizev1.GroupVersion.String(),
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      appName,
			Namespace: ns.Name,
		},
		Spec: kustomizev1.KustomizationSpec{
			Path: "/foo",
			SourceRef: kustomizev1.CrossNamespaceSourceReference{
				Kind: "GitRepository",
			},
		},
	}

	ctx := context.Background()
	opt := []client.PatchOption{
		client.ForceOwnership,
		client.FieldOwner("test"),
	}
	g.Expect(clustersClient.Patch(ctx, clusterName, kust, client.Apply, opt...)).To(Succeed())

	k := &kustomizev1.Kustomization{}
	g.Expect(k8sEnv.Client.Get(ctx, types.NamespacedName{Name: appName, Namespace: ns.Name}, k)).To(Succeed())
	g.Expect(k.Spec.Path).To(Equal("/foo"))
}

func createNamespace(g *GomegaWithT) *corev1.Namespace {
	ns := &corev1.Namespace{}
	ns.Name = "kube-test-" + rand.String(5)

	g.Expect(k8sEnv.Client.Create(context.Background(), ns)).To(Succeed())

	return ns
}
