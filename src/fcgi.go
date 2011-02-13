package web

import (
        "bufio"
        "bytes"
        "encoding/binary"
        "fmt"
        "http"
        "io"
        //"io/ioutil"
        "log"
        "os"
        "runtime"
        "strconv"
        //"strings"
        "syscall"
)

type FCGIModel struct {
        isCGI bool
}

// TODO: fix this '= 0'
var flagCGI int = 2 // 0 = Unchecked, 1 = CGI, 2 = FCGI

var logger *log.Logger

func NewFCGIModel() (am AppModel, err os.Error) {
        if flagCGI == 0 {
                sa, ec := syscall.Getpeername(FCGI_LISTENSOCK_FILENO)
                //fmt.Printf("CGI: %v, %x, %v\n", sa, ec, os.Errno(ec))
                if sa == nil && (ec == syscall.ENOTSOCK || ec == syscall.ENOTCONN) {
                        flagCGI = 1
                } else {
                        flagCGI = 2
                }
        }

        if flagCGI == 1 {
                am, err = NewCGIModel()
        } else {
                m := &FCGIModel{}
                am = AppModel(m)
        }
        return
}

type VersionCode uint8
type RequestId uint16

const (
        FCGI_LISTENSOCK_FILENO = 0
        FCGI_VERSION_1 VersionCode = 1
        FCGI_HEADER_LEN = 8
        FCGI_MAX_LENGTH = 0xffff
        FCGI_NULL_REQUEST_ID RequestId = 0
)

type RecordType uint8

const (
        FCGI_BEGIN_REQUEST RecordType = 1 + iota
        FCGI_ABORT_REQUEST
        FCGI_END_REQUEST
        FCGI_PARAMS
        FCGI_STDIN
        FCGI_STDOUT
        FCGI_STDERR
        FCGI_DATA
        FCGI_GET_VALUES
        FCGI_GET_VALUES_RESULT
        FCGI_UNKNOWN_TYPE
        FCGI_MAXTYPE = FCGI_UNKNOWN_TYPE
)

type RecordHeader struct { // see fastcgi.h
        Version VersionCode
        Type RecordType
        RequestId RequestId
        ContentLength uint16
        PaddingLength uint8
        reserved uint8
}

type Role uint16

const (
        FCGI_RESPONDER Role = 1 + iota
        FCGI_AUTHORIZER
        FCGI_FILTER
)

type BeginFlag uint8

const (
        FCGI_KEEP_CONN BeginFlag = 1
)

type BeginRequest struct {
        header *RecordHeader
        body BeginRequestBody
}

type BeginRequestBody struct {
        Role Role
        Flags BeginFlag
        reserved [5]uint8
}

type ProtocolStatus uint8

const ( // Protocol Status
        FCGI_REQUEST_COMPLETE ProtocolStatus = iota
        FCGI_CANT_MPX_CONN
        FCGI_OVERLOADED
        FCGI_UNKNOWN_ROLE
)

type EndRequest struct {
        header *RecordHeader
        body EndRequestBody
}

type EndRequestBody struct {
        AppStatus uint32
        ProtocolStatus ProtocolStatus
        reserved [3]uint8
}

type UnknownTypeRecord struct {
        header *RecordHeader
        body UnknownTypeBody
}

type UnknownTypeBody struct {
        Type uint8
        reserved [7]uint8
}

func (v VersionCode) String() string {
        return fmt.Sprintf("Version(%d)", uint8(v))
}

func (ri RequestId) String() string {
        return fmt.Sprintf("RequestId(%d)", uint16(ri))
}

func (t RecordType) String() string {
        names := []string{
        FCGI_BEGIN_REQUEST: "FCGI_BEGIN_REQUEST",
        FCGI_ABORT_REQUEST: "FCGI_ABORT_REQUEST",
        FCGI_END_REQUEST: "FCGI_END_REQUEST",
        FCGI_PARAMS: "FCGI_PARAMS",
        FCGI_STDIN: "FCGI_STDIN",
        FCGI_STDOUT: "FCGI_STDOUT",
        FCGI_STDERR: "FCGI_STDERR",
        FCGI_DATA: "FCGI_DATA",
        FCGI_GET_VALUES: "FCGI_GET_VALUES",
        FCGI_GET_VALUES_RESULT: "FCGI_GET_VALUES_RESULT",
        FCGI_UNKNOWN_TYPE: "FCGI_UNKNOWN_TYPE",
        }
        if len(names) <= int(t) {
                return fmt.Sprintf("FCGI_UNKNOWN_TYPE<%v>", uint8(t))
        }
        return names[t]
}

