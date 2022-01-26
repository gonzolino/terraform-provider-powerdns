package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPowerdnsZoneResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccPowerdnsZoneResourceConfig("example.org.", "localhost", "Native"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerdns_zone.test", "id", "example.org."),
					resource.TestCheckResourceAttr("powerdns_zone.test", "server_id", "localhost"),
					resource.TestCheckResourceAttr("powerdns_zone.test", "name", "example.org."),
					resource.TestCheckResourceAttr("powerdns_zone.test", "kind", "Native"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "powerdns_zone.test",
				ImportStateId:     "localhost/example.org.",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccPowerdnsZoneResourceConfig("example.org.", "localhost", "Master"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerdns_zone.test", "id", "example.org."),
					resource.TestCheckResourceAttr("powerdns_zone.test", "server_id", "localhost"),
					resource.TestCheckResourceAttr("powerdns_zone.test", "name", "example.org."),
					resource.TestCheckResourceAttr("powerdns_zone.test", "kind", "Master"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccPowerdnsZoneResourceConfig(serverId, name, kind string) string {
	return fmt.Sprintf(`
resource "powerdns_zone" "test" {
  name = %[1]q
  server_id = %[2]q
  kind = %[3]q
}
`, serverId, name, kind)
}
