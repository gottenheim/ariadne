package card

type FakeCardBuilder struct {
	section    string
	entry      string
	artifacts  []CardArtifact
	activities CardActivity
}

func NewFakeCard() *FakeCardBuilder {
	return &FakeCardBuilder{}
}

func (b *FakeCardBuilder) WithSection(section string) *FakeCardBuilder {
	b.section = section
	return b
}

func (b *FakeCardBuilder) WithEntry(entry string) *FakeCardBuilder {
	b.entry = entry
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

func (b *FakeCardBuilder) WithActivityChain(activities CardActivity) *FakeCardBuilder {
	b.activities = activities
	return b
}

func (b *FakeCardBuilder) Build() *Card {
	return FromExisting(b.section, b.entry, b.artifacts, b.activities)
}

func ExtractBriefCards(cards []*Card) []BriefCard {
	var result []BriefCard
	for _, card := range cards {
		result = append(result, BriefCard{
			Section:    card.Section(),
			Entry:      card.Entry(),
			Activities: card.Activities(),
		})
	}
	return result
}
