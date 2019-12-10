<html>
  <head>
    <title>up2date</title>
    <style>
      body {
      padding: 1em;
      }

      .masonry {
      column-count: 4;
      column-gap: 1em;
      }

      .card {
      width: 100%;
      display: inline-block;
      border: 2px solid black;
      border-radius: 4px;
      margin: 0 0 1em;
      }

      .card-content {
      padding: 1em;
      }
    </style>
  </head>
  <body>
    <div>
      <h1>Are you up2date?</h1>
      Updated at {{ .Updated }}
    </div>
    <hr />
    <div class="masonry">
      {{ range $element := .TaskList }}
      <div class="card">
        <div class="card-content">
          <h2>{{$element.Name}}</h2>
          <hr />
          Using Image: {{$element.Image}}<br />
          Current Version: {{$element.Version}}<br />
          {{ if $element.NoData }}
          No data could be loaded for this image.
          {{ else }}
          {{ if $element.Newer }}
          Newer Versions:<br />
          <ul>
            {{ range $version := $element.Newer }}
            <li>{{ $version }}</li>
            {{ end }}
          </ul>
          {{ else }}
          You are on the newest version.
          {{ end }}
          {{ end }}
        </div>
      </div>
      {{ end }}
    </div>
  </body>
</html>
