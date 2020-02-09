import http from "k6/http";

export default function() {
  var url = "http://localhost:1337/v1/ti/";
    var payload = JSON.stringify({ hostname: "seume-ss1.ne.liero.net", port: "gigabitethernet0/0/1" });
	  var params =  { headers: { "Content-Type": "application/json" } }
	    http.post(url, payload, params);
		};
