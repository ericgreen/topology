#!/bin/bash
curl -i -G -H "Accept: application/json" http://$1:9192/topology/discovery/instances/$2/$3
