package connector_test

import (
	"testing"

	"github.com/fkgi/diameter/connector"
)

func TestParse(t *testing.T) {
	parse("sctp://realm.test.com/host.realm.test.com:123", t)
	parse("tcp://hogehoge/google.com:0000", t)
	parse("udp://hogehoge/google.com:3868", t)
	parse("hogehoge/hogehoge.mogemoge", t)
	parse("mogemoge.mogemoge/yahoo.com:9999", t)
	parse("google.co.jp", t)
	parse("yahoo.com:1234", t)
	parse("sctp://google.co.jp", t)
	parse("tcp://google.co.jp:3969", t)
}
func parse(uri string, t *testing.T) {
	scheme, host, realm, ips, port, err := connector.ResolveIdentiry(uri)
	t.Logf("original=%s", uri)
	t.Logf("scheme=%s, host=%s, realm=%s, ip=%v, port=%d, error=%v",
		scheme, host, realm, ips, port, err)
}
