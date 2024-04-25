#! /bin/bash

echo -e 'HTTP/1.1 200 OK\r
Content-Type: application/json\r
Connection: close\r
\r
{
  "Session-Id": "mme.epc.mnc99.mcc999.3gppnetwork.org;2023-09-05T22:33:43+09:00",
  "Auth-Session-State": "NO_STATE_MAINTAINED",
  "Origin-Host": "hss.dummy.net",
  "Origin-Realm": "epc.mnc99.mcc999.3gppnetwork.org",
  "Result-Code": 3002,
  "Vendor-Specific-Application-Id": {
    "Auth-Application-Id": 16777251,
    "Vendor-Id": 10415
  }
}'| nc -l 8081
