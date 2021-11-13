package v1

import (
	"strconv"

	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/models"
)

func MatchMessageType(input *models.FormCompleted, sample *models.FormPattern) (map[string]string, *rules.AdvancedError) {
	completed := make(map[string]string)
	for _, field := range sample.Fields {
		for _, choice := range input.Input {
			if choice.Key == field.Key {
				if len(choice.Value) > field.Length && field.Length > 0 {
					return completed, rules.ErrChoiceValueLength
				}
				switch field.Type {
				case string(rules.TextField):
					break

				case string(rules.DateField):
					if _, err := strconv.ParseInt(choice.Value, 10, 64); err != nil {
						return completed, rules.ErrInvalidChoiceDate
					}

				case string(rules.EmailField):
					if !validateEmail(choice.Value) {
						return completed, rules.ErrInvalidEmail
					}

				case string(rules.LinkField):
					if !validateLink(choice.Value) {
						return completed, rules.ErrInvalidLink
					}

				case string(rules.NumericField):
					if _, err := strconv.Atoi(choice.Value); err != nil {
						return completed, rules.ErrInvalidChoiceValue
					}

				case string(rules.RadiobuttonField):
					if _, err := strconv.ParseBool(choice.Value); err != nil {
						return completed, rules.ErrInvalidChoiceValue
					}

				default:
					return completed, rules.ErrDataRetrieved
				}
				completed[field.Key] = choice.Value
			}

		}
		_, ok := completed[field.Key]
		if !(ok || field.Optional) {
			return completed, rules.ErrMissingChoicePair
		}

	}
	return nil, nil
}
