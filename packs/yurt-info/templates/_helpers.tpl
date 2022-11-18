// only deploys to a region if specified
[[ define "region" -]]
[[- if not (eq .my.region "") -]]
  region = [[ .my.region | quote]]
[[- end -]]
[[- end -]]

[[ define "namespace" ]]
[[- if not (eq .my.namespace "") -]]
  namespace = [[ .my.namespace | quote]]
[[- end -]]
[[- end ]]
