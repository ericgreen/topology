#!/bin/bash
curl -i -G -H "Accept: application/json" -H "X-Spirent-Chassis-Session-ID: $2" http://$1:9090/chassisController/equipment/timestamp
