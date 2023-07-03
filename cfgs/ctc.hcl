

dnc {
  addr = "localhost:1339"
}

snmp "v2c" {
  port = 161
  timeout = "20s"
  version = 2
  retries = 3
  dynamic_repetitions = true
}

transport "telnet" {
  port = 23
  screen_length = "terminal length 0"
  default_errors = ""
  default_prompt = "\\n([<\\[]|)(\\S+)[>#\\]](\\s|)$"
  read_dead_line = "30s"
  write_dead_line = "10s"
}

transport "ssh" {
  port = 22
  screen_length = "terminal length 0"
  default_prompt = "\\n([<\\[]|)(\\S+)[>#\\]](\\s|)$"
  default_errors = ""
  read_dead_line = "30s"
  write_dead_line = "10s"

}

logger {
  level = "DEBUG"
  as_json = false
}
