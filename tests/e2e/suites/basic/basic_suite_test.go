package basic

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	helmv2 "github.com/fluxcd/helm-controller/api/v2"
	"github.com/giantswarm/apptest-framework/v2/pkg/state"
	"github.com/giantswarm/apptest-framework/v2/pkg/suite"
	"github.com/giantswarm/clustertest/pkg/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	isUpgrade       = false
	testNamespace   = "default"
	httpPort        = 80
	zoneIDLabelKey  = "topology.kubernetes.io/zone"
	zoneIDLabelKey2 = "failure-domain.beta.kubernetes.io/zone" // legacy label
)

func TestBasic(t *testing.T) {
	suite.New().
		// Empty string forces the framework to install it in the created cluster org namespace
		WithInstallNamespace("kube-system").
		// If this is an upgrade test or not.
		WithIsUpgrade(isUpgrade).
		WithValuesFile("./values.yaml").
		Tests(func() {
			It("should check the HelmRelease is ready", func() {
				Eventually(func() (bool, error) {
					appNamespace := state.GetCluster().Organization.GetNamespace()
					appName := fmt.Sprintf("%s-cloud-provider-aws", state.GetCluster().Name)

					mcKubeClient := state.GetFramework().MC()

					logger.Log("HelmRelease: %s/%s", appNamespace, appName)

					release := &helmv2.HelmRelease{}
					err := mcKubeClient.Get(state.GetContext(), types.NamespacedName{Name: appName, Namespace: appNamespace}, release)
					if err != nil {
						return false, err
					}

					for _, c := range release.Status.Conditions {
						if c.Type == "Ready" {
							if c.Status == "True" {
								return true, nil
							} else {
								return false, errors.New(fmt.Sprintf("HelmRelease not ready [%s]: %s", c.Reason, c.Message))
							}
						}
					}

					return false, errors.New("HelmRelease not ready")
				}).
					WithTimeout(5 * time.Minute).
					WithPolling(15 * time.Second).
					Should(BeTrue())
			})

			It("should provision Classic Load Balancer and allow external HTTP connectivity", func() {
				testLoadBalancer("classic-lb-test", "", false)
			})

			It("should provision Network Load Balancer and allow external HTTP connectivity", func() {
				testLoadBalancer("nlb-test", "nlb", false)
			})

			It("should provision internal Load Balancer and allow traffic from inside cluster", func() {
				testInternalLoadBalancer("internal-lb-test")
			})

			It("should add AWS zone-id labels to all nodes", func() {
				testNodeZoneLabels()
			})
		}).
		Run(t, "AWS CCM E2E Tests")
}

// testLoadBalancer tests external load balancer functionality (Classic or NLB)
func testLoadBalancer(testName, lbType string, internal bool) {
	wcClient, err := state.GetFramework().WC(state.GetCluster().Name)
	Expect(err).Should(Succeed())

	// Create backend deployment
	deployment := createTestDeployment(testName, testNamespace)
	Eventually(func() error {
		return wcClient.Create(state.GetContext(), deployment)
	}).
		WithTimeout(1 * time.Minute).
		WithPolling(5 * time.Second).
		Should(Succeed())

	// Wait for deployment to be ready
	Eventually(func() (bool, error) {
		return deploymentIsReady(wcClient, testNamespace, testName)
	}).
		WithTimeout(3 * time.Minute).
		WithPolling(5 * time.Second).
		Should(BeTrueBecause("We expect the backend deployment to be ready"))

	// Create LoadBalancer service
	service := createLoadBalancerService(testName, testNamespace, lbType, internal)
	Eventually(func() error {
		return wcClient.Create(state.GetContext(), service)
	}).
		WithTimeout(1 * time.Minute).
		WithPolling(5 * time.Second).
		Should(Succeed())

	// Wait for LB hostname to be assigned
	var lbHostname string
	Eventually(func() (bool, error) {
		hostname, ready, err := getLoadBalancerHostname(wcClient, testNamespace, testName)
		if ready {
			lbHostname = hostname
		}
		return ready, err
	}).
		WithTimeout(10 * time.Minute).
		WithPolling(10 * time.Second).
		Should(BeTrueBecause("We expect the LoadBalancer hostname to be set"))

	// Test HTTP connectivity from outside the cluster
	Eventually(func() (bool, error) {
		return testHTTPConnectivity(lbHostname, httpPort)
	}).
		WithTimeout(5 * time.Minute).
		WithPolling(10 * time.Second).
		Should(BeTrueBecause("We expect HTTP connectivity to the LoadBalancer"))

	// Cleanup
	cleanupResources(wcClient, testNamespace, testName)
}

