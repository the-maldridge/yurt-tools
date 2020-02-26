{{ define "content" }}
    <div class="row">
  {{ with .metadata }}
    <div class="thirdcolumn">
      <h3>Nomad</h3>
      <ul>
        <li><b>Job</b>: {{.Job}}<li>
            <li><b>Group</b>: {{.Group}}<li>
                <li><b>Task</b>: {{.Name}}<li>
                    <li><b>Driver</b>: {{.Driver}}<li>
      </ul>
      </div>
      {{ if eq .Driver "docker" }}
        {{ with .Docker }}
    <div class="thirdcolumn">
          <h3>Docker</h3>
          <ul>
            <li><b>Owner</b>: {{.Owner}}<li>
                <li><b>Image</b>: {{.Image}}<li>
                    <li><b>Tag</b>: {{.Tag}}<li>
          </ul>
    </div>
        {{ end }}
      {{ end }}
  {{ end }}
  {{ if .versions }}
    {{ with .versions }}
    <div class="thirdcolumn">
        <h3>Versions</h3>
        <ul>
          <li><b>Current</b>: {{ .Current }}</li>
          {{if not .UpToDate}}
            <li><b>Available versions</b>:
              <ul>
                {{range $_, $version := .Available}}
                  <li>{{$version}}</li>
                {{end}}
              </ul>
            </li>
          {{end}}
        </ul>
      </div>
    </div>
    {{ end }}
  {{ end }}
  </div>
  <hr/>
  {{ if .trivy }}
    {{with .trivy}}
        <h3>Trivy Results</h3>
        <br />
        {{range $_, $target := .}}
          <h4>{{ $target.Target }}</h4>
          {{if not $target.Vulnerabilities}}
            <br/>
            Target is secure<br/>
          {{else}}
            <table>
              <thead>
                <tr>
                  <th scope="col">Package</th>
                  <th scope="col">ID</th>
                  <th scope="col">Installed Version</th>
                  <th scope="col">Fixed Version</th>
                  <th scope="col">Title</th>
                  <th scope="col">Severity</th>
                </tr>
              </thead>
              <tbody>
                {{range $_, $vuln := $target.Vulnerabilities}}
                  <tr>
                    <td scope="row">{{$vuln.PkgName}}</td>
                    <td style="max-width: 15em">{{$vuln.VulnerabilityID}}
                      <span class="references">
                        {{- range $i, $ref := $vuln.References}}<a href="{{$ref}}">{{$i}}</a>&thinsp;{{- end}}
                        </span>
                    </td>
                    <td>{{$vuln.InstalledVersion}}</td>
                    <td>{{$vuln.FixedVersion}}</td>
                    <td class="tooltip">{{$vuln.Title}}
                      <span class="tooltiptext">{{$vuln.Description}}</span></td>
                    <td>
                      <span style="color:
                    {{if eq $vuln.Severity "HIGH"}}
                      #aa1523
                    {{else if eq $vuln.Severity "MEDIUM"}}
                      #dd9900
                    {{else}}
                      black
                    {{end}}">
                    {{$vuln.Severity}}
                      </span>
                    </td>
                  </tr>
                {{end}}
              </tbody>
            </table>
          {{end}}
          <br/>
        {{end}}
      </div>
    {{ end }}
  {{ end }}
{{ end }}
