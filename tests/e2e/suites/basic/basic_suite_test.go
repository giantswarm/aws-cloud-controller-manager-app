package basic

import (
	"testing"
	"time"

	"github.com/giantswarm/apptest-framework/v3/pkg/state"
	"github.com/giantswarm/apptest-framework/v3/pkg/suite"
	"github.com/giantswarm/clustertest/v3/pkg/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	daemonSetName      = "aws-cloud-controller-manager"
	daemonSetNamespace = "kube-system"
)

func TestBasic(t *testing.T) {
	suite.New().
		WithInstallNamespace(daemonSetNamespace).
		WithIsUpgrade(false).
		Tests(func() {
			It("should have the aws-cloud-controller-manager DaemonSet running", func() {
				wcClient, err := state.GetFramework().WC(state.GetCluster().Name)
				Expect(err).NotTo(HaveOccurred())

				Eventually(func() (bool, error) {
					var ds appsv1.DaemonSet
					err := wcClient.Get(state.GetContext(), types.NamespacedName{
						Name:      daemonSetName,
						Namespace: daemonSetNamespace,
					}, &ds)
					if err != nil {
						logger.Log("DaemonSet %s not found yet: %s", daemonSetName, err.Error())
						return false, err
					}

					logger.Log("DaemonSet %s: desired=%d, ready=%d, available=%d",
						daemonSetName,
						ds.Status.DesiredNumberScheduled,
						ds.Status.NumberReady,
						ds.Status.NumberAvailable,
					)

					if ds.Status.DesiredNumberScheduled == 0 {
						return false, nil
					}

					return ds.Status.NumberReady == ds.Status.DesiredNumberScheduled, nil
				}).
					WithTimeout(10 * time.Minute).
					WithPolling(10 * time.Second).
					Should(BeTrue())
			})

			It("should have all DaemonSet pods running and ready", func() {
				wcClient, err := state.GetFramework().WC(state.GetCluster().Name)
				Expect(err).NotTo(HaveOccurred())

				Eventually(func() (bool, error) {
					var pods corev1.PodList
					err := wcClient.List(state.GetContext(), &pods,
						client.InNamespace(daemonSetNamespace),
						client.MatchingLabels{"app.kubernetes.io/name": daemonSetName},
					)
					if err != nil {
						return false, err
					}

					if len(pods.Items) == 0 {
						logger.Log("No pods found for %s", daemonSetName)
						return false, nil
					}

					for _, pod := range pods.Items {
						if pod.Status.Phase != corev1.PodRunning {
							logger.Log("Pod %s is in phase %s", pod.Name, pod.Status.Phase)
							return false, nil
						}
						for _, cs := range pod.Status.ContainerStatuses {
							if !cs.Ready {
								logger.Log("Pod %s container %s is not ready", pod.Name, cs.Name)
								return false, nil
							}
						}
					}

					logger.Log("All %d pods for %s are running and ready", len(pods.Items), daemonSetName)
					return true, nil
				}).
					WithTimeout(10 * time.Minute).
					WithPolling(10 * time.Second).
					Should(BeTrue())
			})

			It("should have nodes without the cloud-provider uninitialized taint", func() {
				wcClient, err := state.GetFramework().WC(state.GetCluster().Name)
				Expect(err).NotTo(HaveOccurred())

				Eventually(func() (bool, error) {
					var nodes corev1.NodeList
					err := wcClient.List(state.GetContext(), &nodes)
					if err != nil {
						return false, err
					}

					if len(nodes.Items) == 0 {
						logger.Log("No nodes found yet")
						return false, nil
					}

					for _, node := range nodes.Items {
						for _, taint := range node.Spec.Taints {
							if taint.Key == "node.cloudprovider.kubernetes.io/uninitialized" {
								logger.Log("Node %s still has uninitialized taint", node.Name)
								return false, nil
							}
						}
					}

					logger.Log("All %d nodes are initialized (no cloud-provider uninitialized taint)", len(nodes.Items))
					return true, nil
				}).
					WithTimeout(10 * time.Minute).
					WithPolling(10 * time.Second).
					Should(BeTrue())
			})
		}).
		Run(t, "Basic Test")
}
