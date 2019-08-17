package utils

import (
	"errors"
	"fmt"
	"strings"
)

func MapifyLabels(cid string, headers string) (map[string]string, error) {
	var err error

	mapLabels := map[string]string{}

	parsedLabels := strings.Split(headers, ",")

	for _, value := range parsedLabels {
		parsedLabel := strings.Split(value, "=")
		mapLabels[parsedLabel[0]] = parsedLabel[1]
	}

	if len(mapLabels) == 0 {
		return map[string]string{}, err
	}

	return mapLabels, nil
}

// StringifyLabelSelector returns a k8s selector string based on given map. Ex: "key=value,key2=value2"
func StringifyLabelSelector(cid string, labelSelector map[string]string) (string, error) {

	var labelsPair []string

	for key, value := range labelSelector {
		labelsPair = append(labelsPair, fmt.Sprintf("%s=%s", key, value))
	}

	if len(labelsPair) == 0 {
		return "", errors.New("got an empty labelSelector")
	}

	return strings.Join(labelsPair[:], ","), nil
}

func ValidateCobraStringFlag(flag string) error {
	var err error

	if flag == "" {
		return err
	}
	return nil
}
