package visual

import "testing"

func TestSelectViewportsRejectsUnknown(t *testing.T) {
	if _, err := SelectViewports("desktop,watch"); err == nil {
		t.Fatal("unknown viewport was accepted")
	}
	got, err := SelectViewports("mobile,desktop")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Name != "mobile" || got[1].Width != 1440 {
		t.Fatalf("%+v", got)
	}
}