// testInternalLoadBalancer tests internal load balancer with traffic from inside the cluster
func testInternalLoadBalancer(testName string) {
	wcClient, err := state.GetFramework().WC(state.GetCluster().Name)
	Expect(err).Should(Succeed())

	// Create backend deployment
	deployment := createTestDeployment(testName, testNamespace)
	Eventually(func() error {
		return wcClient.Create(state.GetContext(), deployment)
	}).
		WithTimeout(1 * time.Minute).
		WithPolling(5 * time.Second).
		Should(Succeed())

	// Wait for deployment to be ready
	Eventually(func() (bool, error) {
		return deploymentIsReady(wcClient, testNamespace, testName)
	}).
		WithTimeout(3 * time.Minute).
		WithPolling(5 * time.Second).
		Should(BeTrueBecause("We expect the backend deployment to be ready"))

	// Create internal LoadBalancer service
	service := createLoadBalancerService(testName, testNamespace, "", true)
	Eventually(func() error {
		return wcClient.Create(state.GetContext(), service)
	}).
		WithTimeout(1 * time.Minute).
		WithPolling(5 * time.Second).
		Should(Succeed())

	// Wait for LB hostname to be assigned
	var lbHostname string
	Eventually(func() (bool, error) {
		hostname, ready, err := getLoadBalancerHostname(wcClient, testNamespace, testName)
		if ready {
			lbHostname = hostname
		}
		return ready, err
	}).
		WithTimeout(10 * time.Minute).
		WithPolling(10 * time.Second).
		Should(BeTrueBecause("We expect the internal LoadBalancer hostname to be set"))

	// Create a test pod to make requests from inside the cluster
	testPodName := testName + "-client"
	testPod := createCurlPod(testPodName, testNamespace, lbHostname, httpPort)
	Eventually(func() error {
		return wcClient.Create(state.GetContext(), testPod)
	}).
		WithTimeout(1 * time.Minute).
		WithPolling(5 * time.Second).
		Should(Succeed())

	// Wait for test pod to complete and verify it succeeded
	Eventually(func() (bool, error) {
		return podCompletedSuccessfully(wcClient, testNamespace, testPodName)
	}).
		WithTimeout(5 * time.Minute).
		WithPolling(10 * time.Second).
		Should(BeTrueBecause("We expect the internal connectivity test pod to succeed"))

	// Cleanup
	cleanupResources(wcClient, testNamespace, testName)
	cleanupPod(wcClient, testNamespace, testPodName)
}

// testNodeZoneLabels validates that CCM adds AWS zone-id labels to all nodes
func testNodeZoneLabels() {
	wcClient, err := state.GetFramework().WC(state.GetCluster().Name)
	Expect(err).Should(Succeed())

	// Wait a bit for CCM to update node labels
	time.Sleep(30 * time.Second)

	// Check that all nodes have zone labels
	Eventually(func() (bool, error) {
		return allNodesHaveZoneLabels(wcClient)
	}).
		WithTimeout(5 * time.Minute).
		WithPolling(10 * time.Second).
		Should(BeTrueBecause("We expect all nodes to have AWS zone-id labels"))
}

func createTestDeployment(name, namespace string) *appsv1.Deployment {
	replicas := int32(2)
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:stable-alpine",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: httpPort,
								},
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/",
										Port: intstr.FromInt(httpPort),
									},
								},
								InitialDelaySeconds: 5,
								PeriodSeconds:       5,
							},
						},
					},
				},
			},
		},
	}
}

func deploymentIsReady(wcClient client.Client, namespace, name string) (bool, error) {
	deployment := &appsv1.Deployment{}
	err := wcClient.Get(state.GetContext(), types.NamespacedName{Name: name, Namespace: namespace}, deployment)
	if err != nil {
		return false, err
	}

	if deployment.Status.ReadyReplicas == *deployment.Spec.Replicas {
		logger.Log("Deployment %s is ready with %d replicas", name, deployment.Status.ReadyReplicas)
		return true, nil
	}
	return false, nil
}

func createLoadBalancerService(name, namespace, lbType string, internal bool) *corev1.Service {
	annotations := make(map[string]string)

	// Set LB type annotation
	if lbType == "nlb" {
		annotations["service.beta.kubernetes.io/aws-load-balancer-type"] = "nlb"
	}

	// Set internal annotation
	if internal {
		annotations["service.beta.kubernetes.io/aws-load-balancer-internal"] = "true"
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: annotations,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeLoadBalancer,
			Ports: []corev1.ServicePort{
				{
					Port:       httpPort,
					TargetPort: intstr.FromInt(httpPort),
					Protocol:   corev1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"app": name,
			},
		},
	}
}

