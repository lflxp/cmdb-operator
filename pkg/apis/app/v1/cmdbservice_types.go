package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CmdbServiceSpec defines the desired state of CmdbService
// 自定义 API
// 需要我们根据我们的需求去自定义结构体,我们最上面预定义的资源清单中就有 size、image、ports 这些属性，所有我们需要用到的属性都需要在这个结构体中进行定义
type CmdbServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Size *int32 `json:"size"`
	// Image    string                      `json:"image"`
	Resource corev1.ResourceRequirements `json:"resources,omitempty"`
	// Envs       []corev1.EnvVar             `json:"envs,omitempty"`
	Services         []corev1.ServicePort          `json:"ports,omitempty"`
	Containers       []corev1.Container            `json:"containers"`
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
}

// CmdbServiceStatus defines the observed state of CmdbService
// 描述资源的状态
type CmdbServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	appsv1.DeploymentStatus `json:",inline"`
	PodNames                []string `json:"podnames"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CmdbService is the Schema for the cmdbservices API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=cmdbservices,scope=Namespaced
type CmdbService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CmdbServiceSpec   `json:"spec,omitempty"`
	Status CmdbServiceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CmdbServiceList contains a list of CmdbService
type CmdbServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CmdbService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CmdbService{}, &CmdbServiceList{})
}
