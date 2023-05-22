package config

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestReadRootOptionValue(t *testing.T) {
	config, err := FromYamlReader(strings.NewReader(`imagesTag: "master.yoda"`))

	if err != nil {
		t.Fatal(err)
	}

	val, err := config.GetValues()

	if err != nil {
		t.Error("Failed to read configuration value")
	}

	if val["imagesTag"] != "master.yoda" {
		t.Error(fmt.Sprintf("Value should be 'master.yoda' but found %s", val))
	}
}

func TestReadNestedOptionValue(t *testing.T) {
	config, err := FromYamlReader(strings.NewReader(
		`global:
           imagesTag: "master.yoda"`))

	if err != nil {
		t.Fatal(err)
	}

	val, err := config.GetValues()

	if err != nil {
		t.Fatal(err)
	}

	global, ok := val["global"].(MapStr)

	if !ok {
		t.Fatal("Global block has to be map[string]interface{}")
	}

	if global["imagesTag"] != "master.yoda" {
		t.Fatal(fmt.Sprintf("Value should be 'master.yoda' but found %s", val))
	}
}

func TestBlockWithSeveralNestedConfigurations(t *testing.T) {
	config, err := FromYamlReader(strings.NewReader(`
                postgres, rmq:
                    imagesTag: "master.yoda"
        `))

	if err != nil {
		t.Fatal(err)
	}

	nestedConfig, err := config.NestedConfig("postgres")

	if err != nil {
		t.Error(err, "Failed to create nested config")
	}

	if nestedConfig == nil {
		t.Error("Empty nested config is returned")
	}

	val, err := nestedConfig.GetValues()

	if err != nil {
		t.Error("Failed to read configuration value")
	}

	if val["imagesTag"] != "master.yoda" {
		t.Error(fmt.Sprintf("Value should be 'master.yoda' but found %s", val))
	}
}

func TestMergeSimpleOptionsFromCascadeConfiguration(t *testing.T) {
	config1, err := FromYamlReader(strings.NewReader(`imageRegistry: "nuget-nsk.oneinc.local"`))

	if err != nil {
		t.Fatal(err)
	}

	config2, err := FromYamlReader(strings.NewReader(`imagesTag: "master.yoda"`))

	if err != nil {
		t.Fatal(err)
	}

	cascadeConfig, _ := Combine([]Configuration{config1, config2})

	val, err := cascadeConfig.GetValues()

	if err != nil {
		t.Fatal(err)
	}

	if val["imagesTag"] != "master.yoda" {
		t.Fatal(fmt.Sprintf("Value should be 'master.yoda' but found %s", val))
	}

	if val["imageRegistry"] != "nuget-nsk.oneinc.local" {
		t.Fatal(fmt.Sprintf("Value should be 'nuget-nsk.oneinc.local' but found %s", val))
	}
}

func TestOverridingNonEmptyOptionWithEmptyOne(t *testing.T) {
	config1, err := FromYamlReader(strings.NewReader(`MaxConcurrentConnection: 15`))

	if err != nil {
		t.Fatal(err)
	}

	config2, err := FromYamlReader(strings.NewReader(`MaxConcurrentConnection: ""`))

	if err != nil {
		t.Fatal(err)
	}

	cascadeConfig, _ := Combine([]Configuration{config1, config2})

	val, err := cascadeConfig.GetValues()

	if err != nil {
		t.Fatal(err)
	}

	if val["MaxConcurrentConnection"] != "" {
		t.Fatal(fmt.Sprintf("Value should be empty but found %s", val["MaxConcurrentConnection"]))
	}
}

func TestMergeArrayOptionsFromCascadeConfiguration(t *testing.T) {
	config1, err := FromYamlReader(strings.NewReader(`flags: ["a", "b"]`))

	if err != nil {
		t.Fatal(err)
	}

	config2, err := FromYamlReader(strings.NewReader(`flags: ["c"]`))

	if err != nil {
		t.Fatal(err)
	}

	cascadeConfig, _ := Combine([]Configuration{config1, config2})

	val, err := cascadeConfig.GetValues()

	if err != nil {
		t.Fatal(err)
	}

	flags, ok := val["flags"].([]interface{})

	if !ok {
		t.Fatal("Flags should be an array")
	}

	if flags[0] != "a" || flags[1] != "b" || flags[2] != "c" {
		t.Fatal("Flags should contain merged elements")
	}
}

