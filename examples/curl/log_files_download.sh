#!/bin/bash
curl -i -G -H "Accept: application/octet-stream" -H "X-Spirent-Chassis-Session-ID: $2" -o $3 http://$1:9090/chassisController/logFiles/$3
cat logfile_sysmgrd.log | more
