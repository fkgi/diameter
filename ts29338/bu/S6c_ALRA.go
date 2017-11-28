/*
AlertServiceCentreRequest is ALR message.
 <ALR> ::= < Diameter Header: 8388648, REQ, PXY, 16777312 >
           < Session-Id >
           [ DRMP ]
           [ Vendor-Specific-Application-Id ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
           [ Destination-Host ]
           { Destination-Realm }
           { SC-Address }
           { User-Identifier }
           [ SMSMI-Correlation-ID ]
           [ Maximum-UE-Availability-Time ]
           [ SMS-GMSC-Alert-Event ]
           [ Serving-Node ]
         * [ Supported-Features ]
         * [ AVP ]
         * [ Proxy-Info ]
         * [ Route-Record ]
*/
/*
 AlertServiceCentreAnswer is ALA message.
 <ALA> ::= < Diameter Header: 8388648, PXY, 16777312 >
           < Session-Id >
           [ DRMP ]
           [ Vendor-Specific-Application-Id ]
           [ Result-Code ]
           [ Experimental-Result ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
         * [ Supported-Features ]
         * [ AVP ]
         * [ Failed-AVP ]
         * [ Proxy-Info ]
         * [ Route-Record ]
*/

