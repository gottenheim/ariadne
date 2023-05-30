package card

type CardTemplateRepository interface {
	GetTemplate() (*CardTemplate, error)
}
