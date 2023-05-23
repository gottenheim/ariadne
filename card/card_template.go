package card

type CardTemplate struct {
	artifacts []CardArtifact
}

func NewCardTemplate(artifacts []CardArtifact) *CardTemplate {
	return &CardTemplate{
		artifacts: artifacts,
	}
}

func (t *CardTemplate) Artifacts() []CardArtifact {
	return t.artifacts
}

func (t *CardTemplate) FindArtifactByName(name string) *CardArtifact {
	for _, artifact := range t.artifacts {
		if artifact.Name() == name {
			return &artifact
		}
	}
	return nil
}
