package http

type RemoteMethod int

const (
	RemoteMethodGet RemoteMethod = iota
	RemoteMethodPost
	RemoteMethodPut
	RemoteMethodDelete
	RemoteMethodPatch
)

// Header is a simple name/value header.
type Header struct {
	Name  string
	Value string
}

type Request struct {
	ID         string
	Method     RemoteMethod
	URL        string
	Headers    []Header
	Body       []byte
	TimeoutMs  *uint32 // nil => no explicit timeout
	CacheTtlMs *uint32 // hint; host may ignore
}

// Response mirrors tangent:remote@0.1.0.response.
type Response struct {
	ID      string
	Status  uint16
	Headers []Header
	Body    []byte
	Error   *string
}