func TestMergeObjectOptionsFromCascadeConfiguration(t *testing.T) {
	config1, err := FromYamlReader(strings.NewReader(`
            global: 
                imagesTag: "master.yoda"`))

	if err != nil {
		t.Fatal(err)
	}

	config2, err := FromYamlReader(strings.NewReader(`
            global:
                imageRegistry: "nuget-nsk.oneinc.local"`))

	if err != nil {
		t.Fatal(err)
	}

	cascadeConfig, _ := Combine([]Configuration{config1, config2})

	val, err := cascadeConfig.GetValues()

	if err != nil {
		t.Fatal(err)
	}

	global, ok := val["global"].(MapStr)

	if !ok {
		t.Fatal("Global block has to be map[string]interface{}")
	}

	if global["imagesTag"] != "master.yoda" {
		t.Fatal(fmt.Sprintf("Value should be 'master.yoda' but found %s", val))
	}

	if global["imageRegistry"] != "nuget-nsk.oneinc.local" {
		t.Fatal(fmt.Sprintf("Value should be 'nuget-nsk.oneinc.local' but found %s", val))
	}
}

func TestCascadeConfigurationOptionOverriding(t *testing.T) {
	config1, err := FromYamlReader(strings.NewReader(`logLevel: "info"`))

	if err != nil {
		t.Fatal(err)
	}

	config2, err := FromYamlReader(strings.NewReader(`logLevel: "debug"`))

	if err != nil {
		t.Fatal(err)
	}

	cascadeConfig, _ := Combine([]Configuration{config1, config2})

	val, err := cascadeConfig.GetValues()

	if err != nil {
		t.Fatal(err)
	}

	if val["logLevel"] != "debug" {
		t.Fatal("Option from last value should have higher priority")
	}
}

type VersionInfo struct {
	ImagesTag string
}

func TestSimpleFieldMaterialization(t *testing.T) {
	config, err := FromYamlReader(strings.NewReader(`imagesTag: "master.yoda"`))

	if err != nil {
		t.Fatal(err)
	}

	versionInfo := &VersionInfo{}

	err = config.Materialize(&versionInfo)

	if err != nil {
		t.Fatal(err)
	}

	if versionInfo.ImagesTag != "master.yoda" {
		t.Error(fmt.Sprintf("Value should be 'master.yoda' but found %s", versionInfo.ImagesTag))
	}
}

type ReleaseInfo struct {
	ServiceName string
	Version     VersionInfo
}

func TestNestedObjectMaterialization(t *testing.T) {
	config, err := FromYamlReader(
		strings.NewReader(`
            serviceName: "billing"
            version: 
                imagesTag: "master.yoda"`))

	if err != nil {
		t.Fatal(err)
	}

	releaseInfo := &ReleaseInfo{}

	err = config.Materialize(&releaseInfo)

	if err != nil {
		t.Fatal(err)
	}

	if releaseInfo.ServiceName != "billing" {
		t.Error(fmt.Sprintf("Value should be 'billing' but found %s", releaseInfo.ServiceName))
	}

	if releaseInfo.Version.ImagesTag != "master.yoda" {
		t.Error(fmt.Sprintf("Value should be 'master.yoda' but found %s", releaseInfo.Version.ImagesTag))
	}
}

func TestMaterializationByPath(t *testing.T) {
	config, err := FromYamlReader(
		strings.NewReader(`
            injection: 
                appConfig: 
                    targetContainer: "postgres"`))

	if err != nil {
		t.Fatal(err)
	}

	type AppConfig struct {
		TargetContainer string
	}

	appConfig := &AppConfig{}

	ok, err := config.StrictMaterializeAt("injection/appConfig", &appConfig)

	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("No configuration found by path injection/appConfig")
	}

	if appConfig.TargetContainer != "postgres" {
		t.Error(fmt.Sprintf("Value should be 'postgres' but found %s", appConfig.TargetContainer))
	}
}

func TestFailingMaterializationByPathIfUnusedPropertyFound(t *testing.T) {
	config, err := FromYamlReader(
		strings.NewReader(`
            injection: 
                appConfig: 
                    knownProperty: "some value"
                    unknownProperty: "another value"`))

	if err != nil {
		t.Fatal(err)
	}

	type AppConfig struct {
		KnownProperty string
	}

	appConfig := &AppConfig{}

	_, err = config.StrictMaterializeAt("injection/appConfig", &appConfig)

	if err == nil || !strings.HasSuffix(err.Error(), "invalid keys: unknownProperty") {
		t.Fatal(err)
	}
}

