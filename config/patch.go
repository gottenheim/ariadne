package config

import (
	"fmt"
	"io"
	"strings"
)

func ApplyPatchToYamlFile(yamlReader io.Reader, patches []string) ([]byte, error) {
	yamlMap, err := ReadYamlToMap(yamlReader)
	if err != nil {
		return nil, err
	}

	// Apply patches
	for _, patch := range patches {
		err = ApplyPatchToMapStr(&yamlMap, patch)
		if err != nil {
			return nil, err
		}
	}

	// Serialize back
	patchedByteData, err := MapToYamlBytes(&yamlMap)
	if err != nil {
		return nil, err
	}

	return patchedByteData, nil
}

func ApplyPatchToMapStr(yamlMap *MapStr, patch string) error {
	// Split path and value
	patchRefParts := strings.Split(patch, "=")
	if len(patchRefParts) != 2 {
		return fmt.Errorf("YAML patch reference format is incorrect, should be <yaml_path>=<value>")
	}

	yamlPath := patchRefParts[0]
	yamlValue := patchRefParts[1]

	// Split path and end property
	yamlBranchPath, yamlLeaf := CutLastPart(yamlPath, "/")

	// Make path (find or create) and set value
	path := yamlMap.MakePath(yamlBranchPath)
	path[yamlLeaf] = yamlValue

	return nil
}

func CutLastPart(val string, sep string) (string, string) {
	length := len(val)
	if length == 0 {
		return "", ""
	}

	index := strings.LastIndex(val, sep)
	if index < 0 {
		return "", val
	}

	return val[:index], val[index+1:]
}
