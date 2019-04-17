package mdconv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: fix the behavior for this test
func TestToHTML(t *testing.T) {
	assertions := []struct {
		name   string
		mdtext string
		exp    string
	}{
		{
			"Simple test",
			"This _i_s a t*est* __*string*__ in [which](we) ***will*** ___try__ `to` `**list**` _`all`_ (possible) [cases] (.)",
			"This <em>i</em>s a t<em>est</em> <strong><em>string</em></strong> in <a href=\"we\" target=\"_blank\">which</a> <em></em><em>will</em><em></em> <strong>_try</strong> <font face=\"Courier New\">to</font> <font face=\"Courier New\"><strong>list</strong></font> <em><font face=\"Courier New\">all</font></em> (possible) [cases] (.)",
		},
	}

	for i, a := range assertions {
		actual := ToHTML(a.mdtext)
		assert.Equal(t, a.exp, actual, `%d) Test "%s": expected "%s" not equal \t"%s"`, i+1, a.name, a.exp, actual)
	}
}
