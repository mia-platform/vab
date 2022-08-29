package apply

import (
	"bytes"
	"testing"
)

func TestGetResources(t *testing.T) {
	s := `apiVersion: v1
kind: Service
metadata:
 name: test
---
apiVersion: v1
kind: Service
metadata:
 name: test2
---
apiVersion: v1
kind: CustomResourceDefinition
metadata:
 name: testcrd
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: configmap
data:
  key1: aaaa---dldl
  key2: aaaa---
  a.json: '{"key":"\n---\nfoo\nbar\\n---\n"}'
  b.json: "{\"key\":\"\n---\nfoo\nbar\\n---\\n\"}"
`

	createResourcesFiles("./output", "./output/crds", "./output/res", *bytes.NewBuffer([]byte(s)))

}
