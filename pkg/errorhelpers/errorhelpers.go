package errorhelpers

import "strings"

type MyErrorMap map[string]error

func JoinErrorMap(errMap MyErrorMap, key string, err error) MyErrorMap {
	if err == nil {
		return errMap
	}
	if errMap == nil {
		errMap = make(MyErrorMap)
	}
	errMap[key] = err
	return errMap
}

func (errMap MyErrorMap) Error() string {
	strBuilder := strings.Builder{}
	for key, err := range errMap {
		strBuilder.WriteString(key)
		strBuilder.WriteString(": ")
		strBuilder.WriteString(err.Error())
		strBuilder.WriteString("\n")
	}
	return strBuilder.String()
}

func (errMap MyErrorMap) Unwrap() []error {
	errs := make([]error, 0, len(errMap))
	for _, err := range errMap {
		errs = append(errs, err)
	}
	return errs
}
