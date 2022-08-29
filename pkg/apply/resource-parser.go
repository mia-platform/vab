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
	fmt.Println(crdPath, resourcesPath)

	fileContent, err := ioutil.ReadAll(&buffer)
	if err != nil {
		return err
	}
	fmt.Println(string(fileContent))

	re := regexp.MustCompile(`\n---\n`)

	item := make(map[string]interface{})

	for _, doc := range re.Split(string(fileContent), -1) {
		yaml.Unmarshal([]byte(doc), &item)
		fmt.Println(item)
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
	fmt.Println(resources)

	err = os.MkdirAll(outputDir, folderPermissions)
	if err != nil {
		return err
	}

	fmt.Println("creating Crds at ", crdPath)
	err = ioutil.WriteFile(crdPath, customResources.Bytes(), filesPermissions)
	if err != nil {
		return err
	}

	fmt.Println("creating Resources at ", resourcesPath)
	err = ioutil.WriteFile(resourcesPath, resources.Bytes(), filesPermissions)
	if err != nil {
		return err
	}

	return
}
