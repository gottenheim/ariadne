package card

type CardArtifact struct {
	name    string
	content []byte
}

func NewCardArtifact(name string, content []byte) CardArtifact {
	return CardArtifact{
		name:    name,
		content: content,
	}
}

func (a *CardArtifact) Name() string {
	return a.name
}

func (a *CardArtifact) Content() []byte {
	return a.content
}
