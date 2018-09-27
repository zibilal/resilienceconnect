package resilienceconnect

import (
	"testing"
	"fmt"
)

func TestConnectorOption(t *testing.T) {
	t.Log("Testing ConnectorOption")
	{
		cOptions := make(ConnectorOption)
		cOptions.Put("app", "testing1")
		cOptions.Put("version", "v0.0.1")
		cOptions.Put("equal", "data_equal")

		expected := "map[app:testing1 version:v0.0.1 equal:data_equal]"

		if fmt.Sprintf("%v", cOptions) == expected {
			t.Logf("%s expected optinos %s", success, expected)
		} else {
			t.Fatalf("%s expected options %s, got %v", failed, expected, cOptions)
		}
	}
}
