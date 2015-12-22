package gae

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Manbeardo/mtgdocs.info/parse"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
)

func init() {
	readFiles()

	ws := new(restful.WebService).Path("/")

	paramDisplay := ws.QueryParameter("Display", "display mode (HTML, TXT, or JSON)").DefaultValue("HTML")

	ws.Route(ws.GET("/").To(serveIndex))
	ws.Route(ws.GET("/cr").
		Param(paramDisplay).
		To(serveCR))
	ws.Route(ws.GET("/cr/{Title}").
		Param(ws.PathParameter("Title", "title of the rule")).
		Param(paramDisplay).
		To(serveCR))

	restful.Add(ws)

	swagger.InstallSwaggerService(swagger.Config{
		WebServices: restful.DefaultContainer.RegisteredWebServices(),
		ApiPath:     "/apidocs.json",
	})
}

var cr parse.ComprehensiveRules
var funcMap = template.FuncMap{
	"sum": func(values ...int) int {
		sum := 0
		for _, val := range values {
			sum += val
		}
		return sum
	},
}
var crTpl = template.Must(template.New("cr.tpl").Funcs(funcMap).ParseFiles("etc/cr.tpl"))
var index []byte

func readFiles() {
	file, err := os.Open("etc/C15.cr.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	cr = parse.ParseCR(file)

	indexBytes, err := ioutil.ReadFile("etc/index.html")
	if err != nil {
		panic(err)
	}
	index = indexBytes
}

func serveIndex(req *restful.Request, resp *restful.Response) {
	resp.Write(index)
}

func serveCR(req *restful.Request, resp *restful.Response) {
	display := req.QueryParameter("Display")
	if display == "" {
		display = "HTML"
	} else if display != "HTML" && display != "JSON" && display != "TXT" {
		resp.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("illegal display mode: %v", display))
		return
	}

	var rule *parse.Rule
	if title := parse.Title(req.PathParameter("Title")); title == "" {
		rule = cr["HEAD"]
	} else if foundRule, ok := cr[title]; ok {
		rule = foundRule
	} else if foundRule, ok = cr[title+"."]; ok {
		rule = foundRule
	} else {
		resp.WriteErrorString(http.StatusNotFound, "no rule with that title found")
	}

	if display == "HTML" {
		crTpl.Execute(resp, rule)
	} else if display == "TXT" {
		resp.Write([]byte(rule.CompleteText()))
	} else if display == "JSON" {
		resp.WriteAsJson(rule)
	}
}
