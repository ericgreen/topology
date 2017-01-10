#!/bin/bash
curl -i -X POST -H "Accept: application/json" -H "X-Spirent-Chassis-Session-ID: $2" http://$1:9090/chassisController/equipment/activateProfile/$3/$4
