package eudore


import (
	"io"
	"bytes"
	"io/ioutil"
	"net/http"
	"crypto/tls"
)

type (
	// Get the method, version, uri, header, body from the RequestReader according to the http protocol request body. (There is no host in the golang net/http library header)
	//
	// Read the remote connection address and TLS information from the net.Conn connection.
	//
	// 根据http协议请求体，从RequestReader获取方法、版本、uri、header、body。(golang net/http库header中没有host)
	//
	// 从net.Conn连接读取远程连接地址和TLS信息。
	RequestReader interface {
		// http protocol data
		Method() string
		Proto() string
		RequestURI() string
		Header() Header
		Read([]byte) (int, error)
		Host() string
		// conn data
		RemoteAddr() string
		TLS() *tls.ConnectionState
	}
	RequestReadSeeker interface {
		RequestReader
		io.Seeker
	}
	// Convert RequestReader to the net.http.Request object interface.
	//
	// 将RequestReader转换成net.http.Request对象接口。
	RequestConvertNetHttp interface {
		GetNetHttpRequest() *http.Request
	}
	// Convert net.http.Request to RequestReader.
	//
	// 将net/http.Request转换成RequestReader。
	RequestReaderHttp struct {
		http.Request
	}
	// Modify the RequestReader method and request uri inside the internal redirect.
	//
	// 内部重定向内修改RequestReader的方法和请求uri。
	RequestReaderRedirect struct {
		RequestReader
		method string
		uri	string
	}
	RequestReaderSeeker struct {
		RequestReader
		reader *bytes.Reader
	}
	RequestReaderEudore struct {
		method string
		proto string
		requestURI string
		remoteAddr string
		header http.Header
		body []byte
		tls *tls.ConnectionState
	}
)


var _ RequestReader		=	&RequestReaderHttp{}

func NewRequestReaderHttp(r *http.Request) RequestReader {
	return &RequestReaderHttp{
		Request:	*r,
	}
}

func ResetRequestReaderHttp(r *RequestReaderHttp, req *http.Request) RequestReader {
	r.Request = *req
	return r
}

func (r *RequestReaderHttp) Read(p []byte) (int, error) {
	return r.Request.Body.Read(p)
}

func (r *RequestReaderHttp) Method() string {
	return r.Request.Method 
} 

func (r *RequestReaderHttp) Proto() string {
	return r.Request.Proto
}

func (r *RequestReaderHttp) Host() string {
	return r.Request.Host	
}

func (r *RequestReaderHttp) RequestURI() string {
	return r.Request.RequestURI
}

func (r *RequestReaderHttp) Header() Header {
	return r.Request.Header
} 

func (r *RequestReaderHttp) RemoteAddr() string {
	return r.Request.RemoteAddr
}

func (r *RequestReaderHttp) TLS() *tls.ConnectionState {
	return r.Request.TLS
}

func (r *RequestReaderHttp) GetNetHttpRequest() *http.Request {
	return &r.Request
}


func NewRequestReaderRedirect(r RequestReader, method, uri string) (RequestReader) {
	return &RequestReaderRedirect{
		RequestReader:	r,
		method:			method,
		uri:			uri,
	}
}

func (r *RequestReaderRedirect) Method() string {
	return r.method
}

func (r *RequestReaderRedirect) RemoteAddr() string {
	return r.uri
}

func NewRequestReaderSeeker(r RequestReader) (RequestReadSeeker) {
	rs, ok := r.(RequestReadSeeker)
	if ok {
		return rs
	}
	bts, _ := ioutil.ReadAll(r)
	return &RequestReaderSeeker{
		RequestReader:	r,
		reader:			bytes.NewReader(bts),
	}
}

func (r *RequestReaderSeeker) Read(p []byte) (int, error) {
	return r.reader.Read(p)
}

func (r *RequestReaderSeeker) Seek(offset int64, whence int) (int64, error) {
	return r.reader.Seek(offset, whence)
}