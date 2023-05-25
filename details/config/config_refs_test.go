package config

import (
	"strings"
	"testing"
)

func TestResolveReferenceWithoutContext(t *testing.T) {
	config, err := FromYamlReader(strings.NewReader(`
                baseUrl: acme.com
                serviceUrl: ${baseUrl}`))

	if err != nil {
		t.Fatal(err)
	}

	type ServiceConfig struct {
		ServiceUrl string
	}

	svcConfig := ServiceConfig{}

	err = config.Materialize(&svcConfig)

	if err != nil {
		t.Fatal(err)
	}

	if svcConfig.ServiceUrl != "acme.com" {
		t.Fatal("Reference is not resolved")
	}
}

func TestResolveReferenceWithPartialContext(t *testing.T) {
	config, err := FromYamlReader(strings.NewReader(`
                baseUrl: acme.com
                serviceUrl: api.${baseUrl}`))

	if err != nil {
		t.Fatal(err)
	}

	type ServiceConfig struct {
		ServiceUrl string
	}

	svcConfig := ServiceConfig{}

	err = config.Materialize(&svcConfig)

	if err != nil {
		t.Fatal(err)
	}

	if svcConfig.ServiceUrl != "api.acme.com" {
		t.Fatal("Reference is not resolved")
	}
}

func TestResolveReferenceWithFullContext(t *testing.T) {
	config, err := FromYamlReader(strings.NewReader(`
                baseUrl: acme.com
                serviceUrl: http://api.${baseUrl}/products/get/11`))

	if err != nil {
		t.Fatal(err)
	}

	type ServiceConfig struct {
		ServiceUrl string
	}

	svcConfig := ServiceConfig{}

	err = config.Materialize(&svcConfig)

	if err != nil {
		t.Fatal(err)
	}

	if svcConfig.ServiceUrl != "http://api.acme.com/products/get/11" {
		t.Fatal("Reference is not resolved")
	}
}

func TestResolveSeveralReferences(t *testing.T) {
	config, err := FromYamlReader(strings.NewReader(`
                chuck: "Chuck Norris"
                silvester: "Silvester Stallone"
                memeText: "${silvester}: How many push-ups can you do? ${chuck}: All of them."`))

	if err != nil {
		t.Fatal(err)
	}

	type ChuckMemes struct {
		MemeText string
	}

	chuckMemes := ChuckMemes{}

	err = config.Materialize(&chuckMemes)

	if err != nil {
		t.Fatal(err)
	}

	if chuckMemes.MemeText != "Silvester Stallone: How many push-ups can you do? Chuck Norris: All of them." {
		t.Fatal("References are not resolved")
	}
}

func TestResolveReferenceInsideUntypedMap(t *testing.T) {
	config, err := FromYamlReader(strings.NewReader(`
                baseUrl: acme.com
                values: 
                    service:
                        url: "service.${baseUrl}"`))

	if err != nil {
		t.Fatal(err)
	}

	type Config struct {
		Values map[string]interface{}
	}

	cfg := Config{}

	err = config.Materialize(&cfg)

	if err != nil {
		t.Fatal(err)
	}

	service := cfg.Values["service"].(MapStr)

	if service["url"] != "service.acme.com" {
		t.Fatal("Reference is not resolved")
	}
}

func TestResolveReferenceInsideUntypedMapArray(t *testing.T) {
	config, err := FromYamlReader(strings.NewReader(`
                baseUrl: acme.com
                values: 
                    ingress:
                        routes: 
                        - path: service.${baseUrl}`))

	if err != nil {
		t.Fatal(err)
	}

	type Config struct {
		Values map[string]interface{}
	}

	cfg := Config{}

	err = config.Materialize(&cfg)

	if err != nil {
		t.Fatal(err)
	}

	ingress := cfg.Values["ingress"].(MapStr)

	routes := ingress["routes"].([]interface{})

	route := routes[0].(MapStr)

	if route["path"] != "service.acme.com" {
		t.Fatal("Reference is not resolved")
	}
}

func TestResolveReferenceInsideUntypedMapArrayWithBooleanConversion_TrueValue(t *testing.T) {
	config, err := FromYamlReader(strings.NewReader(`
                highAvailability: true
                values: 
                    pdb:
                        create: ${highAvailability}`))

	if err != nil {
		t.Fatal(err)
	}

	type Config struct {
		Values map[string]interface{}
	}

	cfg := Config{}

	err = config.Materialize(&cfg)

	if err != nil {
		t.Fatal(err)
	}

	pdb := cfg.Values["pdb"].(MapStr)

	createPdb, isBool := pdb["create"].(bool)

	if !isBool {
		t.Fatal("Variable should be casted to bool")
	}

	if !createPdb {
		t.Fatal("Reference is not resolved")
	}
}