func TestMaterializationByMissingPath(t *testing.T) {
	config, err := FromYamlReader(
		strings.NewReader(`
            injection: {}`))

	if err != nil {
		t.Fatal(err)
	}

	type AppConfig struct {
		TargetContainer string
	}

	appConfig := &AppConfig{}

	ok, err := config.StrictMaterializeAt("injection/appConfig", &appConfig)

	if err != nil {
		t.Fatal(err)
	}

	if ok {
		t.Fatal("Materialization should not succeed because given key is missing")
	}
}

type DockerConfig struct {
	ImagesTag     string
	ImageRegistry string
}

func TestMergedCascadeConfigurationMaterialization(t *testing.T) {
	config1, err := FromYamlReader(strings.NewReader(`imagesTag: "master.yoda"`))

	if err != nil {
		t.Fatal(err)
	}

	config2, err := FromYamlReader(strings.NewReader(`imageRegistry: "nuget-nsk.oneinc.local"`))

	if err != nil {
		t.Fatal(err)
	}

	cascadeConfig, _ := Combine([]Configuration{config1, config2})

	dockerConfig := &DockerConfig{}

	err = cascadeConfig.Materialize(&dockerConfig)

	if err != nil {
		t.Fatal(err)
	}

	if dockerConfig.ImagesTag != "master.yoda" {
		t.Fatal(fmt.Sprintf("Value should be 'master.yoda' but found %s", dockerConfig.ImagesTag))
	}

	if dockerConfig.ImageRegistry != "nuget-nsk.oneinc.local" {
		t.Fatal(fmt.Sprintf("Value should be 'nuget-nsk.oneinc.local' but found %s", dockerConfig.ImageRegistry))
	}
}

func TestCascadeMaterializationByPath(t *testing.T) {
	config1, err := FromYamlReader(strings.NewReader(`
                            docker:
                                images:
                                    imageRegistry: "nuget-nsk.oneinc.local"`))
	if err != nil {
		t.Fatal(err)
	}

	config2, err := FromYamlReader(strings.NewReader(`
                                docker:
                                    images:
                                        imagesTag: "master.yoda"`))

	if err != nil {
		t.Fatal(err)
	}

	cascadeConfig, _ := Combine([]Configuration{config1, config2})

	dockerConfig := &DockerConfig{}

	ok, err := cascadeConfig.StrictMaterializeAt("docker/images", &dockerConfig)

	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("No configuration found by path docker/images")
	}

	if dockerConfig.ImagesTag != "master.yoda" {
		t.Fatal(fmt.Sprintf("Value should be 'master.yoda' but found %s", dockerConfig.ImagesTag))
	}

	if dockerConfig.ImageRegistry != "nuget-nsk.oneinc.local" {
		t.Fatal(fmt.Sprintf("Value should be 'nuget-nsk.oneinc.local' but found %s", dockerConfig.ImageRegistry))
	}
}

func TestFailingCascadeMaterializationByPathIfUnusedPropertyFound(t *testing.T) {
	config1, err := FromYamlReader(strings.NewReader(`
                            docker:
                                images:
                                    imageRegistry: "nuget-nsk.oneinc.local"`))
	if err != nil {
		t.Fatal(err)
	}

	config2, err := FromYamlReader(strings.NewReader(`
                                docker:
                                    images:
                                        imagesTag: "master.yoda"
                                        unknownProperty: "some value"`))

	if err != nil {
		t.Fatal(err)
	}

	cascadeConfig, _ := Combine([]Configuration{config1, config2})

	dockerConfig := &DockerConfig{}

	_, err = cascadeConfig.StrictMaterializeAt("docker/images", &dockerConfig)

	if err == nil || !strings.HasSuffix(err.Error(), "invalid keys: unknownProperty") {
		t.Fatal(err)
	}
}

