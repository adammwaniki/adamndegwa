package content

import "testing"

// TestRealContentLoads guards the actual content/ directory against
// malformed frontmatter or duplicate/missing data.
func TestRealContentLoads(t *testing.T) {
	withContentDir(t, "../../content")
	for _, s := range Sections() {
		if len(s.Cards) == 0 {
			t.Errorf("section %q has no cards", s.ID)
		}
		for _, c := range s.Cards {
			if c.Title == "" || c.Description == "" || c.Date == "" {
				t.Errorf("%s/%s: missing required frontmatter: %+v", s.ID, c.ID, c)
			}
			if c.Detail == "" {
				t.Errorf("%s/%s: empty rendered detail", s.ID, c.ID)
			}
		}
	}
}
