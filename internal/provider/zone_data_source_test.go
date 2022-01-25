package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPowerdnsZoneDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccPowerdnsZoneDataSourceConfig("example.net.", "localhost"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerdns_zone.test", "id", "example.net."),
					resource.TestCheckResourceAttr("data.powerdns_zone.test", "server_id", "localhost"),
					resource.TestCheckResourceAttr("data.powerdns_zone.test", "name", "example.net."),
					resource.TestCheckResourceAttr("data.powerdns_zone.test", "kind", "Native"),
				),
			},
			{
				Config:      testAccPowerdnsZoneDataSourceConfig("unknown.net.", "localhost"),
				ExpectError: regexp.MustCompile(`Unable to get zone 'unknown.net.': .*`),
			},
			{
				Config:      testAccPowerdnsZoneDataSourceConfig("example.net.", "unknownhost"),
				ExpectError: regexp.MustCompile(`Unable to get zone 'example.net.': .*`),
			},
		},
	})
}

func testAccPowerdnsZoneDataSourceConfig(id, serverId string) string {
	return fmt.Sprintf(`
data "powerdns_zone" "test" {
  id = %[1]q
  server_id = %[2]q
}
`, id, serverId)
}
