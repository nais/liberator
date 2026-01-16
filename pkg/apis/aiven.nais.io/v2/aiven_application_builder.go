package aiven_nais_io_v2

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type AivenApplicationBuilder struct {
	application AivenApplication
}

func NewAivenApplicationBuilder(name, namespace string) AivenApplicationBuilder {
	return AivenApplicationBuilder{
		application: AivenApplication{
			TypeMeta: metav1.TypeMeta{
				Kind:       "AivenApplication",
				APIVersion: "aiven.nais.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec:   AivenApplicationSpec{},
			Status: AivenApplicationStatus{},
		},
	}
}

func (b AivenApplicationBuilder) WithSpec(spec AivenApplicationSpec) AivenApplicationBuilder {
	b.application.Spec = spec
	return b
}

func (b AivenApplicationBuilder) WithStatus(status AivenApplicationStatus) AivenApplicationBuilder {
	b.application.Status = status
	return b
}

func (b AivenApplicationBuilder) WithAnnotation(key, value string) AivenApplicationBuilder {
	b.application.SetAnnotations(map[string]string{
		key: value,
	})
	return b
}

func (b AivenApplicationBuilder) Build() AivenApplication {
	return b.application
}
