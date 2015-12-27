package webui

import (
	"html/template"
)


func addOne(x int) int {
    return x + 1
}

var funcMap = template.FuncMap{
    "addOne": addOne,
}

var FilesTpl = template.Must(template.New("files").Funcs(funcMap).Parse(`<!DOCTYPE html>
<html>
  <head>
    <title>KMS</title>
	<link rel="icon" href="favicon.png" sizes="32x32" />  
    <link rel="stylesheet" href="/css/bootstrap.min.css">
	<script type="text/javascript" src="/js/jquery-2.1.3.min.js"></script>
	<style>
	#jqstooltip{
		height: 28px !important;
		width: 150px !important;
	}
	</style>
  </head>
  <body>
    <div class="container">
      <div class="page-header">
	    <h1>
	      上传者 <small>{{ .Author }}</small>
	    </h1>
      </div>
      <div class="row">
        <h2>文件列表</h2>
        <table class="table table-striped">
          <thead>
            <tr>
              <th>Id</th>
              <th>Title</th>
              <th>Fid</th>
            </tr>
          </thead>
          <tbody>
           {{ $master := .Master }}
          {{ range $index, $p := .Files }}
         
              <tr>
              <td><code>{{$index | addOne}}</code></td>
              <td>{{ $p.name }}</td>
              <td><a href="{{ $master }}/{{$p.fid}}">View</a></td>
            </tr>
          {{ end }}
          </tbody>
        </table>
      </div>

    </div>
  </body>
</html>
`))
