package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/http/cgi"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type doc struct {
	DocTitle  string
	DocFooter string
	Sections  []section
}

type section struct {
	SectionTitle string
	SectionLink  string
	Items        []map[string]string
}

func main() {
	cgi.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		f, err := os.Create("/tmp/cgi.log")
		check(err)
		log.SetOutput(f)
		entry := r.FormValue("file")
		depth := r.FormValue("depth")
		log.Print("depth=" + depth)
		if entry == "" {
			entry = "Eclipse.Eclipse"
		}
		entrySlice := strings.Split(entry, ".")
		srv(entrySlice, w)
	}))
}

func Xmain() {
	entry := "Eclipse"
	entry = "Eclipse.Transfer"
	entrySlice := strings.Split(entry, ".")
	srv(entrySlice, os.Stdout)
}

func srv(entrySlice []string, w io.Writer) {
	fmt.Fprintf(os.Stderr, "entrySlice:%+v\n", entrySlice)
	var sections []section
	entry := entrySlice[len(entrySlice)-1]
	log.Print("entry=" + entry)
	hash, pages := read(0, entry)

	for _, i := range pages {
		page := hash[strconv.Itoa(i)]
		nhash, npages := read(i, page)
		sectionTitle := nhash["what"]
		sectionLink := "placeholder" // XXX
		var items []map[string]string

		if len(npages) > 0 {
			for _, j := range npages { // depth = 2
				npage := nhash[strconv.Itoa(j)]
				nnhash, _ := read(j, npage)
				m := make(map[string]string)
				m[mungSpaces(npage)] = nnhash["what"]
				items = append(items, m)

			}
		} else {
			items = append(items, nhash)
		}

		sections = append(sections, section{
			SectionTitle: sectionTitle,
			SectionLink:  sectionLink,
			Items:        items,
		})
	}

	t, err := template.New("webpage").Parse(templ())
	check(err)

	doc := doc{DocTitle: hash["what"], DocFooter: hash["why"]}
	doc.Sections = sections
	err = t.Execute(w, doc)

	fmt.Fprintf(os.Stderr, "doc:%+v\n", doc)
	//	pp.Print(doc)
	check(err)

}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func mungSpaces(s string) string {
	return strings.Replace(s, " ", "_", -1)
}

func read(i int, f string) (map[string]string, []int) {
	f = mungSpaces(f)
	fn := "./pages/" + f + "/current"

	var pages []int
	hash := make(map[string]string)

	if _, err := os.Stat(fn); os.IsNotExist(err) {
		hash[""] = f
	} else {
		fmt.Fprintln(os.Stderr, "\nopening "+fn)
		file, err := os.Open(fn)
		check(err)
		scanner := bufio.NewScanner(file)
		numRegex, _ := regexp.Compile(`^\d+$`)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}
			lineSlice := strings.Split(line, ": ")
			key := lineSlice[0]
			val := lineSlice[1]
			hash[key] = val
			if numRegex.Match([]byte(key)) {
				i, _ := strconv.ParseInt(key, 10, 32)
				pages = append(pages, int(i))
			}
		}
		sort.Ints(pages)
	}

	//pp.Print(hash)
	//pp.Print(pages)
	return hash, pages
}

func templ() string {
	return `
<html>
    <head>
    </head>
    <body>
        <center><br>
            <br>
            <a href="?file=Eclipse&amp;depth=2">{{.DocTitle}}</a>
            <table cellspacing="5" cellpadding="10">

			{{ range $j, $k := .Sections}}
                <tr>
                    <td bgcolor="#EEEEEE"><a href=
                    "?file=Eclipse.{{.SectionLink}}&amp;depth=2">{{.SectionTitle}}</a>

                        <table cellspacing="5" cellpadding="10">
                            <tr>
								{{ range $i, $e := .Items }}
								{{ range $k, $v := . }}
                                <td bgcolor="#CCCCCC">{{ if ne $k ""}}<a href="?file=Eclipse.Eclipse.{{$k}}&amp;depth=2">{{end}}{{$v}}</a></td>
								{{ end }}
								{{ end }}
                            </tr>
                        </table>
                    </td>
                </tr>
				{{end}}
				</table>

				<table width="500">
				  <tr>
				    <td><font color="gray"> {{.DocFooter}}
				    </font><font color="gray">Use <a href=
				"?file=Eclipse&amp;depth=2&amp;edit=on">edit</a> to change this and
				nearby elements.</font></td>
				    </tr>
				</table>
    </body>
</html>
`
}
