package config

import (
	"testing"
)

func TestExpander_EnvVarExpanded(t *testing.T) {
	//testCases := []struct {
	//	name string // test case name (also file name containing config yaml)
	//}{
	//	{name: "no-env.yaml"},
	//	{name: "partial-env.yaml"},
	//	{name: "all-env.yaml"},
	//}
	//
	//const valueExtra = "some string"
	//const valueExtraMapValue = "some map value"
	//const valueExtraListMapValue = "some list map value"
	//const valueExtraListElement = "some list value"
	//
	//t.Setenv("EXTRA", valueExtra)
	//t.Setenv("EXTRA_MAP_VALUE_1", valueExtraMapValue+"_1")
	//t.Setenv("EXTRA_MAP_VALUE_2", valueExtraMapValue+"_2")
	//t.Setenv("EXTRA_LIST_MAP_VALUE_1", valueExtraListMapValue+"_1")
	//t.Setenv("EXTRA_LIST_MAP_VALUE_2", valueExtraListMapValue+"_2")
	//t.Setenv("EXTRA_LIST_VALUE_1", valueExtraListElement+"_1")
	//t.Setenv("EXTRA_LIST_VALUE_2", valueExtraListElement+"_2")
	//
	//for _, test := range testCases {
	//	t.Run(test.name, func(t *testing.T) {
	//		expander := Expander{}
	//
	//		// modifiedConfig, err := expander.Expand()
	//	})
	//}
}
