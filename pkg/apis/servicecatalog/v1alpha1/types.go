/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/pkg/api/v1"
)

// +genclient=true
// +nonNamespaced=true

// Broker represents an entity that provides ServiceClasses for use in the
// service catalog.
type Broker struct {
	metav1.TypeMeta `json:",inline"`
	// Non-namespaced.  The name of this resource in etcd is in ObjectMeta.Name.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BrokerSpec   `json:"spec"`
	Status BrokerStatus `json:"status"`
}

// BrokerList is a list of Brokers.
type BrokerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Broker `json:"items"`
}

// BrokerSpec represents a description of a Broker.
type BrokerSpec struct {
	// URL is the address used to communicate with the Broker.
	URL string `json:"url"`

	// AuthInfo contains the data that the service catalog should use to authenticate
	// with the Broker.
	AuthInfo *BrokerAuthInfo `json:"authInfo,omitempty"`
}

// BrokerAuthInfo is a union type that contains information on one of the authentication methods
// the the service catalog and brokers may support, according to the OpenServiceBroker API
// specification (https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md).
//
// Note that we currently restrict a single broker to have only one of these fields
// set on it.
type BrokerAuthInfo struct {
	// BasicAuthSecret is a reference to a Secret containing auth information the
	// catalog should use to authenticate to this Broker using basic auth.
	BasicAuthSecret *v1.ObjectReference `json:"basicAuthSecret,omitempty"`
}

// BrokerStatus represents the current status of a Broker.
type BrokerStatus struct {
	Conditions []BrokerCondition `json:"conditions"`

	// Checksum is the sha hash of the BrokerSpec that was last successfully
	// reconciled against the broker.
	Checksum *string `json:"checksum,omitempty"`
}

