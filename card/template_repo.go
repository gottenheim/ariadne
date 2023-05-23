package card

type TemplateRepository interface {
	GetTemplate() (*CardTemplate, error)
}