func (r Role) String() string {
        names := []string{
        0: "Role<0>",
        FCGI_RESPONDER: "FCGI_RESPONDER",
        FCGI_AUTHORIZER: "FCGI_AUTHORIZER",
        FCGI_FILTER: "FCGI_FILTER",
        }
        if len(names) <= int(r) {
                return fmt.Sprintf("Role<%v>", uint16(r))
        }
        return names[r]
}

func (f BeginFlag) String() string {
        names := []string{
        0: "BeginRequestFlag(0)",
        FCGI_KEEP_CONN: "FCGI_KEEP_CONN",
        }
        if len(names) <= int(f) {
                return fmt.Sprintf("BeginRequestFlag(%v)", uint16(f))
        }
        return names[f]
}

func (s ProtocolStatus) String() string {
        names := []string{
        FCGI_REQUEST_COMPLETE: "FCGI_REQUEST_COMPLETE",
        FCGI_CANT_MPX_CONN: "FCGI_CANT_MPX_CONN",
        FCGI_OVERLOADED: "FCGI_OVERLOADED",
        FCGI_UNKNOWN_ROLE: "FCGI_UNKNOWN_ROLE",
        }
        if len(names) <= int(s) {
                return fmt.Sprintf("ProtocolStatus<%d>", uint8(s))
        }
        return names[s]
}

func (h *RecordHeader) read(r io.Reader) (ok bool, err os.Error) {
        b := make([]byte, 8)
        l, err := io.ReadFull(r, b)
        if l != 8 {
                if err == nil {
                        err = newError(fmt.Sprintf("bad header read: %d", l))
                }
                return
        }

        h.Version = VersionCode(b[0])
        h.Type = RecordType(b[1])
        h.RequestId = RequestId(binary.BigEndian.Uint16(b[2:4]))
        h.ContentLength = binary.BigEndian.Uint16(b[4:6])
        h.PaddingLength = b[6]
        //h.reserved = b[7]
        ok = true

        return
}

func (h *RecordHeader) readContent(r io.Reader) (b []byte, err os.Error) {
        l := int(h.ContentLength) + int(h.PaddingLength)
        b = make([]byte, l)

        n := 0
        for n < l {
                m, err := r.Read(b[n:])
                if 0 < m { n += m }
                if err != nil {
                        b = b[0:n]
                        if err == os.EOF {
                                err = newError(fmt.Sprintf("Short read %d bytes of %d", n, l))
                        }
                        return
                }
        }

        // Discard the padding bytes
        b = b[0:h.ContentLength]
        return
}

func (rec *BeginRequest) parse(b []byte) (ok bool) {
        rec.body.Role = Role(binary.BigEndian.Uint16(b[0:2]))
        rec.body.Flags = BeginFlag(b[2])
        ok = true
        return
}

func (rec *EndRequest) parse(b []byte) (ok bool) {
        rec.body.AppStatus = binary.BigEndian.Uint32(b[0:4])
        rec.body.ProtocolStatus = ProtocolStatus(b[4])
        ok = true
        return
}

func (rec *UnknownTypeRecord) parse(b []byte) (ok bool) {
        rec.body.Type = b[0]
        ok = true
        return
}

func getNVSize(b []byte) (size int, ob []byte) {
        s := b[0]
        if s>>7 == 1 {
                if 4 <= len(b) {
                        b0 := b[0]
                        b[0] = b[0] & 0x7F // clear the FLAG bit
                        size = int(binary.BigEndian.Uint32(b[0:4]))
                        b[0] = b0 // restore B0
                        ob = b[4:]
                }
        } else {
                if 0 < len(b) {
                        size = int(s)
                        ob = b[1:]
                }
        }
        return
}

func getNVValue(b []byte, size int) (v string, ob []byte) {
        if len(b) < size {
                v = string(b[0:])
                ob = b[0:0]
        } else {
                v = string(b[0:size])
                ob = b[size:]
        }
        return
}

