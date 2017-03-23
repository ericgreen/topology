#!/bin/bash
curl -i -X POST -H "Accept: application/json" http://$1:9192/topology/discovery/discover
