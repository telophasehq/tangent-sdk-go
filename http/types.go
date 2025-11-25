package http

type Method int

const (
	MethodGet Method = iota
	MethodPost
	MethodPut
	MethodDelete
	MethodPatch
)

// Header is a simple name/value header.
type Header struct {
	Name  string
	Value string
}

type Request struct {
	ID         string
	Method     Method
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
