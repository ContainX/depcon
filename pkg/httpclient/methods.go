package httpclient

type Method int

const (
	GET Method = 1 + iota
	POST
	PUT
	DELETE
	HEAD
)

var methods = [...]string{
	"GET",
	"POST",
	"PUT",
	"DELETE",
	"HEAD",
}

func (method Method) String() string {
	return methods[method-1]
}
