package srcedit_test

import (
	"testing"

	"github.com/gsxhq/gsxui/internal/srcedit"
)

func TestApplyReplacesSpansFromTheEndSoOffsetsStayValid(t *testing.T) {
	src := []byte("aaa bbb ccc")
	got, err := srcedit.Apply(src, []srcedit.Edit{
		{Start: 0, End: 3, Value: "A"},
		{Start: 8, End: 11, Value: "CCCC"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "A bbb CCCC" {
		t.Fatalf("got %q", got)
	}
	if string(src) != "aaa bbb ccc" {
		t.Fatal("Apply mutated its input")
	}
}

func TestApplyRejectsOverlapAndOutOfRange(t *testing.T) {
	src := []byte("abcdef")
	if _, err := srcedit.Apply(src, []srcedit.Edit{{Start: 0, End: 3, Value: ""}, {Start: 2, End: 4, Value: ""}}); err == nil {
		t.Fatal("overlapping edits accepted")
	}
	if _, err := srcedit.Apply(src, []srcedit.Edit{{Start: 4, End: 9, Value: ""}}); err == nil {
		t.Fatal("out-of-range edit accepted")
	}
	if _, err := srcedit.Apply(src, []srcedit.Edit{{Start: 3, End: 2, Value: ""}}); err == nil {
		t.Fatal("inverted span accepted")
	}
}

func TestApplyWithNoEditsReturnsAnEqualCopy(t *testing.T) {
	src := []byte("unchanged")
	got, err := srcedit.Apply(src, nil)
	if err != nil || string(got) != "unchanged" {
		t.Fatalf("got %q, %v", got, err)
	}
}
