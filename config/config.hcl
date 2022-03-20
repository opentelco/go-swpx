

/** SWPX Core config **/
switchpoller {

  mongodb "responseCache" {
    server {
      addr = "localhost"
      port = 27017
    }
    database = "swpx"
    timeout = "5s"
  }

  mongodb "interfaceCache" {
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

  snmp {
    community = "xWTyZ9nA158ktJF2"
    timeout = "20s"
    version = 2
    retries = 3
    dynamic_repetitions = true
  }

  transport "ssh" {
    username = ""
    password = ""
    port = 22
    screen_length = ""
    default_prompt = ""
    default_errors = ""
    cache_ttl = ""
    read_dead_line = ""
    write_dead_line = ""
    ssh_key_path = ""
  }

  transport "ssh" {
    username = ""
    password = ""
    port = 23
    screen_length = ""
    default_prompt = ""
    default_errors = ""
    cache_ttl = ""
    read_dead_line = ""
    write_dead_line = ""
    ssh_key_path = ""
  }

}

// resource plugins/drivers
// vrp pins the config to a resource plugin named vrp
resource "vrp" {
  version = "v1.0.0"
  description = "switches from huawei (VRP software)"

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

  field = "field value"
}


// providers plugins/drivers
provider "vx" {
  version = "v1.0.0"
  description = "provider from VX to do lookups and other stuff from"

  field = "field value"
}

// providers plugins/drivers
provider "sait" {
  version = "v1.0.1"
  description = "provider for SAIT"

  field = "field value"

  item {
    field-a = "this is field a-1"
    field-b = "this is field b-1"
  }

  item {
    field-a = "this is field a-2"
    field-b = "this is field b-2"
  }


}
