package gae

import (
	"net/http"
	"os"

	"github.com/Manbeardo/mtgdocs.info/parse"
	"github.com/emicklei/go-restful"
)

func init() {
	readFiles()

	ws := new(restful.WebService).Path("/")
	ws.Route(ws.GET("/cr/").To(serveCR))
	ws.Route(ws.GET("/cr/{Title}").
		Param(ws.PathParameter("Title", "the title of the rule")).
		To(serveCR))

	restful.Add(ws)
}

var cr parse.ComprehensiveRules

func readFiles() {
	file, err := os.Open("etc/C15.cr.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	cr = parse.ParseCR(file)
}

func serveCR(req *restful.Request, resp *restful.Response) {
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

	resp.Write([]byte(rule.CompleteText()))
}
