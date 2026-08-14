package vault

type Group struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	MemberEntityIds []string `json:"member_entity_ids"`
}

func Deduplicate(s []string) []string {
	if len(s) == 0 {
		return s
	}
	seen := make(map[string]bool)
	var result []string
	for _, val := range s {
		if !seen[val] {
			seen[val] = true
			result = append(result, val)
		}
	}
	return result
}