func TestCascadeMaterializationByMissingPath(t *testing.T) {
	config1, err := FromYamlReader(strings.NewReader(`docker: {}`))
	if err != nil {
		t.Fatal(err)
	}

	config2, err := FromYamlReader(strings.NewReader(`docker: {}`))

	if err != nil {
		t.Fatal(err)
	}

	cascadeConfig, _ := Combine([]Configuration{config1, config2})

	dockerConfig := &DockerConfig{}

	ok, err := cascadeConfig.StrictMaterializeAt("docker/images", &dockerConfig)

	if err != nil {
		t.Fatal(err)
	}

	if ok {
		t.Fatal("Materialization should not succeed because given key is missing")
	}
}

func TestNestedConfiguration(t *testing.T) {
	config, err := FromYamlReader(strings.NewReader(`
                imageRegistry: "nuget-nsk.oneinc.local"
                postgres:
                    imagesTag: "master.yoda"
        `))

	if err != nil {
		t.Fatal(err)
	}

	nestedConfig, err := config.NestedConfig("postgres")

	if err != nil {
		t.Error(err, "Failed to create nested config")
	}

	if nestedConfig == nil {
		t.Error("Empty nested config is returned")
	}

	val, err := nestedConfig.GetValues()

	if err != nil {
		t.Error("Failed to read configuration value")
	}

	if val["imageRegistry"] != "nuget-nsk.oneinc.local" {
		t.Error(fmt.Sprintf("Value should be 'master.yoda' but found %s", val))
	}

	if val["imagesTag"] != "master.yoda" {
		t.Error(fmt.Sprintf("Value should be 'master.yoda' but found %s", val))
	}

	if val["postgres"] != nil {
		t.Error("postgres sub-configuration should be removed")
	}
}

func TestCascadeConfigurationNestedConfig(t *testing.T) {
	config1, err := FromYamlReader(strings.NewReader(`imageRegistry: "nuget-nsk.oneinc.local"`))

	if err != nil {
		t.Fatal(err)
	}

	config2, err := FromYamlReader(strings.NewReader(`
                            postgres:
                                imagesTag: "master.yoda"`))

	if err != nil {
		t.Fatal(err)
	}

	cascadeConfig, _ := Combine([]Configuration{config1, config2})

	nestedConfig, err := cascadeConfig.NestedConfig("postgres")

	if err != nil {
		t.Fatal(err, "Failed to create nested configuration")
	}

	dockerConfig := &DockerConfig{}

	err = nestedConfig.Materialize(&dockerConfig)

	if err != nil {
		t.Fatal(err)
	}

	if dockerConfig.ImagesTag != "master.yoda" {
		t.Fatal(fmt.Sprintf("Value should be 'master.yoda' but found %s", dockerConfig.ImagesTag))
	}

	if dockerConfig.ImageRegistry != "nuget-nsk.oneinc.local" {
		t.Fatal(fmt.Sprintf("Value should be 'nuget-nsk.oneinc.local' but found %s", dockerConfig.ImageRegistry))
	}
}

func TestWritingOfSimpleMapStr(t *testing.T) {
	config := make(MapStr)
	config["name"] = "pmnext"
	config["version"] = "future"
	config["id"] = 5

	buffer := bytes.Buffer{}

	// Act
	err := WriteMapToYaml(&config, &buffer)
	if err != nil {
		panic(err)
	}

	encodedBytes := buffer.Bytes()

	// Asserts
	assertYamlTextContainsStringProperty(t, encodedBytes, "name", "pmnext")
	assertYamlTextContainsStringProperty(t, encodedBytes, "version", "future")
	assertYamlTextContainsIntProperty(t, encodedBytes, "id", 5)
}

func assertYamlTextContainsStringProperty(t *testing.T, yamlText []byte, key string, value string) {
	matched, err := regexp.Match(fmt.Sprintf(`\s*%s\s*:\s*"?%s"?`, key, value), yamlText)
	if err != nil {
		t.Fatal(err)
	}

	if !matched {
		t.Logf(fmt.Sprintf("Encoded YAML should contains property with name '%s' and value '%s'", key, value))
		t.Fail()
	}
}

func assertYamlTextContainsIntProperty(t *testing.T, yamlText []byte, key string, value int) {
	matched, err := regexp.Match(fmt.Sprintf(`\s*%s\s*:\s*%d`, key, value), yamlText)
	if err != nil {
		t.Fatal(err)
	}

	if !matched {
		t.Logf(fmt.Sprintf("Encoded YAML should contains property with name '%s' and value '%d'", key, value))
		t.Fail()
	}
}
