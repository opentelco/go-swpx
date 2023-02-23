
dnc {
  addr = "localhost:1339"
}

snmp "v2c" {
  community = "xWTyZ9nA158ktJF2"
  port = 161
  timeout = "20s"
  version = 2
  retries = 3
  dynamic_repetitions = true
}

transport "telnet" {
  username = ""
  password = ""
  port = 22
  screen_length = ""
  default_prompt = "\\n([<\\[]|)(\\S+)[>#\\]](\\s|)$"
  default_errors = ""

  read_dead_line = "30s"
  write_dead_line = "10s"
  ssh_key_path = ""
}

transport "ssh" {
  username = ""
  password = ""
  port = 22
  screen_length = ""
  default_prompt = "\\n([<\\[]|)(\\S+)[>#\\]](\\s|)$"
  default_errors = ""
  read_dead_line = "30s"
  write_dead_line = "10s"
  ssh_key_path = ""
}

logger {
  level = "DEBUG"
  as_json = false
}
