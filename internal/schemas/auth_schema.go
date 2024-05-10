package schemas

type RegisterDataForErr struct {
	Name            string
	Email           string
	Password        string
	ConfirmPassword string
}

type RegisterErrors struct {
	Name            string
	Email           string
	Password        string
	ConfirmPassword string
}

type TemplateForm struct {
	RegisterErrors     RegisterErrors
	RegisterDataForErr RegisterDataForErr
}

type RegisterForm struct {
	TemplateForm *TemplateForm
}
