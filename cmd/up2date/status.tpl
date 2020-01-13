<html>
  <head>
    <title>up2date</title>
    <style>
      body {
        padding: 1em;
        background-color: #F9F9FA;
        font-family: sans-serif;
      }

      h1, h2{
        font-weight: 300;
      }

      .topic-name {
        font-weight: 300;
      }

      .masonry {
        column-count: 4;
        column-gap: 1em;
      }

      .topic-message {
        display: block;
        text-align: center;
        margin-top: 10px;
      }

      .card {
        width: 100%;
        display: inline-block;
        border-radius: 4px;
        box-shadow: 0 1px 3px rgba(0,0,0,0.12), 0 1px 2px rgba(0,0,0,0.24);
        margin: 0 0 1em;
        cursor: pointer;
        transition: all 0.3s cubic-bezier(.25,.8,.25,1);
      }

      a.card {
        color: #000000;
        text-decoration: none;
      }

      .card:hover {
        box-shadow: 0 14px 28px rgba(0,0,0,0.25), 0 10px 10px rgba(0,0,0,0.22);
      }

      .card.upToDate {
        background-color: #ffffff;
      }

      .card.noData {
        background-color: #eceff1;
      }

      .card.outOfDate {
        background-color: #ffebee;
      }

      hr {
        border-top-color: rgba(0,0,0,.12);
        color: transparent;
      }

      .card-content {
        padding: 1em;
      }
    </style>
  </head>
  <body>
    {{ if not .TaskList }}
    {{ if eq .Updated.Year 1 }}
    <div>
      <h1>Loading data</h1>
      Check back soon
    </div>
    {{ else }}
    No Data
    {{ end }}
    {{ else }}
    <div>
      <h1>Are you up2date?</h1>
      Updated at {{ .Updated.Format "Jan 02, 2006 15:04:05" }}
    </div>
    {{ end }}
    <hr />
    <div class="masonry">
      {{ range $element := .TaskList }}
      <a class="card {{ if $element.Newer }}outOfDate{{ else if $element.NoData }}noData{{ else }}upToDate{{ end }}" href="{{ $element.Url }}" target="_blank">
        <div class="card-content">
          <h2>{{$element.Name}}</h2>
          <hr />
          <span class="topic-name">Using Image:</span> {{$element.Image}}<br />
          <span class="topic-name">Current Version:</span> {{$element.Version}}<br />
          {{ if $element.NoData }}
          <span class="topic-message">No data could be loaded for this image.</span>
          {{ else }}
          {{ if $element.Newer }}
          <span class="topic-name">Newer Versions:</span><br />
          <ul>
            {{ range $version := $element.Newer }}
            <li>{{ $version }}</li>
            {{ end }}
          </ul>
          {{ else }}
          <span class="topic-message">You are on the newest version.</span>
          {{ end }}
          {{ end }}
        </div>
      </a>
      {{ end }}
    </div>
  </body>
</html>
