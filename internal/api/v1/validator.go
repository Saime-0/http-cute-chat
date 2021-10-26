package v1

import "github.com/saime-0/http-cute-chat/internal/api/rules"

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
func validateAppSettings(settings string) (valid bool) {

	return len(settings) <= rules.AppSettingsMaxLength

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
