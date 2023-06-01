package gotempmail

type TempMail struct {
	email string
}

func New() *TempMail {
	ret := TempMail{}
	return &ret
}
