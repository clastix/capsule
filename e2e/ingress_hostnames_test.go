//+build e2e

// Copyright 2020-2021 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/clastix/capsule/api/v1alpha1"
)

var _ = Describe("when Tenant handles Ingress hostnames", func() {
	tnt := &v1alpha1.Tenant{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ingress-hostnames",
		},
		Spec: v1alpha1.TenantSpec{
			Owner: v1alpha1.OwnerSpec{
				Name: "hostname",
				Kind: "User",
			},
			IngressHostnames: &v1alpha1.AllowedListSpec{
				Exact:      []string{"sigs.k8s.io", "operator.sdk", "domain.tld"},
				Regex: `.*\.clastix\.io`,
			},
		},
	}

	// scaffold a basic networking.k8s.io Ingress with name and host
	networkingIngress := func(name, hostname string) *networkingv1.Ingress {
		return &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Spec: networkingv1.IngressSpec{
				Rules: []networkingv1.IngressRule{
					{
						Host: hostname,
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										Path: "/",
										PathType: func(v networkingv1.PathType) *networkingv1.PathType {
											return &v
										}(networkingv1.PathTypeExact),
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "foo",
												Port: networkingv1.ServiceBackendPort{Name: "http"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
	}
	// scaffold a basic extensions Ingress with name and host
	extensionsIngress := func(name, hostname string) *extensionsv1beta1.Ingress {
		return &extensionsv1beta1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Spec: extensionsv1beta1.IngressSpec{
				Rules: []extensionsv1beta1.IngressRule{
					{
						Host: hostname,
						IngressRuleValue: extensionsv1beta1.IngressRuleValue{
							HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
								Paths: []extensionsv1beta1.HTTPIngressPath{
									{
										Path: "/",
										PathType: func(v extensionsv1beta1.PathType) *extensionsv1beta1.PathType {
											return &v
										}(extensionsv1beta1.PathTypeExact),
										Backend: extensionsv1beta1.IngressBackend{
											ServiceName: "foo",
											ServicePort: intstr.FromInt(8080),
										},
									},
								},
							},
						},
					},
				},
			},
		}

	}

	JustBeforeEach(func() {
		EventuallyCreation(func() error {
			tnt.ResourceVersion = ""
			return k8sClient.Create(context.TODO(), tnt)
		}).Should(Succeed())
	})

	JustAfterEach(func() {
		Expect(k8sClient.Delete(context.TODO(), tnt)).Should(Succeed())
	})


	It("should block a non allowed Hostname", func() {
		maj, min, v := GetKubernetesSemVer()

		if maj == 1 && min > 18 {
			ns := NewNamespace("disallowed-hostname-networking")
			cs := ownerClient(tnt)

			NamespaceCreation(ns, tnt, defaultTimeoutInterval).Should(Succeed())
			TenantNamespaceList(tnt, defaultTimeoutInterval).Should(ContainElement(ns.GetName()))

			By("testing networking.k8s.io", func() {
				Eventually(func() (err error) {
					obj := networkingIngress("denied-networking", "kubernetes.io")
					_, err = cs.NetworkingV1().Ingresses(ns.GetName()).Create(context.TODO(), obj, metav1.CreateOptions{})
					return
				}, defaultTimeoutInterval, defaultPollInterval).ShouldNot(Succeed())
			})

			return
		}

		Skip("Running test on Kubernetes " + v + ", doesn't provide networking.k8s.io")
	})

	It("should block a non allowed Hostname", func() {
		maj, min, v := GetKubernetesSemVer()

		if maj == 1 && min < 22 {
			By("testing extensions", func() {
				ns := NewNamespace("disallowed-hostname-extensions")
				cs := ownerClient(tnt)

				NamespaceCreation(ns, tnt, defaultTimeoutInterval).Should(Succeed())
				TenantNamespaceList(tnt, defaultTimeoutInterval).Should(ContainElement(ns.GetName()))

				Eventually(func() (err error) {
					obj := extensionsIngress("denied-extensions", "kubernetes.io")
					_, err = cs.ExtensionsV1beta1().Ingresses(ns.GetName()).Create(context.TODO(), obj, metav1.CreateOptions{})
					return
				}, defaultTimeoutInterval, defaultPollInterval).ShouldNot(Succeed())
			})

			return
		}

		Skip("Running test on Kubernetes " + v + ", extensions is deprecated")
	})

	It("should allow Hostnames in list", func() {
		maj, min, v := GetKubernetesSemVer()

		if maj == 1 && min > 18 {
			ns := NewNamespace("allowed-hostname-list-networking")
			cs := ownerClient(tnt)

			NamespaceCreation(ns, tnt, defaultTimeoutInterval).Should(Succeed())
			TenantNamespaceList(tnt, defaultTimeoutInterval).Should(ContainElement(ns.GetName()))

			By("testing networking.k8s.io", func() {
				for i, h := range tnt.Spec.IngressHostnames.Exact {
					Eventually(func() (err error) {
						obj := networkingIngress(fmt.Sprintf("allowed-networking-%d", i), h)
						_, err = cs.NetworkingV1().Ingresses(ns.GetName()).Create(context.TODO(), obj, metav1.CreateOptions{})
						return
					}, defaultTimeoutInterval, defaultPollInterval).Should(Succeed())
				}
			})

			return
		}

		Skip("Running test on Kubernetes " + v + ", doesn't provide networking.k8s.io")
	})

	It("should allow Hostnames in list", func() {
		maj, min, v := GetKubernetesSemVer()

		if maj == 1 && min < 22 {
			ns := NewNamespace("allowed-hostname-list-extensions")
			cs := ownerClient(tnt)

			NamespaceCreation(ns, tnt, defaultTimeoutInterval).Should(Succeed())
			TenantNamespaceList(tnt, defaultTimeoutInterval).Should(ContainElement(ns.GetName()))

			By("testing extensions", func() {
				for i, h := range tnt.Spec.IngressHostnames.Exact {
					Eventually(func() (err error) {
						obj := extensionsIngress(fmt.Sprintf("allowed-extensions-%d", i), h)
						_, err = cs.ExtensionsV1beta1().Ingresses(ns.GetName()).Create(context.TODO(), obj, metav1.CreateOptions{})
						return
					}, defaultTimeoutInterval, defaultPollInterval).Should(Succeed())
				}
			})

			return
		}

		Skip("Running test on Kubernetes " + v + ", extensions is deprecated")
	})

	It("should allow Hostnames in regex", func() {
		maj, min, v := GetKubernetesSemVer()

		if maj == 1 && min > 18 {
			ns := NewNamespace("allowed-hostname-regex-networking")
			cs := ownerClient(tnt)

			NamespaceCreation(ns, tnt, defaultTimeoutInterval).Should(Succeed())
			TenantNamespaceList(tnt, defaultTimeoutInterval).Should(ContainElement(ns.GetName()))

			By("testing networking.k8s.io", func() {
				for _, h := range []string{"foo", "bar", "bizz"} {
					Eventually(func() (err error) {
						obj := networkingIngress(fmt.Sprintf("allowed-networking-%s", h), fmt.Sprintf("%s.clastix.io", h))
						_, err = cs.NetworkingV1().Ingresses(ns.GetName()).Create(context.TODO(), obj, metav1.CreateOptions{})
						return
					}, defaultTimeoutInterval, defaultPollInterval).Should(Succeed())
				}
			})

			return
		}

		Skip("Running test on Kubernetes " + v + ", doesn't provide networking.k8s.io")
	})

	It("should allow Hostnames in regex", func() {
		maj, min, v := GetKubernetesSemVer()

		if maj == 1 && min < 22 {
			By("testing extensions", func() {
				ns := NewNamespace("allowed-hostname-regex-extensions")
				cs := ownerClient(tnt)

				NamespaceCreation(ns, tnt, defaultTimeoutInterval).Should(Succeed())
				TenantNamespaceList(tnt, defaultTimeoutInterval).Should(ContainElement(ns.GetName()))

				for _, h := range []string{"foo", "bar", "bizz"} {
					Eventually(func() (err error) {
						obj := extensionsIngress(fmt.Sprintf("allowed-extensions-%s", h), fmt.Sprintf("%s.clastix.io", h))
						_, err = cs.ExtensionsV1beta1().Ingresses(ns.GetName()).Create(context.TODO(), obj, metav1.CreateOptions{})
						return
					}, defaultTimeoutInterval, defaultPollInterval).Should(Succeed())
				}
			})
		}

		Skip("Running test on Kubernetes " + v + ", extensions is deprecated")
	})
})
