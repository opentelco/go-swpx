

/** SWPX Core config **/

mongo {
  server {
    addr = "localhost"
    port = 27017
  }
  database = "swpx"
  timeout = "5s"
}

logger {
  level = "DEBUG"
  as_json = false
}

/** Nats servers for the DNC **/
nats {
  server {
    addr = "localhost"
    port = "14222"
  }

  server {
    addr = "localhost"
    port = "24222"
  }

  server {
    addr = "localhost"
    port = "34222"
  }
}


// resource plugins/drivers
// vrp pins the config to a resource plugin named vrp
resource "vrp" {
  version = "v1.0.0"
  description = "switches from huawei (VRP software)"

  // pushed to plugin and dynamically loaded in the plugin
  // dnc config
  dnc {
    snmp {

    }
    connection "ssh" {

    }
    connection "telnet" {

    }

    nats {
      username = "test"
      password = "testPassword"
      server {
        addr = "localhost"
        port = "14222"
      }

      server {
        addr = "localhost"
        port = "24222"
      }

      server {
        addr = "localhost"
        port = "34222"
      }
    }
  }


}


resource "raycore" {
  version = "v1.0.0"
  description = "cpe from raycore"
}


// providers plugins/drivers
provider "vx" {
  version = "v1.0.0"
  description = "provider from VX to do lookups and other stuff from"

}