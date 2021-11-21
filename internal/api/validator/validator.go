package v1

import (
	"regexp"
	"strconv"
	"time"

	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/models"
)

func validateDomain(domain string) (valid bool) {
	if (len(domain) > rules.DomainMaxLength) || (len(domain) < rules.DomainMinLength) {
		return false
	}
	return true
}
func validateName(name string) (valid bool) {
	if (len(name) > rules.NameMaxLength) || (len(name) < rules.NameMinLength) {
		return false
	}
	return true
}
func validatePassword(password string) (valid bool) { // todo

	return len(password) <= rules.MaxPasswordLength && len(password) >= rules.MinPasswordLength
}
func validateEmail(email string) (valid bool) { // todo
	valid, _ = regexp.MatchString(`/^[A-Z0-9._%+-]+@[A-Z0-9-]+.+.[A-Z]{2,4}$/i`, email)
	return
}
func validateAppSettings(settings string) (valid bool) {

	return len(settings) <= rules.AppSettingsMaxLength
}
func validateLifetime(lt int64) (valid bool) {

	return lt >= int64(time.Minute) && lt <= rules.Year
}
func validateAliens(aliens int) (valid bool) {

	return aliens >= 1 && aliens <= 99_999
}
func validateLink(link string) (valid bool) {
	valid, _ = regexp.MatchString(`/^((ftp|http|https):\/\/)?(www\.)?([A-Za-zА-Яа-я0-9]{1}[A-Za-zА-Яа-я0-9\-]*\.?)*\.{1}[A-Za-zА-Яа-я0-9-]{2,8}(\/([\w#!:.?+=&%@!\-\/])*)?/`, link)
	return
}

// func handleDatabaseError(err error) error {
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			responder.Error(w, http.StatusInternalServerError, )
// 			return ErrDataRetrieved
// 		}
// 		responder.Error(w, http.StatusInternalServerError, )
// 		panic(err)
// 	}
// }

// func validateDomain(domain string) (valid bool) {
// return
// }
// func validateDomain(domain string) (valid bool) {
// return
// }
// func validateDomain(domain string) (valid bool) {
// return
// }
// func validateDomain(domain string) (valid bool) {
// 	return
// }

func validateRoomForm(form *models.FormPattern) (valid bool) {
	for _, field := range form.Fields {
		if field.Key == "" ||
			field.Length < 0 {
			return
		}
		// unique key name
		for _, field2 := range form.Fields {
			if field.Key == field2.Key {
				return
			}
		}

		switch field.Type {
		case rules.TextField:
			// nothing

		case rules.DateField:
			for _, v := range field.Items {
				if _, err := strconv.ParseInt(v, 10, 64); err != nil {
					return
				}
			}

		case rules.EmailField:
			for _, v := range field.Items {
				if !validateEmail(v) {
					return
				}
			}

		case rules.LinkField:
			for _, v := range field.Items {
				if !validateLink(v) {
					return
				}
			}

		case rules.NumericField:
			for _, v := range field.Items {
				if _, err := strconv.Atoi(v); err != nil {
					return
				}
			}

		default:
			return
		}

	}
	return true
}
