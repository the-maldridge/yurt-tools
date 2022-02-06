package filters

import (
	"reflect"

	"github.com/flosch/pongo2/v4"
)

func init() {
	pongo2.RegisterFilter("countVulnerabilties", filterCountVulnerabilities)
}

// This expects to start from trivy.Results
func filterCountVulnerabilities(resultsArray *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	vulnerabilities := 0
	// this takes us into our list of maps with the key Vulnerabilities
	resultsArray.Iterate(
		func(idx, count int, resultObj, value *pongo2.Value) bool {
			// All hope abandon ye who enter here
			results := reflect.ValueOf(resultObj.Interface())
			for _, t := range results.MapKeys() {
				// Once we get to the Vulnerabilties key, add it to the total
				if t.String() == "Vulnerabilities" {
					vulnerabilities_list := reflect.ValueOf(results.MapIndex(t).Interface())
					vulnerabilities += vulnerabilities_list.Len()
				}
			}
			return true
		},
		func() {
		},
	)
	return pongo2.AsValue(vulnerabilities), nil
}
