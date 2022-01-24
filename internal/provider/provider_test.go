package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"powerdns": func() (tfprotov6.ProviderServer, error) {
		return tfsdk.NewProtocol6Server(New("test")()), nil
	},
}

func testAccPreCheck(t *testing.T) {
	if _, ok := os.LookupEnv("POWERDNS_API_KEY"); !ok {
		t.Fatal("POWERDNS_API_KEY must be set for acceptance tests")
	}
	if _, ok := os.LookupEnv("POWERDNS_SERVER_URL"); !ok {
		t.Fatal("POWERDNS_SERVER_URL must be set for acceptance tests")
	}
}
