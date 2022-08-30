package apply

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

func createResourcesFiles(outputDir string, crdPath string, resourcesPath string, buffer bytes.Buffer) (err error) {
	customResources := new(bytes.Buffer)
	resources := new(bytes.Buffer)

	fileContent, err := ioutil.ReadAll(&buffer)
	if err != nil {
		return fmt.Errorf("error reading Kustomize output: %s", err)
	}

	re := regexp.MustCompile(`\n---\n`)

	item := make(map[string]interface{})

	for _, doc := range re.Split(string(fileContent), -1) {
		yaml.Unmarshal([]byte(doc), &item)
		if item["kind"] == "CustomResourceDefinition" {
			crdYaml, err := yaml.Marshal(&item)
			if err != nil {
				return err
			}
			fmt.Fprint(customResources, string(crdYaml), "---\n")
		} else {
			resYaml, err := yaml.Marshal(&item)
			if err != nil {
				return err
			}
			fmt.Fprint(resources, string(resYaml), "---\n")
		}
	}

	err = os.MkdirAll(outputDir, folderPermissions)
	if err != nil {
		return fmt.Errorf("error creating output folder (%s): %s", outputDir, err)
	}

	if customResources.Len() != 0 {
		fmt.Println("creating Crds at ", crdPath)
		err = ioutil.WriteFile(crdPath, customResources.Bytes(), filesPermissions)
		if err != nil {
			return fmt.Errorf("error creating crd file (%s): %s", crdPath, err)
		}
	}

	fmt.Println("creating Resources at ", resourcesPath)
	if resources.Len() != 0 {
		err = ioutil.WriteFile(resourcesPath, resources.Bytes(), filesPermissions)
		if err != nil {
			return fmt.Errorf("error creating resources file (%s): %s", resourcesPath, err)
		}
	}

	return
}
