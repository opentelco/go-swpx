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

logger {
  level = "DEBUG"
  as_json = false
}
