package openapi2kong

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const fixturePath = "./oas3_testfiles/"

// findFilesBySuffix returns a list of files in the fixturePath
// that end with the given suffix.
func findFilesBySuffix(t *testing.T, suffix string) []fs.DirEntry {
	files, err := os.ReadDir(fixturePath)
	if err != nil {
		t.Error("failed reading test data: %w", err)
	}

	// loop over all files, and remove anything that doesn't end with ".yaml"
	for i := 0; i < len(files); i++ {
		if !strings.HasSuffix(files[i].Name(), suffix) {
			files = append(files[:i], files[i+1:]...)
			i--
		}
	}

	return files
}

func Test_Openapi2kong(t *testing.T) {
	files := findFilesBySuffix(t, ".yaml")

	for _, file := range files {
		fileNameIn := file.Name()
		fileNameExpected := strings.TrimSuffix(fileNameIn, ".yaml") + ".expected.json"
		fileNameOut := strings.TrimSuffix(fileNameIn, ".yaml") + ".generated.json"
		dataIn, _ := os.ReadFile(fixturePath + fileNameIn)
		dataOut, err := Convert(dataIn, O2kOptions{
			Tags: []string{"OAS3_import", "OAS3file_" + fileNameIn},
		})
		if err != nil {
			t.Error(fmt.Sprintf("'%s' didn't expect error: %%w", fixturePath+fileNameIn), err)
		} else {
			JSONOut, _ := json.MarshalIndent(dataOut, "", "  ")
			os.WriteFile(fixturePath+fileNameOut, JSONOut, 0o600)
			JSONExpected, _ := os.ReadFile(fixturePath + fileNameExpected)
			assert.JSONEq(t, string(JSONExpected), string(JSONOut),
				"'%s': the JSON blobs should be equal", fixturePath+fileNameIn)
		}
	}
}
