package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/cgi"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func main() {
	entry := "Eclipse"
	hash, pages := read(entry)
	fmt.Printf("\n%s", hash["what"])
	for _, i := range pages {
		page := hash[strconv.Itoa(i)]
		nhash, npages := read(page)
		fmt.Printf("\n%v\n", nhash["what"])
		for _, j := range npages { // depth = 2
			npage := nhash[strconv.Itoa(j)]
			nnhash, _ := read(npage)
			fmt.Printf("\n%s", nnhash["what"])
			fmt.Printf("[%s]\n", mungSpaces(npage))

		}
	}
	fmt.Printf("\n%s", hash["why"])

	fmt.Print("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

	//	t, err := template.New("webpage").Parse(templ())
	t, err := template.New("webpage").Parse(`{{range $index, $element := .}}<{{$index}}>|{{$element}}|{{end}}`)
	check(err)

	//	data := struct {
	//Entry string
	//}{
	//Entry: entry,
	//}

	data := map[string]string{"Mechanism": "Computer Interaction Mechanisms",
		"Browsing": "Browsing Web Sites", "Collaboration": "Collaborating Online"}

	err = t.Execute(os.Stdout, data)
	check(err)

	cgi.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		//		fmt.Fprintf(w, "Hello %s", r.FormValue("name"))

		fmt.Print(read(entry))

	}))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func mungSpaces(s string) string {
	return strings.Replace(s, " ", "_", -1)
}

func read(f string) (map[string]string, []int) {
	f = mungSpaces(f)
	fn := "./pages/" + f + "/current"
	//fmt.Println("\nopening " + fn)
	file, err := os.Open(fn)
	check(err)
	scanner := bufio.NewScanner(file)
	numRegex, _ := regexp.Compile(`^\d+$`)
	hash := make(map[string]string)
	var pages []int
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
            <a href="?file=Eclipse&amp;depth=2">Wikitect</a>
            <table cellspacing="5" cellpadding="10">

                <tr>
                    <td bgcolor="#EEEEEE"><a href=
                    "?file=Eclipse.Activity&amp;depth=2">Problem Solving Activity</a>
                        <table cellspacing="5" cellpadding="10">
                            <tr>
                                <td bgcolor="#CCCCCC"><a href="?file=Eclipse.Transfer&amp;depth=2">Transfer</a></td>
                                <td bgcolor="#CCCCCC"><a href="?file=Eclipse.Expression&amp;depth=2">Expression</a></td>
                                <td bgcolor="#CCCCCC"><a href="?file=Eclipse.Knowhow&amp;depth=2">Knowhow</a></td>
                            </tr>
                        </table>
                    </td>
                </tr>
                <tr>
    </body>
</html>
`
}
