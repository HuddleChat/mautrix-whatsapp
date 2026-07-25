package connector

import (
	"testing"

	"go.mau.fi/whatsmeow/types"
)

const safeDisplaynameTemplate = `{{or .BusinessName .PushName .Phone .RedactedPhone "Unknown user"}}`

func newTestConfig(t *testing.T, privateChatNameTemplate string) *Config {
	t.Helper()
	cfg := &Config{
		DisplaynameTemplate:     safeDisplaynameTemplate,
		PrivateChatNameTemplate: privateChatNameTemplate,
	}
	if err := cfg.PostProcess(); err != nil {
		t.Fatalf("PostProcess failed: %v", err)
	}
	return cfg
}

func testJID() types.JID {
	return types.JID{User: "15551234567", Server: types.DefaultUserServer}
}

// The whole point of the separate template: the room name may use address-book fields,
// while the (global) ghost displayname must not.
func TestPrivateChatNamePrefersAddressBookName(t *testing.T) {
	cfg := newTestConfig(t, `{{or .FullName .FirstName .BusinessName .PushName}}`)
	contact := types.ContactInfo{FullName: "Brother", PushName: "mo"}

	if got := cfg.FormatPrivateChatName(testJID(), "", contact); got != "Brother" {
		t.Errorf("private chat name = %q, want %q", got, "Brother")
	}
	if got := cfg.FormatDisplayname(testJID(), "", contact); got != "mo" {
		t.Errorf("ghost displayname = %q, want %q (address-book name must not leak)", got, "mo")
	}
}

// An empty render must stay empty so wrapDMInfo leaves ChatInfo.Name nil. Returning a
// pointer to "" instead would set NameIsCustom and permanently blank the room name.
func TestPrivateChatNameEmptyWhenNoAddressBookEntry(t *testing.T) {
	cfg := newTestConfig(t, `{{or .FullName .FirstName ""}}`)

	if got := cfg.FormatPrivateChatName(testJID(), "", types.ContactInfo{PushName: "mo"}); got != "" {
		t.Errorf("private chat name = %q, want empty so the ghost name is kept", got)
	}
}

func TestPrivateChatNameDisabledByDefault(t *testing.T) {
	cfg := newTestConfig(t, "")

	if got := cfg.FormatPrivateChatName(testJID(), "", types.ContactInfo{FullName: "Brother"}); got != "" {
		t.Errorf("private chat name = %q, want empty when no template is configured", got)
	}
}

func TestPrivateChatNameFallbackOrder(t *testing.T) {
	cfg := newTestConfig(t, `{{or .FullName .FirstName .BusinessName .PushName .Phone}}`)

	for _, tc := range []struct {
		name    string
		contact types.ContactInfo
		want    string
	}{
		{"full name wins", types.ContactInfo{FullName: "Full", FirstName: "First", PushName: "Push"}, "Full"},
		{"first name next", types.ContactInfo{FirstName: "First", PushName: "Push"}, "First"},
		{"business before push", types.ContactInfo{BusinessName: "Biz", PushName: "Push"}, "Biz"},
		{"push name next", types.ContactInfo{PushName: "Push"}, "Push"},
		{"phone last", types.ContactInfo{}, "+15551234567"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := cfg.FormatPrivateChatName(testJID(), "", tc.contact); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

// A DM left without an explicit name keeps following the global ghost, so any other user's
// contact sync can rename it. The deployed template must therefore always render something
// for a phone-number JID, even with no contact info at all.
func TestDeployedTemplateNeverRendersEmptyForPhoneJID(t *testing.T) {
	cfg := newTestConfig(t, `{{or .FullName .FirstName .BusinessName .PushName .Phone}}`)

	if got := cfg.FormatPrivateChatName(testJID(), "", types.ContactInfo{}); got == "" {
		t.Error("template rendered empty with no contact info; the room would stay ghost-named")
	}
}

func TestInvalidPrivateChatNameTemplateIsRejected(t *testing.T) {
	cfg := &Config{
		DisplaynameTemplate:     safeDisplaynameTemplate,
		PrivateChatNameTemplate: `{{.Nope`,
	}
	if err := cfg.PostProcess(); err == nil {
		t.Error("PostProcess accepted an invalid private chat name template")
	}
}
