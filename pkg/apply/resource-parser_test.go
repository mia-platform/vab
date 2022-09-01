package apply

import (
	"bytes"
	"os"
	"testing"

	"github.com/mia-platform/vab/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestGetResources(t *testing.T) {
	testPath := testutils.GetTestFile("apply", "resource-parser-test", "resources.yaml")
	contentFile, err := os.ReadFile(testPath)
	assert.NoError(t, err)

	err = createResourcesFiles("./output", "./output/crds", "./output/res", *bytes.NewBuffer(contentFile))
	assert.NoError(t, err)

	assert.FileExists(t, "./output/crds")
	assert.FileExists(t, "./output/res")

	outFile1, err := os.ReadFile("./output/crds")
	assert.NoError(t, err)
	outFile2, err := os.ReadFile("./output/res")
	assert.NoError(t, err)

	expectedFile1, err := os.ReadFile("../../tests/apply/resource-parser-test/expected-crds.yaml")
	assert.NoError(t, err)
	expectedFile2, err := os.ReadFile("../../tests/apply/resource-parser-test/expected-resources.yaml")
	assert.NoError(t, err)

	assert.Equal(t, outFile1, expectedFile1)
	assert.Equal(t, outFile2, expectedFile2)

	os.RemoveAll("./output")
}
