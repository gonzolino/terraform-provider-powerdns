resource "powerdns_zone" "example_org" {
  name      = "example.org."
  server_id = "localhost"
  kind      = "Native"
}
