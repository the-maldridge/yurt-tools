{{define "content"}}
  {{range $job, $jobdata := .}}
    {{range $group, $groupdata := $jobdata}}
      {{range $task, $taskdata := $groupdata}}
        {{$color := "grey"}}
        <div class="card">
          {{$upToDate := "unknown"}}
          {{$secure := "unknown"}}
          {{if $taskdata.versions}}
            {{if not $taskdata.versions.NonComparable}}
              {{if $taskdata.versions.UpToDate}}
                {{$upToDate = "true"}}
              {{else}}
                {{$upToDate = "false"}}
              {{end}}
            {{end}}
          {{end}}
          {{if $taskdata.trivy}}
            {{$secure = "true"}}
            {{range $target, $targetdata := $taskdata.trivy}}
              {{if $targetdata.Vulnerabilities}}
                {{$secure = "false"}}
              {{end}}
            {{end}}
          {{end}}
          {{if or (eq $secure "false") (eq $upToDate "false")}}
            {{$color = "red"}}
          {{else if or (eq $secure "true") (eq $upToDate "true")}}
            {{$color = "green"}}
          {{end}}
          <div class="ring-container">
            <div class="ringring {{$color}}"></div>
            <div class="circle {{$color}}"></div>
          </div>
          <div class="card-content">
            <a class="card-title" href="/detail/{{$job}}/{{$group}}/{{$task}}">
              {{$job}} &#187; {{$task}}
            </a>
            <hr/>
            {{if or (eq $secure "false") (eq $upToDate "false")}}
              {{if and (eq $secure "false") (eq $upToDate "false")}}
                {{$taskdata.metadata.Docker.Tag}} has <em>multiple issues</em>
              {{else if (eq $secure "false")}}
                {{$taskdata.metadata.Docker.Tag}} has <em>security issues</em>
              {{else}}
                {{$taskdata.metadata.Docker.Tag}} is <em>out of date</em>
              {{end}}
            {{else if and (eq $secure "unknown") (eq $upToDate "unknown")}}
              {{$taskdata.metadata.Docker.Tag}} has unknown status
            {{else}}
              {{$taskdata.metadata.Docker.Tag}} is looking good
            {{end}}
          </div>
        </div>
      {{end}}
    {{end}}
  {{end}}
{{end}}
