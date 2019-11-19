package resource

import (
	"fmt"

	rabbitmqv1beta1 "github.com/pivotal/rabbitmq-for-kubernetes/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type RabbitmqCluster struct {
	Instance                 *rabbitmqv1beta1.RabbitmqCluster
	StatefulSetConfiguration StatefulSetConfiguration
	DefaultConfiguration     DefaultConfiguration
}

type DefaultConfiguration struct {
	ServiceAnnotations     map[string]string
	ServiceType            string
	OperatorRegistrySecret *corev1.Secret
}

func (cluster *RabbitmqCluster) Resources() (resources []runtime.Object, err error) {
	serverConf := cluster.ServerConfigMap()
	resources = append(resources, serverConf)

	ingressService := cluster.IngressService()
	resources = append(resources, ingressService)

	headlessService := cluster.HeadlessService()
	resources = append(resources, headlessService)

	adminSecret, err := cluster.AdminSecret()
	if err != nil {
		return nil, fmt.Errorf("failed to generate admin secret: %v ", err)
	}
	resources = append(resources, adminSecret)

	erlangCookie, err := cluster.ErlangCookie()
	if err != nil {
		return nil, fmt.Errorf("failed to generate erlang cookie: %v ", err)
	}
	resources = append(resources, erlangCookie)

	serviceAccount := cluster.ServiceAccount()
	resources = append(resources, serviceAccount)

	role := cluster.Role()
	resources = append(resources, role)

	roleBinding := cluster.RoleBinding()
	resources = append(resources, roleBinding)

	if cluster.DefaultConfiguration.OperatorRegistrySecret != nil {
		clusterRegistrySecret := cluster.RegistrySecret()
		resources = append(resources, clusterRegistrySecret)
	}

	statefulSet, err := cluster.StatefulSet()
	if err != nil {
		return nil, fmt.Errorf("failed to generate StatefulSet: %v ", err)
	}

	resources = append(resources, statefulSet)

	return resources, nil
}
