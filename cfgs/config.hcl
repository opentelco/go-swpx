http_addr = ":1337"
grpc_addr = ":1338"

// request is the default values for all request to the poller
request {
  default_request_timeout = "90s"
  default_task_queue_prefix = "VX_SE1"
  default_provider = "vx"
  default_cache_ttl = "30s"
  override_ok_list = ["mulbarton-migration", "only-for-migration", "deltapark-a18"]
}

mongodb-cache "response" {
  collection = "response_cache"
  database = "swpx"
}


mongodb-cache "interface" {
  database = "swpx"
  collection = "request_cache"
}

mongodb {
  addr = "localhost"
  port = 27019
  timeout = "5s"
}


logger {
  level = "DEBUG"
  as_json = false
}
