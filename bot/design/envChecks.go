package design

import (
	"slices"
	"strconv"

	botStructures "PanoptisMouthNew/structures/bot"
)

func RunChecks(variables []botStructures.VariableType) (map[string]string, []string, bool) {
	var notOkay []string

	for _, variable := range variables {
		if len(variable.RawValue) == 0 {
			notOkay = append(notOkay, variable.Name)
		}

		if variable.IsNumeric {
			isValid := IsNumeric(variable.RawValue, variable.SignedInteger)
			if !isValid {
				notOkay = append(notOkay, variable.Name)
			}
		}
	}

	var okayVariables = make(map[string]string)

	for _, variable := range variables {
		if variable.Name != "" && !slices.Contains(notOkay, variable.Name) {
			okayVariables[variable.Name] = variable.RawValue
		}
	}

	return okayVariables, notOkay, len(notOkay) == 0
}

func IsNumeric(str string, signed bool) bool {
	if signed {
		_, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return false
		}
	} else {
		_, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return false
		}
	}
	return true
}