// BrokerCondition contains condition information for a Broker.
type BrokerCondition struct {
	// Type of the condition, currently ('Ready').
	Type BrokerConditionType `json:"type"`

	// Status of the condition, one of ('True', 'False', 'Unknown').
	Status ConditionStatus `json:"status"`

	// LastTransitionTime is the timestamp corresponding to the last status
	// change of this condition.
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`

	// Reason is a brief machine readable explanation for the condition's last
	// transition.
	Reason string `json:"reason"`

	// Message is a human readable description of the details of the last
	// transition, complementing reason.
	Message string `json:"message"`
}

// BrokerConditionType represents a broker condition value.
type BrokerConditionType string

const (
	// BrokerConditionReady represents the fact that a given broker condition
	// is in ready state.
	BrokerConditionReady BrokerConditionType = "Ready"
)

// ConditionStatus represents a condition's status.
type ConditionStatus string

// These are valid condition statuses. "ConditionTrue" means a resource is in
// the condition; "ConditionFalse" means a resource is not in the condition;
// "ConditionUnknown" means kubernetes can't decide if a resource is in the
// condition or not. In the future, we could add other intermediate
// conditions, e.g. ConditionDegraded.
const (
	// ConditionTrue represents the fact that a given condition is true
	ConditionTrue ConditionStatus = "True"

	// ConditionFalse represents the fact that a given condition is false
	ConditionFalse ConditionStatus = "False"

	// ConditionUnknown represents the fact that a given condition is unknown
	ConditionUnknown ConditionStatus = "Unknown"
)

// ServiceClassList is a list of ServiceClasses.
type ServiceClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []ServiceClass `json:"items"`
}

// +genclient=true
// +nonNamespaced=true

// ServiceClass represents an offering in the service catalog.
type ServiceClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// BrokerName is the reference to the Broker that provides this
	// ServiceClass.
	//
	// Immutable.
	BrokerName string `json:"brokerName"`

	// Description is a short description of this ServiceClass.
	Description string `json:"description"`

	// Bindable indicates whether a user can create bindings to an Instance
	// provisioned from this service. ServicePlan has an optional field called
	// Bindable which overrides the value of this field.
	Bindable bool `json:"bindable"`

	// Plans is the list of ServicePlans for this ServiceClass.  All
	// ServiceClasses have at least one ServicePlan.
	Plans []ServicePlan `json:"plans"`

	// PlanUpdatable indicates whether instances provisioned from this
	// ServiceClass may change ServicePlans after being provisioned.
	PlanUpdatable bool `json:"planUpdatable"`

	// ExternalID is the identity of this object for use with the OSB API.
	//
	// Immutable.
	ExternalID string `json:"externalID"`

	// ExternalMetadata is a blob of information about the ServiceClass, meant
	// to be user-facing content and display instructions.  This field may
	// contain platform-specific conventional values.
	ExternalMetadata *runtime.RawExtension `json:"externalMetadata, omitempty"`

	// Currently, this field is ALPHA: it may change or disappear at any time
	// and its data will not be migrated.
	//
	// Tags is a list of strings that represent different classification
	// attributes of the ServiceClass.  These are used in Cloud Foundry in a
	// way similar to Kubernetes labels, but they currently have no special
	// meaning in Kubernetes.
	AlphaTags []string `json:"alphaTags,omitempty"`

	// Currently, this field is ALPHA: it may change or disappear at any time
	// and its data will not be migrated.
	//
	// AlphaRequires exposes a list of Cloud Foundry-specific 'permissions'
	// that must be granted to an instance of this service within Cloud
	// Foundry.  These 'permissions' have no meaning within Kubernetes and an
	// Instance provisioned from this ServiceClass will not work correctly.
	AlphaRequires []string `json:"alphaRequires,omitempty"`
}

// ServicePlan represents a tier of a ServiceClass.
type ServicePlan struct {
	// Name is the CLI-friendly name of this ServicePlan.
	Name string `json:"name"`

	// ExternalID is the identity of this object for use with the OSB API.
	//
	// Immutable.
	ExternalID string `json:"externalID"`

	// Description is a short description of this ServicePlan.
	Description string `json:"description"`

	// Bindable indicates whether a user can create bindings to an Instance
	// using this ServicePlan.  If set, overrides the value of the
	// ServiceClass.Bindable field.
	Bindable *bool `json:"bindable,omitempty"`

	// Free indicates whether this plan is available at no cost.
	Free bool `json:"free"`

	// ExternalMetadata is a blob of information about the plan, meant to be
	// user-facing content and display instructions.  This field may contain
	// platform-specific conventional values.
	ExternalMetadata *runtime.RawExtension `json:"externalMetadata, omitempty"`

	// Currently, this field is ALPHA: it may change or disappear at any time
	// and its data will not be migrated.
	//
	// AlphaInstanceCreateParameterSchema is the schema for the parameters
	// that may be supplied when provisioning a new Instance on this plan.
	AlphaInstanceCreateParameterSchema *runtime.RawExtension `json:"alphaInstanceCreateParameterSchema,omitempty"`

	// Currently, this field is ALPHA: it may change or disappear at any time
	// and its data will not be migrated.
	//
	// AlphaInstanceUpdateParameterSchema is the schema for the parameters
	// that may be updated once an Instance has been provisioned on this plan.
	// This field only has meaning if the ServiceClass is PlanUpdatable.
	AlphaInstanceUpdateParameterSchema *runtime.RawExtension `json:"alphaInstanceUpdateParameterSchema,omitempty"`

	// Currently, this field is ALPHA: it may change or disappear at any time
	// and its data will not be migrated.
	//
	// AlphaBindingCreateParameterSchema is the schema for the parameters that
	// may be supplied binding to an Instance on this plan.
	AlphaBindingCreateParameterSchema *runtime.RawExtension `json:"alphaBindingCreateParameterSchema,omitempty"`
}

// InstanceList is a list of instances.
type InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Instance `json:"items"`
}

// +genclient=true

// Instance represents a provisioned instance of a ServiceClass.
type Instance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InstanceSpec   `json:"spec"`
	Status InstanceStatus `json:"status"`
}

// InstanceSpec represents the desired state of an Instance.
type InstanceSpec struct {
	// ServiceClassName is the reference to the ServiceClass this Instance
	// should be provisioned from.
	//
	// Immutable.
	ServiceClassName string `json:"serviceClassName"`

	// PlanName is the name of the ServicePlan this Instance should be
	// provisioned from.
	PlanName string `json:"planName"`

	// Parameters is a YAML representation of the properties to be
	// passed to the underlying broker.
	Parameters *runtime.RawExtension `json:"parameters,omitempty"`

	// ExternalID is the identity of this object for use with the OSB SB API.
	//
	// Immutable.
	ExternalID string `json:"externalID"`
}

// InstanceStatus represents the current status of an Instance.
type InstanceStatus struct {
	// Conditions is an array of InstanceConditions capturing aspects of an
	// Instance's status.
	Conditions []InstanceCondition `json:"conditions"`

	// AsyncOpInProgress is set to true if there is an ongoing async operation
	// against this Service Instance in progress.
	AsyncOpInProgress bool `json:"asyncOpInProgress"`

	// LastOperation is the string that the broker may have returned when
	// an async operation started, it should be sent back to the broker
	// on poll requests as a query param.
	LastOperation *string `json:"lastOperation,omitempty"`

	// DashboardURL is the URL of a web-based management user interface for
	// the service instance.
	DashboardURL *string `json:"dashboardURL,omitempty"`

	// Checksum is the checksum of the InstanceSpec that was last successfully
	// reconciled against the broker.
	Checksum *string `json:"checksum,omitempty"`
}

// InstanceCondition contains condition information about an Instance.
type InstanceCondition struct {
	// Type of the condition, currently ('Ready').
	Type InstanceConditionType `json:"type"`

	// Status of the condition, one of ('True', 'False', 'Unknown').
	Status ConditionStatus `json:"status"`

	// LastTransitionTime is the timestamp corresponding to the last status
	// change of this condition.
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`

	// Reason is a brief machine readable explanation for the condition's last
	// transition.
	Reason string `json:"reason"`

	// Message is a human readable description of the details of the last
	// transition, complementing reason.
	Message string `json:"message"`
}

