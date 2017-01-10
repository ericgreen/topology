#!/bin/bash
curl -i -X DELETE -H "Accept: application/json" -H "X-Spirent-Chassis-Session-ID: $2" http://$1:9090/chassisController/testPackages/$3
