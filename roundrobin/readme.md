# Round-Robin Diameter debugger
Round-Robin can accept or connect to specific Diameter peer node and connect Diameter connection.

Round-Robin can send any Diameter request to peer node.
It make Diameter request message from received HTTP REST request, then send the message to peer node.
It receive Diameter answer message, then make HTTP REST answer from received Diameter answer message.

Round-Robin can receive any Diameter request from peer.
It make HTTP REST request from received Diameter request message, then send the request to pre-configured HTTP server.
It receive HTTP REST answer, then make Diameter answer message from received HTTP REST answer.

HTTP REST request/answer must have specific format JSON document. Dictionary file for converting Diameter message and JSON document must pre-configured.

# How to run Round-Robin
Commandline options.

```
roundrobin [OPTION]... DIAMETER_PEER
DIAMETER_PEER = [(tcp|sctp)://][realm/]hostname[:port]
```

Commandline example

```
roundrobin -l mme.epc.mcc99.mnc999.3gppnetwork.org -i :8080 -b mockserver:8080 -d ./s6a.xml sctp://hss.ecp.mcc99.mnc999.3gppnetwork.org
```

## Args
- `DIAMETER_PEER`  
Diameter peer host definition.
Round-Robin connect to specified Diameter peer.

## Options
- `-l`  
Diameter local host definition.
Value must have format `[realm/]host[:port]`.
Hostname of operating system is used as default of `host`.
Refer following section about other parameters.

- `-i`  
Local listening address and port for receiving HTTP REST request.
Value must have format `host[:port]`.
`host` is hostname or IP address.
IP address is resolved from hostname if hostname is specified.
`port` is port number.

- `-b`  
Peer address and port for sending HTTP REST request.
Value must have format `host[:port]`.
`host` is hostname or IP address.
IP address is resolved from hostname if hostname is specified.
`port` is port number.

- `-d`  
Path for dictionary XML file.
`dictionary.xml` file in current directory is used as default.

- `-t`  
Duration of Diameter request timeout in second.
Not only service message but also control message like DWR follow this duration.

## Format of Diameter node identity

```
[(tcp|sctp)://][realm/]hostname[:port]
```
Transport layer protocol is specified by first item `[(tcp|sctp)://]`.
`tcp` and `sctp` is acceptable value. `tcp` is used as default if this item is omitted.

Diameter Realm is specified by second item `[realm/]`.
The text must has Diameter-Identity format that is defined by RFC 6733.
It is generated from `hostname` if the item is empty.

Diameter Host is specified by third item `hostname`.
The text must has Diameter-Identity format that is defined by RFC 6733.
IP address is derived from hostname.

Transport layer port number is specified by last item `[:port]`.
It must available port number digit.
`3868` is used as default if this item is omitted.

Port `0` is used as any for source port.
If local port is 0, local port is automaticaly selected by system.

# Format of Dictionary file
Dictionary file is XML document.

```
<dictionary>
    <vendor name="3GPP" id="10415">
        <application name="S6a" id="16777251">
            <command name="Update-Location" id="316" />
            <command name="Cancel-Location" id="317" />
        </application>
        <avp name="Subscription-Data" id="1400" type="Grouped" mandatory="true" />
        <avp name="Terminal-Information" id="1401" type="Grouped" mandatory="true" />
    </vendor>
    <vendor name="IETF" id="0">
        <application name="base" id="0">
            <command name="Capabilities-Exchange" id="257" />
        </application>
        <avp name="Acct-Interim-Interval" id="85" type="Unsigned32" mandatory="true" />
    </vendor>
</dictionary>
```
Root element is `<dictionary>`.

Second element is `<vendor>` that include Vendor layer data.
Vendor layer data defines vendor.
It has `name` and `id` attributes, `application` and `avps` element data.
- `name` in Vendor layer is name of the vendor
- `id` in Vendor layer is vendor ID digit that is assigned in IANA
- `application` contains elements of Application layer data
- `avps` contains elements of AVP layer data

Application layer data defines Diameter application for specified vendor.
It has `name` and `id` attributes, `command` element data.
- `name` in Application layer is name of the application
- `id` in Application layer is application ID digit that is assigned in IANA
- `command` contains elements of Command layer data

Command layer data defines Diameter command in specified application.
It has `name` and `id` attributes.
- `name` in Command layer is name of the command
- `id` in Command layer is command code digit that is assigned in IANA

AVP layer data defines Diameter AVP for specified vendor.
It has `name`, `id`, `mandatory`, `type` attributes, `enum` element data.
- `name` in AVP layer is name of the avp
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
- `enum` define Enumerated value mapping  
`value` attribute is Enumerated digit value  
Text data of the element is Enumerated value name  

# Format of REST message
Only POST method is acceptable for HTTP REST request.
HTTP URI path has prefix `/diamsg/v1`.
HTTP URI path has required Diameter command by `/{vendor name}/{application name}/{command name}`. Each names are defined in dictionary file.
```
POST http://roundrobin:8080/diamsg/v1/3GPP/S6a/Update-Location
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

If there are multiple AVPs that have same AVP ID, array is used for Value of the MAP.
Each elements of array are value of independent AVP that has same AVP ID.

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
    "Homogeneous-Support-of-IMS-Voice-Over-PS-Sessions": "NOT_SUPPORTED",
    "Supported-Features": [
        {  
            "Vendor-Id": 10415,
            "Feature-List-ID": 1,
            "Feature-List": 2
        },
        {  
            "Vendor-Id": 10415,
            "Feature-List-ID": 2,
            "Feature-List": 1
        }
    ]
}
```

# Behavior for specific AVP
## Session-ID
If Session-ID AVP is exist but the value is empty, Round-Robbin generate session ID automatically and fill in to empty Session-ID.
New session ID is generated for Session-ID AVP in request message.
Session ID is copied from request message for Session-ID AVP in response message.
```
{
    ...
    "Session-Id": "",
    ...
}
```
## Origin-Host/Realm
If Origin-Host AVP is exist but the value is empty, Round-Robbin set own local Diameter Host to it.
If Origin-Realm AVP is exist but the value is empty, Round-Robbin set own local Diameter Realm to it.
```
{
    ...
    "Origin-Host": "",
    "Origin-Realm": "",
    ...
}
```

# Behavior for specific HTTP result code
## 200 OK
If peer HTTP server returns 200 OK response with correct JSON data, Round-Robbin make Diameter response.
Result code of Diameter response is in part of JSON data.

## 503 Service Unavailable
If peer HTTP server returns 503 Service Unabailable response, Round-Robbin discard Diameter transaction and does not response to originator Diameter peer.

## Others
If peer HTTP server returns any other response, Round-Robbin make Diameter response with 5012 UNABLE_TO_COMPLY response code..

