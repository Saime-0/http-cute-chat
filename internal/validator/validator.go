package validator

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/res"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"regexp"
	"strconv"
)

func ValidateDomain(domain string) (valid bool) {
	if (len(domain) > rules.DomainMaxLength) || (len(domain) < rules.DomainMinLength) {
		return false
	}
	return true
}
func ValidateName(name string) (valid bool) {
	if (len(name) > rules.NameMaxLength) || (len(name) < rules.NameMinLength) {
		return false
	}
	return true
}
func ValidateNote(note string) (valid bool) {
	if len(note) > rules.NoteMaxLength {
		return false
	}
	return true
}
func ValidatePassword(password string) (valid bool) { // todo

	return len(password) <= rules.MaxPasswordLength && len(password) >= rules.MinPasswordLength
}
func ValidateEmail(email string) (valid bool) { // todo
	valid, _ = regexp.MatchString(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`, email)
	return
}
func ValidateLifetime(lt int64) (valid bool) {
	return lt >= 60 && lt <= res.Year
}
func ValidateAliens(aliens int) (valid bool) {

	return aliens >= 1 && aliens <= 99_999
}
func ValidateLink(link string) (valid bool) {
	valid, _ = regexp.MatchString(`/^((ftp|http|https):\/\/)?(www\.)?([A-Za-zА-Яа-я0-9]{1}[A-Za-zА-Яа-я0-9\-]*\.?)*\.{1}[A-Za-zА-Яа-я0-9-]{2,8}(\/([\w#!:.?+=&%@!\-\/])*)?/`, link)
	return
}
func ValidateOffset(offset int) (valid bool) {
	return offset >= 0
}
func ValidateLimit(limit int) (valid bool) {
	return limit >= 1 && limit <= 20
}
func ValidateNameFragment(fragment string) (valid bool) {
	return len(fragment) >= 1 && len(fragment) <= rules.NameMaxLength
}
func ValidateID(id int) (valid bool) {
	return id > 0
}
func ValidateIDs(ids []int) (valid bool) {
	for _, id := range ids {
		if !ValidateID(id) {
			return false
		}
	}
	return true
}
func ValidateAllowInput(allow *model.AllowInput) (valid bool) {
	_, err := strconv.Atoi(allow.Value)
	if allow.Group == model.AllowGroupRole || allow.Group == model.AllowGroupMember {
		if err != nil {
			return false
		}
	} else if allow.Group == model.AllowGroupChar {
		if model.CharTypeModer.String() != allow.Value &&
			model.CharTypeAdmin.String() != allow.Value {
			return false
		}
	}
	return true
}
func ValidateAllowsInput(allows *model.AllowsInput) (valid bool) {
	for _, v := range allows.Allows {
		if !ValidateAllowInput(v) {
			return
		}
	}
	return true
}

func ValidateRoomForm(form *model.UpdateFormInput) (valid bool, err error) {
	if len(form.Fields) > rules.MaxFormFields {
		return false, cerrors.New("превышен лимит полей")
	}
	for _, field := range form.Fields {
		// todo valid name
		if !ValidateName(field.Key) {
			return false, cerrors.New("невалидное значение ключа")
		}
		if len(field.Items) > rules.MaxFielditems {
			return false, cerrors.New("exceeded the limit of items")
		}
		// unique key name
		count := 0
		for _, field2 := range form.Fields {
			if field.Key == field2.Key {
				count += 1
			}
		}
		if count > 1 {
			return false, cerrors.New("повторяющиеся значения ключей")
		}
		if field.Length != nil && *field.Length < 1 {
			return false, cerrors.New("длина не может быть меньше 1")
		}

		// handling fields type
		switch field.Type {
		case model.FieldTypeText:
			// nothing

		case model.FieldTypeDate:
			for _, v := range field.Items {
				if _, err := strconv.ParseInt(v, 10, 64); err != nil {
					return false, cerrors.New("it is not possible to convert a item value to int64")
				}
			}

		case model.FieldTypeEmail:
			for _, v := range field.Items {
				if !ValidateEmail(v) {
					return false, cerrors.New("item is not email")

				}
			}

		case model.FieldTypeLink:
			for _, v := range field.Items {
				if !ValidateLink(v) {
					return false, cerrors.New("item is not link")

				}
			}

		case model.FieldTypeNumeric:
			for _, v := range field.Items {
				if _, err := strconv.Atoi(v); err != nil {
					return false, cerrors.New("item is not numeric type")

				}
			}

		default:
			return
		}

	}
	return true, nil
}

func ValidateSessionKey(sessionKey string) (valid bool) {
	return regexp.MustCompile(`^[a-zA-Z0-9\-\=]{20}$`).MatchString(sessionKey)
}
