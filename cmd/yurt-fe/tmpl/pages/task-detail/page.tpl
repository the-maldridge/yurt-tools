{{ define "content" }}
{{ with .metadata }}
<div id="task-meta">
  Job: {{.Job}}<br />
  Group: {{.Group}}<br />
  Task: {{.Name}}<br />
  Driver: {{.Driver}}<br />
  <br />
  {{ if eq .Driver "docker" }}
  {{ with .Docker }}
  Owner: {{.Owner}}<br />
  Image: {{.Image}}<br />
  Tag: {{.Tag}}<br />
  {{ end }}
  {{ end }}
</div>
{{ end }}
{{ if .versions }}
{{ with .versions }}
<div id="task-versions">
  {{ . }}
</div>
{{ end }}
{{ end }}
{{ end }}