func parseNameValuePares(b []byte) (params map [string]string) {
        var kl, vl int
        var k, v string
        params = make(map[string]string)
        for 0 < len(b) {
                if kl, b = getNVSize(b);     b == nil { break }
                if vl, b = getNVSize(b);     b == nil { break }
                if kl < 0 || vl < 0 {
                        logger.Printf("error: wrong Name-Value pair size: klen=%d, vlen=%d, (%d bytes left)\n", kl, vl, len(b))
                        logger.Printf("error: Name-Value pares: %v", string(b))
                        break
                }
                if k, b = getNVValue(b, kl); b == nil { break }
                if v, b = getNVValue(b, vl); b == nil { break }
                //logger.Printf("kv: %v, %v, %d\n", k, v, len(b))
                params[k] = v
        }
        return
}

func printErrorCallStack(err interface{}, str io.Writer) {
        stack := make([]uintptr, 30)

        n := runtime.Callers(0, stack)
        stack = stack[0:n]

        fmt.Fprintf(str, "Content-Type: text/html; charset=utf-8\n\n")
        fmt.Fprintf(str, `<font color="red"><b>%v</b></font><p>`, err)
        for i := range stack {
                pc := stack[len(stack)-i-1]
                f := runtime.FuncForPC(pc)
                if f != nil {
                        file, line := f.FileLine(pc)
                        fmt.Fprintf(str, `%s:%d: <font color="red">%s</font><br/>`, file, line, f.Name())
                }
        }
        fmt.Fprintf(str, `</p>`)
}

func initRequest(request *Request, params map[string]string) (err os.Error) {
        request.Header = params

        request.Method = params["REQUEST_METHOD"]
        request.Proto = params["SERVER_PROTOCOL"]

        ok := false
        if request.Proto == "" {
                logger.Printf("empty SERVER_PROTOCOL, force to HTTP/1.1")
                request.Proto = "HTTP/1.1"
        }

        request.ProtoMajor, request.ProtoMinor, ok = parseHTTPVersion(request.Proto)
        if !ok {
                err = newError("malformed HTTP version: '"+request.Proto+"'")
                return
        }

        request.RawURL = params["REQUEST_URI"]
        request.URL, err = http.ParseURL(request.RawURL)
        if err != nil {
                logger.Printf("error: REQUEST_URI: '%v' (%v)", request.RawURL, err)
                //return
        }

        request.Host = params["HTTP_HOST"]
        request.Referer = params["HTTP_REFERER"]
        request.UserAgent = params["HTTP_USER_AGENT"]

        request.Close = false
        //request.Body
        request.ContentLength, _ = strconv.Atoi64(params["HTTP_CONTENT_LENGTH"])

        request.Path = params["PATH_INFO"]
        request.QueryString = params["QUERY_STRING"]
        request.ScriptName = params["SCRIPT_NAME"]
        request.HttpCookie = params["HTTP_COOKIE"]
        request.cookies = ParseCookies(request.HttpCookie)

        if request.session != nil { return }
        if err = request.initSession(); err != nil {
                logger.Printf("error: request.initSession: '%v'", err)
                //return
        }

        return
}

func (fcgi *FCGIModel) sendResult(out io.Writer, rm RequestManager, h *RecordHeader) (err os.Error) {
        logger.Printf("sending result: %v\n", h.RequestId)

        hh := &RecordHeader{
        Version: h.Version,
        Type: FCGI_STDOUT,
        RequestId: h.RequestId,
        ContentLength: 0,
        PaddingLength: 0,
        }

        defer func() {
                if err := recover(); err != nil {
                        str := bytes.NewBuffer(make([]byte, 0, 2048))
                        printErrorCallStack(err, str)
                        hh.ContentLength = uint16(str.Len())
                        if err = binary.Write(out, binary.BigEndian, hh); err != nil { return }
                        if _, err = io.WriteString(out, str.String()); err != nil { return }
                }
        }()
                
        var request *Request
        request, err = rm.GetRequest(fmt.Sprintf("%v", uint16(h.RequestId)))
        if err == nil && request != nil {
                if request.session == nil {
                        // TODO: init special session for FCGI, no FS and DB session storage
                        if err = request.initSession(); err != nil {
                                logger.Printf("error: request.initSession: '%v'", err)
                                //return
                        }
                }

                var response *Response
                response, err = rm.ProcessRequest(request)
                if err == nil && response != nil {
                        str := bytes.NewBuffer(make([]byte, 0, 2048))
                        if err = response.writeHeader(str); err != nil { return }
                        if _, err = io.Copy(str, response.Body); err != nil { return }
                        //logger.Printf("result:\n%v", str)
                        //logger.Printf("result:==========\n")

                        hh.ContentLength = uint16(str.Len())
                        if err = binary.Write(out, binary.BigEndian, hh); err != nil { return }
                        if _, err = io.Copy(out, str); err != nil { return }
                }
        } else {
                panic(fmt.Sprintf("<b>unknown request</b>: %v", h.RequestId))
        }

        return
}

