#!/bin/bash
curl -i -G -H "Accept: application/json" http://$1:9192/topology/discovery/hypervisors/$2
