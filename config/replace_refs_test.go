package config

import (
	"strings"
	"testing"
)

func TestReplaceReferencesInStringVariable(t *testing.T) {
	ctx, err := FromYamlReader(strings.NewReader(`BaseUrl: policynext.com`))

	if err != nil {
		t.Fatal(err, "Failed to parse deployment configuration")
	}

	cfg, err := FromYamlReader(strings.NewReader(`FrontEndUrl: api.${BaseUrl}`))

	if err != nil {
		t.Fatal(err, "Failed to parse deployment configuration")
	}

	replacedCfg, err := ResolveReferences(cfg, ctx)

	if err != nil {
		t.Fatal(err, "Failed to replace references")
	}

	values, err := replacedCfg.GetValues()

	if err != nil {
		t.Fatal(err, "Failed to get configuration values")
	}

	if values["FrontEndUrl"] != "api.policynext.com" {
		t.Fatalf("Reference is not replaced")
	}
}

func TestReplaceReferencesInObjects(t *testing.T) {
	ctx, err := FromYamlReader(strings.NewReader(`BaseUrl: policynext.com`))

	if err != nil {
		t.Fatal(err, "Failed to parse deployment configuration")
	}

	cfg, err := FromYamlReader(strings.NewReader(`
                                IdentityServer:
                                    IssuerUri: "https://identity.${BaseUrl}"`))

	if err != nil {
		t.Fatal(err, "Failed to parse deployment configuration")
	}

	replacedCfg, err := ResolveReferences(cfg, ctx)

	if err != nil {
		t.Fatal(err, "Failed to replace references")
	}

	values, err := replacedCfg.GetValues()

	if err != nil {
		t.Fatal(err, "Failed to get configuration values")
	}

	if values.FindPath("IdentityServer/IssuerUri") != "https://identity.policynext.com" {
		t.Fatalf("Reference is not replaced")
	}
}

func TestReplaceReferencesInArrays(t *testing.T) {
	ctx, err := FromYamlReader(strings.NewReader(`BaseUrl: policynext.com`))

	if err != nil {
		t.Fatal(err, "Failed to parse deployment configuration")
	}

	cfg, err := FromYamlReader(strings.NewReader(`
                                Clients:
                                  pm_ui:
                                    RedirectUris:       
                                      - "https://${BaseUrl}/silent-refresh.html"`))

	if err != nil {
		t.Fatal(err, "Failed to parse deployment configuration")
	}

	replacedCfg, err := ResolveReferences(cfg, ctx)

	if err != nil {
		t.Fatal(err, "Failed to replace references")
	}

	values, err := replacedCfg.GetValues()

	if err != nil {
		t.Fatal(err, "Failed to get configuration values")
	}

	redirectUris, ok := values.FindPath("Clients/pm_ui/RedirectUris").([]interface{})

	if !ok {
		t.Fatalf("RedirectUris array could not be found")
	}

	if redirectUris[0] != "https://policynext.com/silent-refresh.html" {
		t.Fatalf("Reference is not replaced")
	}
}
