package provider

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPowerdnsZoneResource(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	zoneName := randomZoneName(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccPowerdnsZoneResourceConfig(zoneName, "localhost", "Native"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerdns_zone.test", "id", zoneName),
					resource.TestCheckResourceAttr("powerdns_zone.test", "server_id", "localhost"),
					resource.TestCheckResourceAttr("powerdns_zone.test", "name", zoneName),
					resource.TestCheckResourceAttr("powerdns_zone.test", "kind", "Native"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "powerdns_zone.test",
				ImportStateId:     "localhost/" + zoneName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccPowerdnsZoneResourceConfig(zoneName, "localhost", "Master"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerdns_zone.test", "id", zoneName),
					resource.TestCheckResourceAttr("powerdns_zone.test", "server_id", "localhost"),
					resource.TestCheckResourceAttr("powerdns_zone.test", "name", zoneName),
					resource.TestCheckResourceAttr("powerdns_zone.test", "kind", "Master"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccPowerdnsZoneResourceConfig(name, serverId, kind string) string {
	return fmt.Sprintf(`
resource "powerdns_zone" "test" {
  name = %[1]q
  server_id = %[2]q
  kind = %[3]q
}
`, name, serverId, kind)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

func randomZoneName(n int) string {
	zoneName := make([]byte, n)
	for i := range zoneName {
		zoneName[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	// Put a '.' between domain and tld and at the end of the zone name.
	// It must be ensured that sep > 0 && sep < n - 2 to get a valid zone name.
	sep := rand.Intn(n-3) + 1
	zoneName[sep] = '.'
	zoneName[n-1] = '.'

	return string(zoneName)
}
