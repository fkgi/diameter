
/*
ReportSMDeliveryStatusRequest is RDR message.
 <RDR> ::= < Diameter Header: 8388649, REQ, PXY, 16777312 >
           < Session-Id >
           [ DRMP ]
           [ Vendor-Specific-Application-Id ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
           [ Destination-Host ]
           { Destination-Realm }
         * [ Supported-Features ]
           { User-Identifier }
           [ SMSMI-Correlation-ID ]
           { SC-Address }
           { SM-Delivery-Outcome }
           [ RDR-Flags ]
         * [ AVP ]
         * [ Proxy-Info ]
         * [ Route-Record ]
*/
/*
ReportSMDeliveryStatusAnswer is RDA message.
 <RDA> ::= < Diameter Header: 8388649, PXY, 16777312 >
           < Session-Id >
           [ DRMP ]
           [ Vendor-Specific-Application-Id ]
           [ Result-Code ]
           [ Experimental-Result ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
         * [ Supported-Features ]
           [ User-Identifier ]
         * [ AVP ]
         * [ Failed-AVP ]
         * [ Proxy-Info ]
         * [ Route-Record ]
*/
