http_addr = ":1337"
grpc_addr = ":1338"

// request is the default values for all request to the poller
request {
  default_request_timeout   = "90s"
  default_task_queue_prefix = "VX_SE1"
  default_provider          = "vx"
  default_cache_ttl         = "30s"
}

mongodb-cache "response" {
  collection = "response_cache"
  database   = "swpx"
}


mongodb-cache "interface" {
  database   = "swpx"
  collection = "request_cache"
}


blacklist_provider = ["vx"]

mongodb {
  addr    = "localhost"
  port    = 27019
  timeout = "5s"
  // user = "SE2_SWPX"
  // password = "RD2!fM@nQzQQ8MgrQJtTxFUJpHp4PBAK-7yEGeBcj.N-!TP-hjxevbPN-vNuxVVoE"

  database = "swpx"
}


logger {
  level   = "DEBUG"
  as_json = false
}

temporal {
  addr       = "localhost:7233"
  task_queue = "swpx-tq"
  namespace  = "TEST_SwitchPoller"
}