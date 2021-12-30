package model

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"io"
	"regexp"
	"strings"
)

type HexColor string

func MarshalHexColor(h HexColor) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, fmt.Sprintf(`"%s"`, h))
	})
}

func UnmarshalHexColor(v interface{}) (HexColor, error) {
	color, err := graphql.UnmarshalString(v)
	if err != nil {
		return "", err
	}

	var re = regexp.MustCompile(`(?m)^#([0-9a-fA-F]{3}){1,2}$`)
	if !re.MatchString(color) {
		return "", fmt.Errorf("%s is not an HexColor", color)
	}

	return HexColor(strings.ToUpper(color)), nil
}
