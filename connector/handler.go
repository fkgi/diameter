package connector

import "github.com/fkgi/diameter"

// Handle wrap diameter.Handle for single connection.
func Handle(code, appID, venID uint32, h diameter.Handler) diameter.Handler {
	return diameter.Handle(
		code, appID, venID, h,
		func() *diameter.Connection { return &con })
}

// DefaultTxHandler wrap diameter.DefaultTxHandler for single connection
func DefaultTxHandler(m diameter.Message) diameter.Message {
	return con.DefaultTxHandler(m)
}
