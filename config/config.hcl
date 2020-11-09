



mongo {
  addr = localhost
  port = 27017
  database = "swpx"

}

logger {
  level = "DEBUG"
  as_json = false
}




// resource plugins/drivers
resource "vrp" {
  name = "huawei"
  version = "v1.0.0"
  description = "switches from huawei (VRP software)"

}


resource "raycore" {
  name = "raycore"
  version = "v1.0.0"
  description = "cpe from raycore"

}


// providers plugins/drivers
provider "vx" {
  name = "vx"
  version = "v1.0.0"
  description = "provider from VX to do lookups and other stuff from"

}