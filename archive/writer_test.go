package archive

import (
	"bytes"
	"testing"

	"github.com/gottenheim/ariadne/fs"
	"github.com/spf13/afero"
)

func TestCompressingAndDecompressingArchive(t *testing.T) {
	fakeFs, err := fs.NewFake([]fs.FakeEntry{
		fs.NewFakeEntry("/dev/frontend", "deploy.yaml", `first: firstValue`),
		fs.NewFakeEntry("/dev/frontend", "app.yaml", `second: secondValue`),
		fs.NewFakeEntry("/dev/backend", "deploy.yaml", `third: thirdValue`),
		fs.NewFakeEntry("/dev/backend", "app.yaml", `fourth: fourthValue`),
	})

	if err != nil {
		t.Fatal(err)
	}

	writer := NewWriter()

	err = writer.AddDir(fakeFs, "/dev")

	if err != nil {
		t.Fatal(err)
	}

	buffer, err := writer.Buffer()

	if err != nil {
		t.Fatal(err)
	}

	reader := bytes.NewReader(buffer.Bytes())

	targetFs := afero.NewMemMapFs()

	err = Uncompress(reader, targetFs, "/layers")

	if err != nil {
		t.Fatal(err)
	}

	first, err := afero.ReadFile(targetFs, "/layers/frontend/deploy.yaml")

	if err != nil {
		t.Fatal(err)
	}

	if string(first) != "first: firstValue" {
		t.Fatal("Wrong first file contents")
	}

	second, err := afero.ReadFile(targetFs, "/layers/frontend/app.yaml")

	if err != nil {
		t.Fatal(err)
	}

	if string(second) != "second: secondValue" {
		t.Fatal("Wrong second file contents")
	}

	third, err := afero.ReadFile(targetFs, "/layers/backend/deploy.yaml")

	if err != nil {
		t.Fatal(err)
	}

	if string(third) != "third: thirdValue" {
		t.Fatal("Wrong third file contents")
	}

	fourth, err := afero.ReadFile(targetFs, "/layers/backend/app.yaml")

	if err != nil {
		t.Fatal(err)
	}

	if string(fourth) != "fourth: fourthValue" {
		t.Fatal("Wrong fourth file contents")
	}
}
