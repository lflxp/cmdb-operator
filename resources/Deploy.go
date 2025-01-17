package resources

import (
	appv1 "test/operator-study/cmdbdemo/pkg/apis/app/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewDeploy(app *appv1.CmdbService) *appsv1.Deployment {
	labels := map[string]string{"app": app.Name}
	selector := &metav1.LabelSelector{MatchLabels: labels}
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,

			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(app, schema.GroupVersionKind{
					Group:   appv1.SchemeGroupVersion.Group,
					Version: appv1.SchemeGroupVersion.Version,
					Kind:    "CmdbService",
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: app.Spec.Size,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					ImagePullSecrets: app.Spec.ImagePullSecrets,
					Containers:       app.Spec.Containers,
				},
			},
			Selector: selector,
		},
	}
}

// func newContainers(app *appv1.CmdbService) []corev1.Container {
// 	containerPorts := []corev1.ContainerPort{}
// 	for _, svcPort := range app.Spec.Ports {
// 		cport := corev1.ContainerPort{}
// 		cport.ContainerPort = svcPort.TargetPort.IntVal
// 		containerPorts = append(containerPorts, cport)
// 	}
// 	return []corev1.Container{
// 		{
// 			Name:            app.Name,
// 			Image:           app.Spec.Image,
// 			Resources:       app.Spec.Resource,
// 			Ports:           containerPorts,
// 			ImagePullPolicy: corev1.PullIfNotPresent,
// 			Env:             app.Spec.Envs,
// 		},
// 	}
// }
