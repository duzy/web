package main

import (
        "os"
        "fmt"
        "bufio"
        "template"
        "strconv"
)

func error(e os.Error) {
        fmt.Printf("error: %s\n", e)
        os.Exit(0)
}

func check_error(e os.Error) {
        if e != nil { error(e) }
}

func make_fields() (map[string]string) {
        vs := make(map[string]string)
        vs["title"] = "hello"
        vs["body"] = "Go"

        vs["REQUEST_METHOD"] = os.Getenv("REQUEST_METHOD")
        vs["PATH_INFO"] = os.Getenv("PATH_INFO")
        vs["QUERY_STRING"] = os.Getenv("QUERY_STRING")

        vs["environ"] = ""
        for i, v := range os.Environ() {
                vs["environ"] += strconv.Itoa(i) + ": " + v + "<br/>"
        }

        return(vs)
}

func main() {
        var err os.Error

        w, err := bufio.NewWriterSize(os.Stdout, 1024*10)
        check_error(err)
        defer w.Flush()

        w.WriteString("Content-Type: text/html;\n\n")

        t, err := template.ParseFile("foo.tpl", nil)
        check_error(err)

        check_error( t.Execute( make_fields(), w ) )
}
