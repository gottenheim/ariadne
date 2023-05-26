package card

type FakeCardBuilder struct {
	key        Key
	artifacts  []CardArtifact
	activities CardActivity
}

func NewFakeCard() *FakeCardBuilder {
	return &FakeCardBuilder{}
}

func (b *FakeCardBuilder) WithKey(key Key) *FakeCardBuilder {
	b.key = key
	return b
}

func (b *FakeCardBuilder) WithArtifact(name string, text string) *FakeCardBuilder {
	b.artifacts = append(b.artifacts, NewCardArtifact(name, []byte(text)))
	return b
}

func (b *FakeCardBuilder) WithActivities(activities ...GenerateActivity) *FakeCardBuilder {
	b.activities = GenerateActivityChain(activities...)
	return b
}

func (b *FakeCardBuilder) Build() *Card {
	card := NewCard(b.key, b.artifacts)
	card.SetArtifacts(b.artifacts)
	card.SetActivities(b.activities)
	return card
}

func ExtractKeysWithActivities(cards []*Card) []*KeyWithActivities {
	var result []*KeyWithActivities
	for _, card := range cards {
		result = append(result, &KeyWithActivities{
			Key:        card.Key(),
			Activities: card.Activities(),
		})
	}
	return result
}