func (fcgi *FCGIModel) processSession(rm RequestManager, fd int) {
        //logger.Printf("FCGI_WEB_SERVER_ADDRS: %s\n", os.Getenv("FCGI_WEB_SERVER_ADDRS"))

        f := os.NewFile(fd, "FCGI_LISTENSOCK_FILENO")
        if f == nil {
                logger.Printf("error: os.NewFile: %v\n", fd)
                return
        }

        defer f.Close()

        h := new(RecordHeader)
        in := bufio.NewReader(f)

        for {
                ok, err := h.read(in)
                if !ok || err != nil {
                        logger.Printf("error: RecordHeader.read: %v, %v\n", ok, err)
                        return
                }

                //logger.Printf("header: %v\n", h)

                content, err := h.readContent(in)
                if err != nil {
                        logger.Printf("error: RecordHeader.readContent: %v\n", err)
                        return
                }

                switch h.Type {
                case FCGI_BEGIN_REQUEST:
                        rec := &BeginRequest{ header: h }
                        ok = rec.parse(content)
                        if !ok {
                                logger.Printf("error: BeginRequest.parse\n")
                                return
                        }
                        //logger.Printf("BeginRequest: %v\n", rec)
                case FCGI_PARAMS:
                        if 0 < h.ContentLength {
                                params := parseNameValuePares(content)
                                //logger.Printf("params: %v\n", params)

                                var request *Request
                                request, err = rm.GetRequest(fmt.Sprintf("%v", uint16(h.RequestId)))
                                if request == nil {
                                        logger.Printf("no request: %v", h.RequestId)
                                        return
                                }
                                if err != nil {
                                        logger.Printf("error: %v", err)
                                        return
                                }

                                err = initRequest(request, params)
                                if err != nil {
                                        logger.Printf("error: %v", err)
                                        //return
                                }
                        }
                case FCGI_STDIN:
                        if h.ContentLength == uint16(0) {
                                fcgi.sendResult(f, rm, h)

                                hh := &RecordHeader{
                                Version: h.Version,
                                Type: FCGI_END_REQUEST,
                                RequestId: h.RequestId,
                                ContentLength: 8,
                                PaddingLength: 0,
                                }
                                er := &EndRequestBody{
                                AppStatus: 0,
                                ProtocolStatus: FCGI_REQUEST_COMPLETE,
                                }
                                if err = binary.Write(f, binary.BigEndian, hh); err != nil { return }
                                if err = binary.Write(f, binary.BigEndian, er); err != nil { return }

                                logger.Printf("request ended\n")
                        } else {
                                logger.Printf("TODO: obtain data from the web server\n")
                                //TODO: request.Body = noCloseReader{ body }
                        }
                        return
                default:
                        logger.Printf("content: %v\n", content)
                }
        }
}

func (fcgi *FCGIModel) ProcessRequests(rm RequestManager) (err os.Error) {
        logFile, _ := os.Open("/tmp/a.go.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
        logger = log.New(logFile, "", log.Lshortfile|log.Ltime)

        logger.Printf("=================================\n")
        logger.Printf("ARGS: %v\n", os.Args)

        logger.Printf("Listenning...\n")
        for {
                fd, _, ec := syscall.Accept(FCGI_LISTENSOCK_FILENO)
                if ec != 0 {
                        logger.Printf("error: %s\n", os.Errno(ec).String())
                        return
                }

                logger.Printf("Accepted: <%v>\n", fd)

                go fcgi.processSession(rm, fd)
        }
        return
}