func TestResolveReferenceInsideUntypedMapArrayWithBooleanConversion_FalseValue(t *testing.T) {
	config, err := FromYamlReader(strings.NewReader(`
                highAvailability: false
                values: 
                    pdb:
                        create: ${highAvailability}`))

	if err != nil {
		t.Fatal(err)
	}

	type Config struct {
		Values map[string]interface{}
	}

	cfg := Config{}

	err = config.Materialize(&cfg)

	if err != nil {
		t.Fatal(err)
	}

	pdb := cfg.Values["pdb"].(MapStr)

	createPdb, isBool := pdb["create"].(bool)

	if !isBool {
		t.Fatal("Variable should be casted to bool")
	}

	if createPdb {
		t.Fatal("Reference is not resolved")
	}
}

func TestResolveBooleanReferenceWithTypeConversion(t *testing.T) {
	config, err := FromYamlReader(strings.NewReader(`
                highAvailability: true
                
                importDefinitions: ${highAvailability}`))

	if err != nil {
		t.Fatal(err)
	}

	type Config struct {
		ImportDefinitions bool
	}

	svcConfig := Config{}

	err = config.Materialize(&svcConfig)

	if err != nil {
		t.Fatal(err)
	}

	if !svcConfig.ImportDefinitions {
		t.Fatal("Reference is not resolved")
	}
}

func TestResolveReferenceWithMaterializeByPath(t *testing.T) {
	config, err := FromYamlReader(strings.NewReader(`
                baseUrl: acme.com
                services:
                  reporting:
                    serviceUrl: ${baseUrl}`))

	if err != nil {
		t.Fatal(err)
	}

	type ServiceConfig struct {
		ServiceUrl string
	}

	svcConfig := ServiceConfig{}

	ok, err := config.StrictMaterializeAt("services/reporting", &svcConfig)

	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("Should find reporting service section")
	}

	if svcConfig.ServiceUrl != "acme.com" {
		t.Fatal("Reference is not resolved")
	}
}

func TestResolveReferenceWithCascadedConfiguration(t *testing.T) {
	baseConfig, err := FromYamlReader(strings.NewReader(`baseUrl: acme.com`))

	if err != nil {
		t.Fatal(err)
	}

	reportingConfig, err := FromYamlReader(strings.NewReader(`serviceUrl: ${baseUrl}`))

	cascadedConfig := newCascade([]Configuration{baseConfig, reportingConfig})

	type ServiceConfig struct {
		ServiceUrl string
	}

	svcConfig := ServiceConfig{}

	err = cascadedConfig.Materialize(&svcConfig)

	if err != nil {
		t.Fatal(err)
	}

	if svcConfig.ServiceUrl != "acme.com" {
		t.Fatal("Reference is not resolved")
	}
}

func TestResolveReferenceWithCascadedConfigurationAndMaterializeByPath(t *testing.T) {
	baseConfig, err := FromYamlReader(strings.NewReader(`baseUrl: acme.com`))

	if err != nil {
		t.Fatal(err)
	}

	reportingConfig, err := FromYamlReader(strings.NewReader(`
                services:
                  reporting:
                    serviceUrl: ${baseUrl}`))

	cascadedConfig := newCascade([]Configuration{baseConfig, reportingConfig})

	type ServiceConfig struct {
		ServiceUrl string
	}

	svcConfig := ServiceConfig{}

	ok, err := cascadedConfig.StrictMaterializeAt("services/reporting", &svcConfig)

	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("Should find reporting service section")
	}

	if svcConfig.ServiceUrl != "acme.com" {
		t.Fatal("Reference is not resolved")
	}
}

func TestFailIfReferenceIsNotResolved(t *testing.T) {
	config, err := FromYamlReader(strings.NewReader(`
                serviceUrl: ${baseUrl}`))

	if err != nil {
		t.Fatal(err)
	}

	type ServiceConfig struct {
		ServiceUrl string
	}

	svcConfig := ServiceConfig{}

	err = config.Materialize(&svcConfig)

	if err == nil || !strings.Contains(err.Error(), "error decoding 'ServiceUrl': Reference 'baseUrl' could not be found") {
		t.Fatal("Should fail with error because reference is not defined")
	}
}

func TestResolveRecursiveReferences(t *testing.T) {
	config, err := FromYamlReader(strings.NewReader(`
                baseUrl: acme.com
                apiUrl: api.${baseUrl}
                policyUrl: ${apiUrl}/api/policy`))

	if err != nil {
		t.Fatal(err)
	}

	type Config struct {
		PolicyUrl string
	}

	cfg := Config{}

	err = config.Materialize(&cfg)

	if err != nil {
		t.Fatal(err)
	}

	if cfg.PolicyUrl != "api.acme.com/api/policy" {
		t.Fatal("Recursive reference is not resolved")
	}
}
