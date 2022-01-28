data "powerdns_zone" "example_net" {
  id        = "example.net."
  server_id = "localhost"
}

data "powerdns_recordset" "www_example_net" {
  zone_id   = data.powerdns_zone.example_net.id
  server_id = data.powerdns_zone.example_net.server_id
  name      = "www.example.net."
}

data "powerdns_recordset" "example_net_soa" {
  zone_id   = data.powerdns_zone.example_net.id
  server_id = data.powerdns_zone.example_net.server_id
  name      = "example.net."
  type      = "SOA"
}
