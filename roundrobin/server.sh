#! /bin/bash

echo -e 'HTTP/1.1 200 OK\r
Content-Type: application/json\r
Connection: close\r
\r
{
  "Session-Id": "",
  "Auth-Session-State": "NO_STATE_MAINTAINED",
  "Origin-Host": "hss.dummy.net",
  "Origin-Realm": "epc.mnc99.mcc999.3gppnetwork.org",
  "Result-Code": 3002,
  "Vendor-Specific-Application-Id": {
    "Auth-Application-Id": 16777251,
    "Vendor-Id": 10415
  }
}'| nc -l 8081
