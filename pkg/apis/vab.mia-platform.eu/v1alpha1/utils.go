package v1alpha1

const (
	kind    = "ClustersConfiguration"
	version = "vab.mia-platform.eu/v1alpha1"
)

func EmptyConfig(name string) ClustersConfiguration {
	return ClustersConfiguration{
		TypeMeta: TypeMeta{
			Kind:       kind,
			APIVersion: version,
		},
		Name: name,
		Spec: ConfigSpec{
			Modules: make(map[string]Module),
			AddOns:  make(map[string]AddOn),
			Groups:  make([]Group, 0),
		},
	}
}
