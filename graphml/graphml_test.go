package graphml

import (
	"testing"
	"strconv"
	"reflect"
	"fmt"
)

func TestNewGraphML(t *testing.T) {
	description := "test"
	gml := NewGraphML(description)

	if gml == nil {
		t.Error("gml == nil ")
		return
	}
	if gml.Description != description {
		t.Error("gml.desc != description", gml.Description)
	}
}

func TestGraphML_CreateGraph(t *testing.T) {
	description := "test graph"
	gml := NewGraphML("")

	// test normal creation
	gr, err := gml.CreateGraph(description, EdgeDirectionDirected)
	if err != nil {
		t.Error(err)
		return
	}
	if gr == nil {
		t.Error("gr == nil")
		return
	}
	if gr.Description != description {
		t.Error("gr.Description != description", gr.Description)
	}
	if gr.edgeDefault != "directed" {
		t.Error("gr.edgeDefault != directed", gr.edgeDefault)
	}
	if len(gml.graphs) != 1 {
		t.Error("len(gml.graphs) != 1", len(gml.graphs))
	}

	// test error
	gr, err = gml.CreateGraph(description, EdgeDirectionDefault)
	if err == nil {
		t.Error("error must be raised when default edge direction not provided")
	}
}

func TestGraphML_RegisterKey(t *testing.T) {
	description := "graphml"
	gml := NewGraphML(description)

	keyName := "weight"
	keyDescr := "the weight function"
	keyDefault := 100.0
	_, err := gml.RegisterKey(KeyForNode, keyName, keyDescr, reflect.TypeOf(keyDefault).Kind(), keyDefault)
	if err != nil {
		t.Error(err)
		return
	}
	if len(gml.keys) != 1 {
		t.Error("len(gml.keys) != 1", len(gml.keys))
		return
	}

	if gml.keys[0].name != keyName {
		t.Error("gml.keys[0].name != keyName", gml.keys[0].name)
	}
	if gml.keys[0].Description != keyDescr {
		t.Error("gml.keys[0].Description != keyDescr", gml.keys[0].Description)
	}
	if gml.keys[0].keyType != "double" {
		t.Error("gml.keys[0].keyType != double", gml.keys[0].keyType)
	}
	if val, err := strconv.ParseFloat(gml.keys[0].defaultValue, 64); err != nil {
		t.Error(err)
	} else if keyDefault != val {
		t.Error("keyDefault != val", keyDefault, val)
	}
	if gml.keys[0].target != KeyForNode {
		t.Error("gml.keys[0].keyFor != KeyForNode", gml.keys[0].target)
	}
	if gml.keys[0].ID != "d0" {
		t.Error("gml.keys[0].ID != d0", gml.keys[0].ID)
	}

	// register key with the same standard identifier and check that error raised
	_, err = gml.RegisterKey(KeyForNode, keyName, keyDescr, reflect.TypeOf(keyDefault).Kind(), keyDefault + 100)
	if err == nil {
		t.Error("error should be raised when attempting to register Key with already existing standard identifier")
	}

	// register another key and check ID
	keyName = "height"
	_, err = gml.RegisterKey(KeyForNode, keyName, keyDescr, reflect.TypeOf(keyDefault).Kind(), keyDefault + 100)
	if err != nil {
		t.Error(err)
		return
	}
	if len(gml.keys) != 2 {
		t.Error("len(gml.keys) != 2", len(gml.keys))
		return
	}
	if gml.keys[1].ID != "d1" {
		t.Error("gml.keys[1].ID != d1", gml.keys[1].ID)
	}
}

func TestGraph_AddNode(t *testing.T) {
	description := "test graph"
	gml := NewGraphML("")
	gr, err := gml.CreateGraph(description, EdgeDirectionDirected)
	if err != nil {
		t.Error(err)
		return
	}

	// add node
	attributes := make(map[string]interface{})
	attributes["double"] = 100.1
	attributes["bool"] = false
	attributes["integer"] = 120
	attributes["string"] = "string data"

	description = "test node"
	node, err := gr.AddNode(attributes, description)
	if err != nil {
		t.Error(err)
		return
	}

	// test results
	if len(gr.nodes) != 1 {
		t.Error("len(gr.nodes) != 1", len(gr.nodes))
	}
	if len(gr.parent.keys) != 4 {
		t.Error("len(gr.parent.keys) != 4", len(gr.parent.keys))
	}
	if len(node.data) != 4 {
		t.Error("len(node.data) != 4", len(node.data))
	}
	count := 0
	for name, val := range attributes {
		keyNameId := keyIdentifier(name, KeyForNode)
		if key := gr.parent.GetKey(name, KeyForNode); key == nil {
			t.Error("failed to find Key with standard identifier:", keyNameId)
			return
		} else {
			// check if attribute data value was saved
			for i, node := range node.data {
				if node.key == key.ID {
					str_val := fmt.Sprint(val)
					if node.value != str_val {
						t.Error("node.value != str_val at:", i)
					}
					// increment counter to count this attribute
					count++
				}
			}
		}
	}
	if count != len(attributes) {
		t.Error("failed to check all nodes")
	}
}

func TestGraphML_stringValueIfSupported(t *testing.T) {
	testBool := true
	res, err := stringValueIfSupported(testBool, "boolean")
	if err != nil {
		t.Error(err)
		return
	}
	if pres, err := strconv.ParseBool(res); err != nil {
		t.Error(err)
	} else if testBool != pres {
		t.Error("testBool != pres", testBool, pres)
	}

	testInt := 42
	res, err = stringValueIfSupported(testInt, "int")
	if err != nil {
		t.Error(err)
		return
	}
	if pres, err := strconv.ParseInt(res, 10, 32); err != nil {
		t.Error(err)
	} else if testInt != int(pres) {
		t.Error("testInt != pres", testInt, pres)
	}

	testLong := int64(12993888475775)
	res, err = stringValueIfSupported(testLong, "long")
	if err != nil {
		t.Error(err)
		return
	}
	if pres, err := strconv.ParseInt(res, 10, 64); err != nil {
		t.Error(err)
	} else if testLong != pres {
		t.Error("testLong != pres", testLong, pres)
	}

	testFloat := float32(0.5)
	res, err = stringValueIfSupported(testFloat, "float")
	if err != nil {
		t.Error(err)
		return
	}
	if pres, err := strconv.ParseFloat(res, 32); err != nil {
		t.Error(err)
	} else if testFloat != float32(pres) {
		t.Error("testFloat != pres", testFloat, pres)
	}

	testDouble := float64(10000.552)
	res, err = stringValueIfSupported(testDouble, "double")
	if err != nil {
		t.Error(err)
		return
	}
	if pres, err := strconv.ParseFloat(res, 64); err != nil {
		t.Error(err)
	} else if testDouble != pres {
		t.Error("testDouble != pres", testDouble, pres)
	}

	testString := "test string"
	res, err = stringValueIfSupported(testString, "string")
	if err != nil {
		t.Error(err)
		return
	}
	if testString != res {
		t.Error("testString != res", testString, res)
	}
}
