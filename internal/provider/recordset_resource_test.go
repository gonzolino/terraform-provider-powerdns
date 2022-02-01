package provider

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPowerdnsRecordsetResource(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	recordsetName := randomRecordsetName(5, "example.net.")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccPowerdnsRecordsetResourceConfig("example.net.", "localhost", recordsetName, "A", 500, []string{"192.168.0.3"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerdns_recordset.test", "zone_id", "example.net."),
					resource.TestCheckResourceAttr("powerdns_recordset.test", "server_id", "localhost"),
					resource.TestCheckResourceAttr("powerdns_recordset.test", "name", recordsetName),
					resource.TestCheckResourceAttr("powerdns_recordset.test", "type", "A"),
					resource.TestCheckResourceAttr("powerdns_recordset.test", "ttl", "500"),
					// resource.TestCheckResourceAttr("powerdns_recordset.test", "records", "[\"192.168.0.3\"]"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "powerdns_recordset.test",
				ImportStateId:     fmt.Sprintf("localhost/example.net./%s/A", recordsetName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccPowerdnsRecordsetResourceConfig("example.net.", "localhost", recordsetName, "A", 800, []string{"192.168.0.2", "192.168.0.4"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerdns_recordset.test", "zone_id", "example.net."),
					resource.TestCheckResourceAttr("powerdns_recordset.test", "server_id", "localhost"),
					resource.TestCheckResourceAttr("powerdns_recordset.test", "name", recordsetName),
					resource.TestCheckResourceAttr("powerdns_recordset.test", "type", "A"),
					resource.TestCheckResourceAttr("powerdns_recordset.test", "ttl", "800"),
					// resource.TestCheckResourceAttr("powerdns_recordset.test", "records", "[\"192.168.0.2\", \"192.168.0.4\"]"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccPowerdnsRecordsetResourceConfig(zoneId, serverId, name, typ string, ttl int64, records []string) string {
	recordBuilder := strings.Builder{}
	for i, r := range records {
		recordBuilder.WriteString(fmt.Sprintf("%q", r))
		if i < len(records)-1 {
			recordBuilder.WriteRune(',')
		}
	}
	return fmt.Sprintf(`
resource "powerdns_recordset" "test" {
  zone_id = %[1]q
  server_id = %[2]q
  name = %[3]q
  type = %[4]q
  ttl = %[5]d
  records = [%[6]s]
}
`, zoneId, serverId, name, typ, ttl, recordBuilder.String())
}

func randomRecordsetName(n int, zoneId string) string {
	recordsetName := make([]byte, n)
	for i := range recordsetName {
		recordsetName[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return fmt.Sprintf("%s.%s", recordsetName, zoneId)
}
