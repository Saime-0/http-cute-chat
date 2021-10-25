package v1

func validateDomain(domain string) (valid bool) {
	if (len(domain) > DomainMaxLength) || (len(domain) < DomainMinLength) {
		return false
	}
	return true
}
func validateName(name string) (valid bool) {
	if (len(name) > NameMaxLength) || (len(name) < NameMinLength) {
		return false
	}
	return true
}
func validateAppSettings(settings string) (valid bool) {

	return len(settings) <= AppSettingsMaxLength

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
