package healer

func (h *Healer) Log(document interface{}) (fail bool) {
	return h.MonitorLogger(h.logger.Log(document))
}

//func (h *Healer) Emergency(document interface{}) (fail bool) {
//	return h.MonitorLogger(h.logger.Log(document))
//}
//func (h *Healer) Alert(document interface{}) (err error) {
//	h.logger.Log()
//}
//func (h *Healer) Critical(document interface{}) (err error) {
//	h.logger.Log()
//}
//func (h *Healer) Error(document interface{}) (err error) {
//	h.logger.Log()
//}
//func (h *Healer) Warning(document interface{}) (err error) {
//	h.logger.Log()
//}
//func (h *Healer) Notice(document interface{}) (err error) {
//	h.logger.Log()
//}
//func (h *Healer) Info(document interface{}) (err error) {
//	h.logger.Log()
//}
//func (h *Healer) Debug(document interface{}) (err error) {
//	h.logger.Log()
//}
