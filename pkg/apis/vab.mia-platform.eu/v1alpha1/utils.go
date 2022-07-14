package v1alpha1

const (
	Kind    = "ClustersConfiguration"
	Version = "vab.mia-platform.eu/v1alpha1"
)

// EmptyConfig generates an empty ClustersConfiguration
func EmptyConfig(name string) *ClustersConfiguration {
	return &ClustersConfiguration{
		TypeMeta: TypeMeta{
			Kind:       Kind,
			APIVersion: Version,
		},
		Name: name,
		Spec: ConfigSpec{
			Modules: make(map[string]Module),
			AddOns:  make(map[string]AddOn),
			Groups:  make([]Group, 0),
		},
	}
}
