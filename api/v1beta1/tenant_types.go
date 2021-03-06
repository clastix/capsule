// Copyright 2020-2021 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TenantSpec defines the desired state of Tenant
type TenantSpec struct {
	Owners []OwnerSpec `json:"owners"`

	//+kubebuilder:validation:Minimum=1
	NamespaceQuota         *int32                       `json:"namespaceQuota,omitempty"`
	NamespacesMetadata     *AdditionalMetadataSpec      `json:"namespacesMetadata,omitempty"`
	ServicesMetadata       *AdditionalMetadataSpec      `json:"servicesMetadata,omitempty"`
	StorageClasses         *AllowedListSpec             `json:"storageClasses,omitempty"`
	IngressClasses         *AllowedListSpec             `json:"ingressClasses,omitempty"`
	IngressHostnames       *AllowedListSpec             `json:"ingressHostnames,omitempty"`
	ContainerRegistries    *AllowedListSpec             `json:"containerRegistries,omitempty"`
	NodeSelector           map[string]string            `json:"nodeSelector,omitempty"`
	NetworkPolicies        *NetworkPolicySpec           `json:"networkPolicies,omitempty"`
	LimitRanges            *LimitRangesSpec             `json:"limitRanges,omitempty"`
	ResourceQuota          *ResourceQuotaSpec           `json:"resourceQuotas,omitempty"`
	AdditionalRoleBindings []AdditionalRoleBindingsSpec `json:"additionalRoleBindings,omitempty"`
	ExternalServiceIPs     *ExternalServiceIPsSpec      `json:"externalServiceIPs,omitempty"`
	ImagePullPolicies      []ImagePullPolicySpec        `json:"imagePullPolicies,omitempty"`
	PriorityClasses        *AllowedListSpec             `json:"priorityClasses,omitempty"`

	//+kubebuilder:default=true
	EnableNodePorts bool `json:"enableNodePorts,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion
// +kubebuilder:resource:scope=Cluster,shortName=tnt
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.state",description="The actual state of the Tenant"
// +kubebuilder:printcolumn:name="Namespace quota",type="integer",JSONPath=".spec.namespaceQuota",description="The max amount of Namespaces can be created"
// +kubebuilder:printcolumn:name="Namespace count",type="integer",JSONPath=".status.size",description="The total amount of Namespaces in use"
// +kubebuilder:printcolumn:name="Node selector",type="string",JSONPath=".spec.nodeSelector",description="Node Selector applied to Pods"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Age"

// Tenant is the Schema for the tenants API
type Tenant struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TenantSpec   `json:"spec,omitempty"`
	Status TenantStatus `json:"status,omitempty"`
}

func (t *Tenant) Hub() {}

//+kubebuilder:object:root=true

// TenantList contains a list of Tenant
type TenantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Tenant `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Tenant{}, &TenantList{})
}
