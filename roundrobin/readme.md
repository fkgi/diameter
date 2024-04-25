# Round-Robin Diameter debugger
Round-Robin can accept or connect to specific Diameter peer node and connect Diameter connection.

Round-Robin can send any Diameter request to peer node as client.
It make Diameter request message from received HTTP REST request, then send the message to peer node.
It receive Diameter answer message, then make HTTP REST answer from received Diameter answer message.

Round-Robin can receive any Diameter request from peer node as server.
It make HTTP REST request from received Diameter request message, then send the request to pre-configured HTTP server.
It receive HTTP REST answer, then make Diameter answer message from received HTTP REST answer.

HTTP REST request/answer must have specific format JSON document. Dictionary file for converting Diameter message and JSON document must pre-configured.

# How to run Round-Robin
Execute binary file with below options.

```
rountdobin -diameter-local sctp://mme.epc.mcc99.mnc999.3gppnetwork.org -diameter-peer sctp://hss.ecp.mcc99.mnc999.3gppnetwork.org -http-local :8080 -http-peer mockserver:8080 -dictionary ./s6a.json
```

## Options
- diameter-local  
Diameter local host definition.
Hostname of operating system is used as default.

- diameter-peer  
Diameter peer host definition.
If this parameter is used, Round-Robin run as client and try to connect to specified peer node.

- http-local  
Local listening address and port for receiving HTTP REST request.
Value must have format `host[:port]`.
`host` is hostname or IP address.
IP address is resolved from hostname if hostname is specified.
`port` is port number.

- http-peer  
Peer connecting address and port for sending HTTP REST request.
Value must have format `host[:port]`.
`host` is hostname or IP address.
IP address is resolved from hostname if hostname is specified.
`port` is port number.

- dictionary  
Path for dictionary JSON file.
`dictionary.json` file in current directory is used as default.

## Format of Diameter node identity
`diameter-local` and `diameter-peer` option must have below format.
```
[tcp|sctp://][realm/]hostname[:port]
```
Transport layer protocol is specified by first item `[tcp|sctp://]`.
`tcp` and `sctp` is acceptable value. `tcp` is used as default if this item is omitted.

Diameter Realm is specified by second item `[realm/]`.
The text must has Diameter-Identity format that is defined by RFC 6733.
It is generated from `hostname` if the item is empty.

Diameter Host is specified by third item `hostname`.
The text must has Diameter-Identity format that is defined by RFC 6733.
IP address is derived from hostname.

Transport layer port number is specified by last item `[:port]`.
It must available port number digit.
3868 is used as default if this item is omitted.

# Format of Dictionary file
Dictionary file is JSON document.

```
{
    "3GPP": {
        "id": 10415,
        "applications": {
            "S6a": {
                "id": 16777251,
                "command": {
                    "Update-Location": {
                        "id": 316
                    }
                }
            }
        },
        "avps": {
            "Subscription-Data": {
                "id": 1400,
                "mandatory": true,
                "type": "Grouped"
            },
            "Terminal-Information": {
                "id": 1401,
                "mandatory": true,
                "type": "Grouped"
            }
        }
    }
}
```
Root is JSON Map object.
Key of the Map is vendor name.
Value of the Map is JSON object for Vendor layer data.

Vendor layer data defines vendor.
It has `id`, `applications`, `avps` data.
- `id` in Vendor layer is vendor ID digit that is assigned in IANA
- `applications` contains JSON Map object of Application layer data
- `avps` contains JSON Map object of AVP layer data

Application layer data defines Diameter application for specified vendor.
Key of the Map is application name.
Value of the Map is JSON object for Application layer data.
It has `id`, `command` data.
- `id` in Application layer is application ID digit that is assigned in IANA
- `command` contains JSON Map object of Command layer data

Command layer data defines Diameter command in specified application.
Key of the Map is command name.
Value of the Map is JSON object for Command layer data.
It has `id` data.
- `id` in Command layer is command code digit that is assigned in IANA

AVP layer data defines Diameter AVP for specified vendor.
Key of the Map is AVP name.
Value of the Map is JSON object for AVP layer data.
It has `id`, `mandatory`, `type` data.
- `id` in AVP layer is AVP ID digit that is assigned in IANA
- `mandatory` is boolean data that indicate the AVP is flagged as mandatory
- `type` is string data that indicate format of the AVP  
Available format are below
  - OctetString
  - Integer32
  - Integer64
  - Unsigned32
  - Unsigned64
  - Float32
  - Float64
  - Grouped
  - Address
  - Time
  - UTF8String
  - DiameterIdentity
  - DiameterURI
  - Enumerated
  - IPFilterRule  
- `map` is JSON Map object for defining Enumerated value  
Key of the Map is Enumerated value name  
Value of the Map is Enumerated digit value

# Format of REST message
Only POST method is acceptable for HTTP REST request.
HTTP URI path has prefix `/msg/v1`.
HTTP URI path has required Diameter command by `/{vendor name}/{application name}/{command name}`. Each names are defined in dictionary file.
```
POST http://roundrobin:8080/msg/v1/3GPP/S6a/Update-Location
```

HTTP body is JSON Map object.
Key of the Map is name of AVP that is defined in dictionary.
Value of the Map is AVP value.
Value is string, number or nested JSON Map object. The value format is defined as below based on dictionary.
  - OctetString : string that hex formatted binary data
  - Integer32 : number
  - Integer64 : number
  - Unsigned32 : number
  - Unsigned64 : number
  - Float32 : number
  - Float64 : number
  - Grouped : JSON Map object
  - Address : string with IP address format
  - Time : string with RFC 3339 time format
  - UTF8String : string
  - DiameterIdentity : string with RFC 6733 Diameter-Identity format
  - DiameterURI : string with RFC 6733 Diameter-URI format
  - Enumerated : string that is defined in dictionary
  - IPFilterRule  : string with RFC 6733 IP-Filter-Rule format

```
{
    "Session-Id": "mme.ecp.mcc99.mnc999.3gppnetwork.org;12345",
    "Vendor-Specific-Application-Id": {
        "Vendor-Id": 10415,
        "Auth-Application-Id": 16777251
    },
    "Auth-Session-State": "NO_STATE_MAINTAINED",
    "Origin-Host": "mme.ecp.mcc99.mnc999.3gppnetwork.org",
    "Origin-Realm": "ecp.mcc99.mnc999.3gppnetwork.org",
    "Destination-Realm": "ecp.mcc99.mnc999.3gppnetwork.org",
    "User-Name": "999990123456789",
    "RAT-Type": "EUTRAN",
    "ULR-Flags": 98,
    "Visited-PLMN-Id": "99F999",
    "Terminal-Information": {
        "Software-Version": "03",
        "IMEI": "01234567890123"
    },
    "UE-SRVCC-Capability": "UE-SRVCC-SUPPORTED",
    "Homogeneous-Support-of-IMS-Voice-Over-PS-Sessions": "NOT_SUPPORTED"
}
```
