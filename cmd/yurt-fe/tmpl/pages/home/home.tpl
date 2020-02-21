{{ define "content" }}
<table>
    <tr>
        <th>Job</th>
        <th>Group</th>
        <th>Task</th>
        <th>up2date?</th>
    </tr>

    {{ range $job, $jobdata := . }}
    {{ range $group, $groupdata := $jobdata }}
    {{ range $task, $taskdata := $groupdata }}
    <tr>
      <td>{{$job}}</td>
      <td>{{$group}}</td>
      <td><a href="/detail/{{$job}}/{{$group}}/{{$task}}">{{$task}}</a></td>
      <td>{{if not $taskdata.versions.NonComparable}}{{$taskdata.versions.UpToDate}}{{else}}Unknown{{end}}</td>
    </tr>
    {{ end }}
    {{ end }}
    {{ end }}
</table>
{{ end }}
