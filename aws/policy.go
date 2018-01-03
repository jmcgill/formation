package aws

import "encoding/json"

func SafeCleanPolicy(policy string) string {
	// ${} references in Policy documents conflict with Terraform. Use &{} instead
	var d map[string]interface{}
	err := json.Unmarshal([]byte(policy), &d)
	if err != nil {
		return policy
	}

	s, _ := json.MarshalIndent(d, "", "    ")
	return string(s) + "\n"
}
