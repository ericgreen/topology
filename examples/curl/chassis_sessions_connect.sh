#!/bin/bash
curl -i -X POST -H "Accept: application/json" http://$1:9090/chassisController/sessions/connect
