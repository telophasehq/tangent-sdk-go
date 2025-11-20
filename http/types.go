package http

type RemoteMethod int

const (
	RemoteMethodGet RemoteMethod = iota
	RemoteMethodPost
	RemoteMethodPut
	RemoteMethodDelete
	RemoteMethodPatch
)

// RemoteHeader is a simple name/value header.
type RemoteHeader struct {
	Name  string
	Value string
}

type RemoteRequest struct {
	ID         string
	Method     RemoteMethod
	URL        string
	Headers    []RemoteHeader
	Body       []byte
	TimeoutMs  *uint32 // nil => no explicit timeout
	CacheTtlMs *uint32 // hint; host may ignore
}

// RemoteResponse mirrors tangent:remote@0.1.0.response.
type RemoteResponse struct {
	ID      string
	Status  uint16
	Headers []RemoteHeader
	Body    []byte
	Error   *string
}
