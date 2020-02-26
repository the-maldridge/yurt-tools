{{ define "base" }}
<html>
  <head>
    <title>Yurt Tools</title>
    <link rel="stylesheet" href="/static/reset.css"> 
    <link rel="stylesheet" href="/static/style.css"> 
  </head>
  <body>
    {{ with .data }}
    {{ block "content" . }}
    No template for this page.
    {{ end }}
    {{ end }}
  </body>
</html>
{{ end }}
