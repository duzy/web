package web

import (
        "bytes"
        "bufio"
        "os"
        "io"
        //"io/ioutil"
        "fmt"
        "log"
        "syscall"
        "encoding/binary"
)

type FCGIModel struct {
        params map[string]string
}

func NewFCGIModel() (am *AppModel) {
        m := &FCGIModel{}
        am = AppModel(m)
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
        return fmt.Sprintf("Version<%d>", uint8(v))
}

func (ri RequestId) String() string {
        return fmt.Sprintf("RequestId<%d>", uint16(ri))
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
        0: "BeginRequestFlag<0>",
        FCGI_KEEP_CONN: "FCGI_KEEP_CONN",
        }
        if len(names) <= int(f) {
                return fmt.Sprintf("BeginRequestFlag<%v>", uint16(f))
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

func (h *RecordHeader) read(logger *log.Logger, r io.Reader) (ok bool, err os.Error) {
        b := make([]byte, 8)
        l, err := io.ReadFull(r, b)
        if l != 8 {
                if err == nil {
                        err = os.NewError(fmt.Sprintf("bad header read: %d", l))
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
                                err = os.NewError(fmt.Sprintf("Short read %d bytes of %d", n, l))
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

func getSize(b []byte) (size uint, ob []byte) {
        s := b[0]
        if s>>7 == 1 {
                if 4 <= len(b) {
                        size = uint(binary.BigEndian.Uint32(b[0:4]))
                        ob = b[4:]
                }
        } else {
                if 0 < len(b) {
                        size = uint(s)
                        ob = b[1:]
                }
        }
        return
}

func getValue(b []byte, size uint) (v string, ob []byte) {
        if len(b) < int(size) {
                v = string(b[0:])
                ob = b[0:0]
        } else {
                v = string(b[0:size])
                ob = b[size:]
        }
        return
}

func parseParams(b []byte) (params map [string]string) {
        var kl, vl uint
        var k, v string
        params = make(map[string]string)
        for 0 < len(b) {
                if kl, b = getSize(b);          b == nil { break }
                if vl, b = getSize(b);          b == nil { break }
                if k, b = getValue(b, kl);      b == nil { break }
                if v, b = getValue(b, vl);      b == nil { break }
                //logger.Printf("kv: %v, %v, %d\n", k, v, len(b))
                params[k] = v
        }
        return
}

var counter = 0

func (fcgi *FCGIModel) processSession(rp RequestProcessor, logger *log.Logger, fd int) {
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
                ok, err := h.read(logger, in)
                if !ok || err != nil {
                        logger.Printf("error: RecordHeader.read: %v, %v\n", ok, err)
                        return
                }

                logger.Printf("header: %v\n", h)

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
                        logger.Printf("BeginRequest: %v\n", rec)
                case FCGI_PARAMS:
                        if 0 < h.ContentLength {
                                params := parseParams(content)
                                logger.Printf("params: %v\n", params)
                        }
                case FCGI_STDIN:
                        if h.ContentLength == uint16(0) {
                                out := bufio.NewWriter(f)
                                {
                                        counter += 1
                                        str := bytes.NewBuffer(make([]byte, 0, 2048))
                                        fmt.Fprintf(str, "Content-Type: text/html\n\n")
                                        fmt.Fprintf(str, "<b>test</b>, num=%d", counter)
                                        hh := &RecordHeader{
                                        Version: h.Version,
                                        Type: FCGI_STDOUT,
                                        RequestId: h.RequestId,
                                        ContentLength: uint16(str.Len()),
                                        PaddingLength: 0,
                                        }
                                        if err = binary.Write(out, binary.BigEndian, hh); err != nil { return }
                                        if _, err = io.WriteString(out, str.String()); err != nil { return }
                                }
                                {
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
                                        if err = binary.Write(out, binary.BigEndian, hh); err != nil { return }
                                        if err = binary.Write(out, binary.BigEndian, er); err != nil { return }
                                }
                                out.Flush()
                                logger.Printf("request ended\n")
                        } else {
                                logger.Printf("TODO: obtain data from the web server\n")
                        }
                        return
                default:
                        logger.Printf("content: %v\n", content)
                }
        }
}

func (fcgi *FCGIModel) ProcessRequests(rp RequestProcessor) (err os.Error) {
        logFile, _ := os.Open("/tmp/a.go.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
        logger := log.New(logFile, "", log.Lshortfile|log.Ltime)

        logger.Printf("=================================\n")
        logger.Printf("ARGS: %v\n", os.Args)

        logger.Printf("Listenning...\n")
        for {
                fd, _, ec := syscall.Accept(FCGI_LISTENSOCK_FILENO)
                if ec != 0 {
                        logger.Printf("error: [%v] %s\n", ec, os.Errno(ec).String())
                        return
                }

                logger.Printf("Accepted: <%v>\n", fd)

                go fcgi.processSession(rp, logger, fd)
        }
}

func (fcgi *FCGIModel) GetRequest() (req *Request) {
        // TODO: ...
        return
}
