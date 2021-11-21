package v1

import (
	"strconv"

	"github.com/saime-0/http-cute-chat/internal/api/validator"

	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/models"
)

func MatchMessageType(input *models.FormCompleted, sample *models.FormPattern) (models.FormCompleted, *rules.AdvancedError) {
	completed := make(map[string]string)
	for _, field := range sample.Fields {
		for _, choice := range input.Input {
			if choice.Key == field.Key {
				var adv_err *rules.AdvancedError
				if len(choice.Value) > field.Length && field.Length > 0 {
					adv_err = rules.ErrChoiceValueLength
				}
				switch field.Type {
				case rules.TextField:
					// nothing

				case rules.DateField:
					if _, err := strconv.ParseInt(choice.Value, 10, 64); err != nil {
						adv_err = rules.ErrInvalidChoiceDate
					}

				case rules.EmailField:
					if !validator.ValidateEmail(choice.Value) {
						adv_err = rules.ErrInvalidEmail
					}

				case rules.LinkField:
					if !validator.ValidateLink(choice.Value) {
						adv_err = rules.ErrInvalidLink
					}

				case rules.NumericField:
					if _, err := strconv.Atoi(choice.Value); err != nil {
						adv_err = rules.ErrInvalidChoiceValue
					}

				default:
					adv_err = rules.ErrDataRetrieved
				}
				if adv_err != nil {
					return models.FormCompleted{}, adv_err
				}
				if len(field.Items) != 0 {
					contains := func(arr []string, str string) bool {
						for _, a := range arr {
							if a == str {
								return true
							}
						}
						return false
					}(field.Items, choice.Value)

					if !contains {
						return models.FormCompleted{}, rules.ErrInvalidChoiceValue
					}
				}
				completed[field.Key] = choice.Value
			}

		}
		_, ok := completed[field.Key]
		if !(ok || field.Optional) {
			return models.FormCompleted{}, rules.ErrMissingChoicePair
		}

	}
	return mapToFormCompleted(&completed), nil
}

func mapToFormCompleted(inp *map[string]string) (form models.FormCompleted) {
	for k, v := range *inp {
		form.Input = append(form.Input, models.FormChoice{
			Key:   k,
			Value: v,
		})
	}
	return
}
