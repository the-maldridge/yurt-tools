{{ define "base" }}
<html>
  <head>
    <title>Yurt Tools</title>
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
