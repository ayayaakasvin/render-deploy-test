package httpserver

import (
	"html/template"
	"net/http"
)

var HelloTemplate = template.Must(template.New("hello").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Hello</title>
</head>
<body>
    <h1>Hello, deployment is working!</h1>
</body>
</html>
`))

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	HelloTemplate.Execute(w, nil)
}
