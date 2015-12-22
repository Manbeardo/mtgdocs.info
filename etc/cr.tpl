{{define "rule"}}
    {{$rule := .}}
    {{if ne .Title "HEAD"}}
        <dt id="{{.Title}}">
            <a href="/cr/{{.Title}}">{{.Title}}</a>
        </dt>
        <dd>{{.Body}}</dd>
        {{range .Notes}}
            <dd>{{.}}</dd>
        {{end}}
        {{range .Examples}}
            <dd>
                <br/>
                <div class="panel panel-default">
                    <div class="panel-heading text-muted">Example</div>
                    <div class="panel-body">{{.}}</div>
                </div>
            </dd>
        {{end}}
        {{if ne 0 (len .Children)}}
            <br/>
            <dd>
                <dl class="dl-horizontal">
                    {{range $index, $child := $rule.Children}}
                        {{template "rule" $child}}
                        {{if and (eq 0 (len $child.Children)) (ne (sum $index 1) (len $rule.Children))}}
                            <br/>
                        {{end}}
                    {{end}}
                </dl>
            </dd>
        {{end}}
    {{else}}
        {{range .Children}}
            {{template "rule" .}}
        {{end}}
    {{end}}
{{end}}
<html>
<head>
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" integrity="sha384-1q8mTJOASx8j1Au+a5WDVnPi2lkFfwwEAa8hDDdjZlpLegxhjVME1fgjWPGmkzs7" crossorigin="anonymous">
<style>
@media (min-width: 768px) {
    .dl-horizontal dt {
        width: 55;
    }
    .dl-horizontal dd {
        margin-left: 70;
    }
    #dl-toplevel {
        margin-right: 70;
    }
}
</style>
</head>
<title>
{{if ne .Title "HEAD"}}
    {{.Title}} - MTG CR
{{else}}
    MTG Comprehensive Rules
{{end}}
</title>
<body>
    <div class="container">
        <h3><a href="/cr" class="text-muted">Magic: the Gathering Comprehensive Rules</a></h3>
        <div class="well">
            <dl class="dl-horizontal" id="dl-toplevel">
                <dt>
                    <div class="pull-left">
                    {{if .Parent}}
                        {{if eq .Parent.Title "HEAD"}}
                            <a href="/cr">..</a>
                        {{else}}
                            <a href="/cr/{{.Parent.Title}}">..</a>
                        {{end}}
                    {{end}}
                    </div>
                </dt><dd></dd>
                {{template "rule" .}}
            </dl>
        </div>
    </div>
</body>
</html>
