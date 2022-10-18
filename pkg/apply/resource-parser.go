package apply

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

func createResourcesFiles(outputDir string, crdPath string, resourcesPath string, buffer bytes.Buffer) (err error) {
	customResources := new(bytes.Buffer)
	resources := new(bytes.Buffer)

	fileContent, err := io.ReadAll(&buffer)
	if err != nil {
		return fmt.Errorf("error reading Kustomize output: %s", err)
	}

	re := regexp.MustCompile(`\n---\n`)

	for _, doc := range re.Split(string(fileContent), -1) {
		item := make(map[string]interface{})

		err := yaml.Unmarshal([]byte(doc), &item)
		if err != nil {
			return fmt.Errorf("error with Kustomize output: %s", err)
		}

		if item["kind"] == "CustomResourceDefinition" {
			crdYaml, err := yaml.Marshal(&item)
			if err != nil {
				return fmt.Errorf("error with Kustomize outpu: %s", err)
			}
			fmt.Fprint(customResources, string(crdYaml), "---\n")
		} else if len(item) != 0 {
			resYaml, err := yaml.Marshal(&item)
			if err != nil {
				return fmt.Errorf("error with Kustomize outpu: %s", err)
			}
			fmt.Fprint(resources, string(resYaml), "---\n")
		}
	}

	err = os.MkdirAll(outputDir, folderPermissions)
	if err != nil {
		return fmt.Errorf("error creating output folder (%s): %s", outputDir, err)
	}

	if customResources.Len() != 0 {
		err = os.WriteFile(crdPath, customResources.Bytes(), filesPermissions)
		if err != nil {
			return fmt.Errorf("error creating crd file (%s): %s", crdPath, err)
		}
	}

	if resources.Len() != 0 {
		err = os.WriteFile(resourcesPath, resources.Bytes(), filesPermissions)
		if err != nil {
			return fmt.Errorf("error creating resources file (%s): %s", resourcesPath, err)
		}
	}

	return
}
