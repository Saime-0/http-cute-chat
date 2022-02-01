package kit

var _ error = (*Errorus)(nil)

type Errorus struct {
	Err error
}

func (e *Errorus) Error() string {
	return e.Err.Error()
}

func (e *Errorus) Fail(err error) (fail bool) {
	return e.FailTh(nil, err)
}

func (e *Errorus) FailTh(_ interface{}, err error) (fail bool) {
	if err != nil {
		e.Err = err
		return true
	}
	return
}

func CreateErrorHandler() *Errorus {
	return &Errorus{}
}
