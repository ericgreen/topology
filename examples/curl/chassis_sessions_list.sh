#!/bin/bash
curl -i -G -H "Accept: application/json" http://$1:9090/chassisController/sessions
