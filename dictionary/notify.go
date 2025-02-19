package dictionary

import (
	"fmt"
	"strings"

	"github.com/fkgi/diameter"
)

func TraceMessageVarbose(prefix string, msg diameter.Message) string {
	buf := new(strings.Builder)

	if com, e := DecodeMessage(msg); e == nil {
		fmt.Fprintf(buf, "%s%s", prefix, com)
	} else {
		fmt.Fprintf(buf,
			"%sunknown_command(appID=%d, code=%d)",
			prefix, msg.AppID, msg.Code)
	}
	if msg.FlgR {
		fmt.Fprint(buf, "-Request")
	} else {
		fmt.Fprint(buf, "-Answer")
	}

	fmt.Fprint(buf, " (")
	if msg.FlgP {
		fmt.Fprint(buf, "P")
	} else {
		fmt.Fprint(buf, "-")
	}
	if msg.FlgE {
		fmt.Fprint(buf, "E")
	} else {
		fmt.Fprint(buf, "-")
	}
	if msg.FlgT {
		fmt.Fprint(buf, "T")
	} else {
		fmt.Fprint(buf, "-")
	}
	fmt.Fprintln(buf, ")")

	fmt.Fprintf(buf,
		"%sHop-by-Hop-ID=%#x, End-to-End-ID=%#x",
		prefix, msg.HbHID, msg.EtEID)
	fmt.Fprintln(buf)

	avps, e := msg.GetAVP()
	if e != nil {
		fmt.Fprintf(buf, "%s%sHEX body=% x", prefix, prefix, msg.AVPs)
		fmt.Fprintln(buf)
	} else {
		for _, a := range avps {
			n, v, e := DecodeAVP(a)
			if e != nil {
				fmt.Fprintf(buf,
					"%s%sunknown AVP(vendorID=%d, code=%d): % x",
					prefix, prefix, a.VendorID, a.Code, a.Data)
				fmt.Fprintln(buf)
			} else {
				printAVP(prefix, n, 2, v, buf)
			}
		}
	}

	return buf.String()
}

func printAVP(prefix, name string, depth int, value any, buf *strings.Builder) {
	gr, ok := value.(map[string]any)
	p := ""
	for i := 0; i < depth; i++ {
		p += prefix
	}
	if ok {
		fmt.Fprintf(buf, "%s%s:", p, name)
		fmt.Fprintln(buf)
		for k, v := range gr {
			printAVP(prefix, k, depth+1, v, buf)
		}
	} else {
		fmt.Fprintf(buf, "%s%s: %v", p, name, value)
		fmt.Fprintln(buf)
	}
}

var NotifyHandlerError func(proto, msg string) = nil