func getLoadBalancerHostname(wcClient client.Client, namespace, name string) (string, bool, error) {
	logger.Log("Checking if Service %s has load balancer hostname set", name)
	service := &corev1.Service{}
	err := wcClient.Get(state.GetContext(), types.NamespacedName{Name: name, Namespace: namespace}, service)
	if err != nil {
		logger.Log("Failed to get Service: %v", err)
		return "", false, err
	}

	if len(service.Status.LoadBalancer.Ingress) > 0 {
		if service.Status.LoadBalancer.Ingress[0].Hostname != "" {
			hostname := service.Status.LoadBalancer.Ingress[0].Hostname
			logger.Log("Load balancer hostname found: %s", hostname)
			return hostname, true, nil
		}
		if service.Status.LoadBalancer.Ingress[0].IP != "" {
			ip := service.Status.LoadBalancer.Ingress[0].IP
			logger.Log("Load balancer IP found: %s", ip)
			return ip, true, nil
		}
	}

	return "", false, nil
}

func testHTTPConnectivity(host string, port int) (bool, error) {
	url := fmt.Sprintf("http://%s:%d/", host, port)
	logger.Log("Testing HTTP connectivity to %s", url)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		logger.Log("HTTP request failed: %v", err)
		return false, nil // Return nil error to allow retry
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		logger.Log("HTTP request successful, status: %d", resp.StatusCode)
		return true, nil
	}

	logger.Log("HTTP request returned status %d: %s", resp.StatusCode, string(body))
	return false, nil
}

func createCurlPod(name, namespace, targetHost string, targetPort int) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,
			Containers: []corev1.Container{
				{
					Name:  "curl",
					Image: "curlimages/curl:latest",
					Command: []string{
						"sh", "-c",
						fmt.Sprintf("for i in $(seq 1 30); do curl -s -o /dev/null -w '%%{http_code}' --connect-timeout 5 http://%s:%d/ && exit 0; sleep 5; done; exit 1", targetHost, targetPort),
					},
				},
			},
		},
	}
}

func podCompletedSuccessfully(wcClient client.Client, namespace, name string) (bool, error) {
	pod := &corev1.Pod{}
	err := wcClient.Get(state.GetContext(), types.NamespacedName{Name: name, Namespace: namespace}, pod)
	if err != nil {
		return false, err
	}

	if pod.Status.Phase == corev1.PodSucceeded {
		logger.Log("Test pod %s completed successfully", name)
		return true, nil
	}

	if pod.Status.Phase == corev1.PodFailed {
		logger.Log("Test pod %s failed", name)
		return false, fmt.Errorf("test pod %s failed", name)
	}

	return false, nil
}

func allNodesHaveZoneLabels(wcClient client.Client) (bool, error) {
	nodes := &corev1.NodeList{}
	err := wcClient.List(state.GetContext(), nodes)
	if err != nil {
		return false, err
	}

	if len(nodes.Items) == 0 {
		logger.Log("No nodes found in cluster")
		return false, nil
	}

	for _, node := range nodes.Items {
		hasZoneLabel := false
		for key, value := range node.Labels {
			if key == zoneIDLabelKey || key == zoneIDLabelKey2 {
				if value != "" && strings.HasPrefix(value, "us-") || strings.HasPrefix(value, "eu-") || strings.HasPrefix(value, "ap-") || strings.HasPrefix(value, "sa-") || strings.HasPrefix(value, "ca-") || strings.HasPrefix(value, "me-") || strings.HasPrefix(value, "af-") {
					hasZoneLabel = true
					logger.Log("Node %s has zone label: %s=%s", node.Name, key, value)
					break
				}
			}
		}
		if !hasZoneLabel {
			logger.Log("Node %s does not have a valid zone label", node.Name)
			return false, nil
		}
	}

	logger.Log("All %d nodes have valid AWS zone-id labels", len(nodes.Items))
	return true, nil
}

func cleanupResources(wcClient client.Client, namespace, name string) {
	ctx := context.Background()

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	if err := wcClient.Delete(ctx, service); err != nil {
		logger.Log("Failed to delete service %s: %v", name, err)
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	if err := wcClient.Delete(ctx, deployment); err != nil {
		logger.Log("Failed to delete deployment %s: %v", name, err)
	}
}

func cleanupPod(wcClient client.Client, namespace, name string) {
	ctx := context.Background()
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	if err := wcClient.Delete(ctx, pod); err != nil {
		logger.Log("Failed to delete pod %s: %v", name, err)
	}
}
