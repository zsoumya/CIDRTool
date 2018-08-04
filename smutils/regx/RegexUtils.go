package regx

import "regexp"

func MatchNamedGroups(pattern string, input string) map[string]string {
	re := regexp.MustCompile(pattern)

	if re.MatchString(input) {
		groups := re.FindStringSubmatch(input)

		match := make(map[string]string)

		for groupNamePos, groupName := range re.SubexpNames() {
			if groupName != "" {
				match[groupName] = groups[groupNamePos]
			}
		}

		return match
	} else {
		return nil
	}
}
