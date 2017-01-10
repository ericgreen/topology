#!/bin/bash
curl -i -X POST -H "Accept: application/json" -H "X-Spirent-Chassis-Session-ID: $2" -d@reserve_target_parameters.json http://$1:9090/chassisController/equipment/reserve
