

/** SWPX Core config **/
switchpoller {

  mongodb "responseCache" {
    server {
      addr = "localhost"
      port = 27017
    }
    collection = "response_cache"
    database = "swpx"
    timeout = "5s"
  }


  mongodb "interfaceCache" {
    server {
      addr = "localhost"
      port = 27017
    }

    collection = "interface_cache"
    database = "swpx"
    timeout = "5s"
  }

  request {
    default_request_timeout = "10s"
    default_provider = "vx"
    default_cache_ttl = "5m"
  }

  logger {
    level = "DEBUG"
    as_json = false
  }


  snmp "v2c"{
    community = "xWTyZ9nA158ktJF2"
    port = 161
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

  snmp "v2c" {
    community = "xWTyZ9nA158ktJF2"
    port = 161
    timeout = "20s"
    version = 2
    retries = 3
    dynamic_repetitions = true
  }

  transport "telnet" {
    username = "telnetUser"
    password = ""
    port = 22
    screen_length = ""
    default_prompt = ""
    default_errors = ""

    read_dead_line = ""
    write_dead_line = ""
    ssh_key_path = ""
  }

  transport "ssh" {
    username = "sshUser"
    password = ""
    port = 23
    screen_length = ""
    default_prompt = ""
    default_errors = ""
    read_dead_line = ""
    write_dead_line = ""
    ssh_key_path = ""
  }

}


// resource plugins/drivers
// ctc pins the config to a resource plugin named ctc
resource "ctc" {
  version = "v1.0.0"
  description = "switches from CTC"

  snmp "v2c" {
    community = "xWTyZ9nA158ktJF2"
    port = 161
    timeout = "20s"
    version = 2
    retries = 3
    dynamic_repetitions = true
  }

  transport "telnet" {
    username = "telnetUser"
    password = ""
    port = 22
    screen_length = ""
    default_prompt = ""
    default_errors = ""

    read_dead_line = ""
    write_dead_line = ""
    ssh_key_path = ""
  }

  transport "ssh" {
    username = "sshUser"
    password = ""
    port = 23
    screen_length = ""
    default_prompt = ""
    default_errors = ""
    read_dead_line = ""
    write_dead_line = ""
    ssh_key_path = ""
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
