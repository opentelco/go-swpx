group = "pollers"
rpc_port = "1337"
receiver_queues = ["requests.cpe"]
response_group = "responses"


nats_server {
    addr = "nats-1"
    port = 4222    
} 

logging {
    as_json = false
    level = "debug"
}


distributed_lock {
    enabled =  true
    ttl = "30s"
}


// dispatcher {
//     receive_queue_cache = 1000
//     send_queue_cache = 1000
    
// }

mongo  {
        server {
            addr = "mongodb"
            port = 27017
        }
        database = "test"
        collection = "distributedLock"
        timeout = "10s"
    }

poller {
    message_cache = 1000
    workers = 10
    snmp "v2c" {
        community = "semipbulic"
        port = 161
        timeout =  "10s"
        retries = 3
        max_iterations = 2
        dynamic_repetitions = true
    }
    connection "ssh" {
            username = "ps2"
            password = "ps2kulator"
            port = 22
            gi = "screen-length disable"
            default_prompt = "\\n[<\\[](\\S+)[>#\\]]$"
            default_errors = "error|failed|unrecognized"
            cache_ttl = "60s"
            read_dead_line = "8s"
            write_dead_line = "8s"
            ssh_key_path = "./keys/"
    }
    connection "telnet" {
            username = "ps2"
            password = "ps2kulator"
            port = 23
            screen_length = "screen-length disable" 
            default_prompt = "\\n[<\\[](\\S+)[>#\\]]$"
            default_errors = "error|failed|unrecognized"
            cache_ttl = "60s"
            read_dead_line = "8s"
            write_dead_line = "8s"
    }
}
