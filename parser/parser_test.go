package parser

import (
	"strings"
	"testing"
)

func TestTranspile(t *testing.T) {
	input := `
@Component
script Unit extends Component

	@Sync
	property number hp = 500

	property string name = "test"

	@ExecSpace("Server")
	method void Spawn(BuffedUnitModel unitData, table skillData)
		self.id = unitData.id
	end
`
	expected := `
Unit = {}

	Unit.hp = 500

	Unit.name = "test"

	-- @method_signature: Spawn(unitData:BuffedUnitModel,skillData:table)
	function Unit:Spawn(unitData, skillData)
		self.id = unitData.id
	end
`
	output, err := Transpile(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.TrimSpace(output) != strings.TrimSpace(expected) {
		t.Errorf("expected:\n%s\n\ngot:\n%s", expected, output)
	}
}
