#!/bin/bash
curl -i -X POST -H "Content-Type: application/json" -H "Accept: application/json" -d@hypervisors_instances_parameters.json http://$1:9192/topology/cloudHypervisorsOvsNetworkTopology/$2