// InstanceConditionType represents a instance condition value.
type InstanceConditionType string

const (
	// InstanceConditionReady represents that a given InstanceCondition is in
	// ready state.
	InstanceConditionReady InstanceConditionType = "Ready"
)

// BindingList is a list of Bindings.
type BindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Binding `json:"items"`
}

// +genclient=true

// Binding represents a "used by" relationship between an application and an
// Instance.
type Binding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BindingSpec   `json:"spec"`
	Status BindingStatus `json:"status"`
}

// BindingSpec represents the desired state of a Binding.
type BindingSpec struct {
	// InstanceRef is the reference to the Instance this Binding is to.
	//
	// Immutable.
	InstanceRef v1.LocalObjectReference `json:"instanceRef"`

	// Parameters is a YAML representation of the properties to be
	// passed to the underlying broker.
	Parameters *runtime.RawExtension `json:"parameters,omitempty"`

	// SecretName is the name of the secret to create in the Binding's
	// namespace that will hold the credentials associated with the Binding.
	SecretName string `json:"secretName,omitempty"`

	// ExternalID is the identity of this object for use with the OSB API.
	//
	// Immutable.
	ExternalID string `json:"externalID"`

	// Currently, this field is ALPHA: it may change or disappear at any time
	// and its data will not be migrated.
	//
	// AlphaPodPresetTemplate describes how a PodPreset should be created once
	// the Binding has been made. If supplied, a PodPreset will be created
	// using information in this field once the Binding has been made in the
	// Broker. The PodPreset will use the EnvFrom feature to expose the keys
	// from the Secret (specified by SecretName) that holds the Binding
	// information into Pods.
	//
	// In the future, we will provide a higher degree of control over the PodPreset.
	AlphaPodPresetTemplate *AlphaPodPresetTemplate `json:"alphaPodPresetTemplate,omitempty"`
}

// AlphaPodPresetTemplate represents how a PodPreset should be created for a
// Binding.
type AlphaPodPresetTemplate struct {
	// Name is the name of the PodPreset to create.
	Name string `json:"name"`
	// Selector is the LabelSelector of the PodPreset to create.
	Selector metav1.LabelSelector `json:"selector"`
}

// BindingStatus represents the current status of a Binding.
type BindingStatus struct {
	Conditions []BindingCondition `json:"conditions"`

	// Checksum is the checksum of the BindingSpec that was last successfully
	// reconciled against the broker.
	Checksum *string `json:"checksum,omitempty"`
}

// BindingCondition condition information for a Binding.
type BindingCondition struct {
	// Type of the condition, currently ('Ready').
	Type BindingConditionType `json:"type"`

	// Status of the condition, one of ('True', 'False', 'Unknown').
	Status ConditionStatus `json:"status"`

	// LastTransitionTime is the timestamp corresponding to the last status
	// change of this condition.
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`

	// Reason is a brief machine readable explanation for the condition's last
	// transition.
	Reason string `json:"reason"`

	// Message is a human readable description of the details of the last
	// transition, complementing reason.
	Message string `json:"message"`
}

// BindingConditionType represents a BindingCondition value.
type BindingConditionType string

const (
	// BindingConditionReady represents a binding condition is in ready state.
	BindingConditionReady BindingConditionType = "Ready"
)

// These are external finalizer values to service catalog, must be qualified name.
const (
	FinalizerServiceCatalog string = "kubernetes-incubator/service-catalog"
)
