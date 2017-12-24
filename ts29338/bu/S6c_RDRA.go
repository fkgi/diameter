package ts29338

import (
	dia "github.com/fkgi/diameter"
	"github.com/fkgi/teldata"
)

/*
RDR is Report-SM-Delivery-Status-Request message.
 <RDR> ::= < Diameter Header: 8388649, REQ, PXY, 16777312 >
           < Session-Id >
           [ DRMP ] // not supported
           [ Vendor-Specific-Application-Id ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
           [ Destination-Host ]
           { Destination-Realm }
         * [ Supported-Features ] // not supported
           { User-Identifier }
           [ SMSMI-Correlation-ID ] // not supported
           { SC-Address }
           { SM-Delivery-Outcome }
           [ RDR-Flags ]
         * [ AVP ]
         * [ Proxy-Info ] // not supported
         * [ Route-Record ]
*/
type RDR struct {
	OriginHost       dia.Identity
	OriginRealm      dia.Identity
	DestinationHost  dia.Identity
	DestinationRealm dia.Identity

	MSISDN teldata.E164
	teldata.IMSI
	SCAddress teldata.E164

	Flags struct {
		SingleAttempt bool
	}

	// DRMP
	// SMSMICorrelationID
	// []SupportedFeatures
	// []ProxyInfo
}

/*
ReportSMDeliveryStatusAnswer is RDA message.
 <RDA> ::= < Diameter Header: 8388649, PXY, 16777312 >
           < Session-Id >
           [ DRMP ] // not supported
           [ Vendor-Specific-Application-Id ]
           [ Result-Code ]
           [ Experimental-Result ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
         * [ Supported-Features ] // not supported
           [ User-Identifier ]
         * [ AVP ]
         * [ Failed-AVP ]
         * [ Proxy-Info ] // not supported
         * [ Route-Record ]
*/
