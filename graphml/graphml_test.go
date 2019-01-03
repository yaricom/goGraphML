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

func TestGraphML_AddGraph(t *testing.T) {
	description := "test graph"
	gml := NewGraphML("")

	// test normal creation
	attributes := make(map[string]interface{})
	attributes["double"] = 100.1
	attributes["bool"] = false
	attributes["integer"] = 120
	attributes["string"] = "string data"

	graph, err := gml.AddGraph(description, EdgeDirectionDirected, attributes)
	if err != nil {
		t.Error(err)
		return
	}
	if graph == nil {
		t.Error("gr == nil")
		return
	}
	if graph.Description != description {
		t.Error("gr.Description != description", graph.Description)
	}
	if graph.edgeDefault != "directed" {
		t.Error("gr.edgeDefault != directed", graph.edgeDefault)
	}
	if len(gml.graphs) != 1 {
		t.Error("len(gml.graphs) != 1", len(gml.graphs))
	}

	// check attributes processing
	if len(graph.parent.keys) != 4 {
		t.Error("len(gr.parent.keys) != 4", len(graph.parent.keys))
	}
	if len(graph.data) != 4 {
		t.Error("len(node.data) != 4", len(graph.data))
	}

	// check attributes
	checkAttributes(attributes, graph.data, KeyForGraph, graph.parent, t)

	// test error
	graph, err = gml.AddGraph(description, EdgeDirectionDefault, nil)
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
	gr, err := gml.AddGraph(description, EdgeDirectionDirected, nil)
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

	// check attributes
	checkAttributes(attributes, node.data, KeyForNode, gr.parent, t)
}

func TestGraph_AddEdge(t *testing.T) {
	description := "test graph"
	gml := NewGraphML("")
	gr, err := gml.AddGraph(description, EdgeDirectionDirected, nil)
	if err != nil {
		t.Error(err)
		return
	}
	// Add nodes
	var n1, n2 *Node
	if n1, err = gr.AddNode(nil, "#1"); err != nil {
		t.Error(err)
		return
	}
	if n2, err = gr.AddNode(nil, "#2"); err != nil {
		t.Error(err)
		return
	}

	// Add graph
	attributes := make(map[string]interface{})
	attributes["double"] = 100.1
	attributes["bool"] = false
	attributes["integer"] = 120
	attributes["string"] = "string data"

	description = "test edge"
	edge, err := gr.AddEdge(n1, n2, attributes, EdgeDirectionDefault, description)
	if err != nil {
		t.Error(err)
		return
	}

	// test results
	if len(gr.edges) != 1 {
		t.Error("len(gr.edges) != 1", len(gr.edges))
	}
	if len(gr.edgesMap) != 1 {
		t.Error("len(gr.edgesMap) != 1", len(gr.edgesMap))
	}
	if edge.Description != description {
		t.Error("edge.Description != description", edge.Description)
	}
	if edge.source != n1.ID {
		t.Error("edge.source != n1.ID ", edge.source, n1.ID)
	}
	if edge.target != n2.ID {
		t.Error("edge.target != n2.ID", edge.target, n2.ID)
	}
	if len(edge.directed) != 0 {
		t.Error("len(edge.directed) != 0", len(edge.directed))
	}
	if _, ok := gr.edgesMap[edgeIdentifier(n1.ID, n2.ID)]; !ok {
		t.Error("edge not found in edges map")
	}

	// check attributes
	checkAttributes(attributes, edge.data, KeyForEdge, gr.parent, t)

	// check error by adding the same node
	edge, err = gr.AddEdge(n1, n2, attributes, EdgeDirectionDefault, description)
	if err == nil {
		t.Error("error must be raised when add same edge")
		return
	}

	// check no error for directed node backward
	edge, err = gr.AddEdge(n2, n1, attributes, EdgeDirectionDefault, description)
	if err != nil {
		t.Error(err)
		return
	}
}

func checkAttributes(attributes map[string]interface{}, data_holders []*Data, target KeyForElement, gml *GraphML, t *testing.T) {
	count := 0
	for name, val := range attributes {
		keyNameId := keyIdentifier(name, target)
		if key := gml.GetKey(name, target); key == nil {
			t.Error("failed to find Key with standard identifier:", keyNameId)
			return
		} else {
			// check if attribute data value was saved
			for i, data := range data_holders {
				if data.key == key.ID {
					str_val := fmt.Sprint(val)
					if data.value != str_val {
						t.Error("data.value != str_val at:", i)
					}
					// increment counter to count this attribute
					count++
				}
			}
		}
	}
	if count != len(attributes) {
		t.Error("failed to check all attributes")
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
