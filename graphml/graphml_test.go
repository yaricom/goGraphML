package graphml

import (
	"testing"
	"strconv"
	"reflect"
	"fmt"
	"bytes"
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

func TestGraphML_Encode(t *testing.T) {
	// build GraphML
	description := "test graph"
	gml := NewGraphML("TestGraphML_Encode")

	// register common data-function for all elements
	if _, err := gml.RegisterKey(KeyForAll, "double", "common double data-function", reflect.Float64, 10.2); err != nil {
		t.Error(err)
		return
	}

	attributes := make(map[string]interface{})
	attributes["double"] = NotAValue
	attributes["bool"] = false
	attributes["integer"] = 120
	attributes["string"] = "string data"
	// add graph
	graph, err := gml.AddGraph(description, EdgeDirectionDirected, attributes)
	if err != nil {
		t.Error(err)
		return
	}
	// add nodes
	description = "test node #1"
	n1, err := graph.AddNode(attributes, description)
	if err != nil {
		t.Error(err)
		return
	}
	description = "test node #2"
	n2, err := graph.AddNode(attributes, description)
	if err != nil {
		t.Error(err)
		return
	}
	// add edge
	description = "test edge"
	_, err = graph.AddEdge(n1, n2, attributes, EdgeDirectionDefault, description)
	if err != nil {
		t.Error(err)
		return
	}

	// encode
	out_buf := bytes.NewBufferString("")
	err = gml.Encode(out_buf, false)

	// check results
	res_string := "<graphml><desc>TestGraphML_Encode</desc><key id=\"d0\" for=\"all\" attr.name=\"double\" attr.type=\"double\"><desc>common double data-function</desc><default>10.2</default></key><key id=\"d1\" for=\"graph\" attr.name=\"bool\" attr.type=\"boolean\"></key><key id=\"d2\" for=\"graph\" attr.name=\"integer\" attr.type=\"int\"></key><key id=\"d3\" for=\"graph\" attr.name=\"string\" attr.type=\"string\"></key><key id=\"d4\" for=\"node\" attr.name=\"bool\" attr.type=\"boolean\"></key><key id=\"d5\" for=\"node\" attr.name=\"integer\" attr.type=\"int\"></key><key id=\"d6\" for=\"node\" attr.name=\"string\" attr.type=\"string\"></key><key id=\"d7\" for=\"edge\" attr.name=\"bool\" attr.type=\"boolean\"></key><key id=\"d8\" for=\"edge\" attr.name=\"integer\" attr.type=\"int\"></key><key id=\"d9\" for=\"edge\" attr.name=\"string\" attr.type=\"string\"></key><graph id=\"g0\" edgedefault=\"directed\"><desc>test graph</desc><node id=\"n0\"><desc>test node #1</desc><data key=\"d0\">10.2</data><data key=\"d4\">false</data><data key=\"d5\">120</data><data key=\"d6\">string data</data></node><node id=\"n1\"><desc>test node #2</desc><data key=\"d0\">10.2</data><data key=\"d4\">false</data><data key=\"d5\">120</data><data key=\"d6\">string data</data></node><edge id=\"e0\" source=\"n0\" target=\"n1\"><desc>test edge</desc><data key=\"d0\">10.2</data><data key=\"d7\">false</data><data key=\"d8\">120</data><data key=\"d9\">string data</data></edge><data key=\"d0\">10.2</data><data key=\"d1\">false</data><data key=\"d2\">120</data><data key=\"d3\">string data</data></graph></graphml>"
	if out_buf.Len()!= len(res_string) {
		t.Error("out_buf.String() != res_string")
		t.Log(out_buf.String())
		t.Log(res_string)
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
	if graph.EdgeDefault != "directed" {
		t.Error("gr.edgeDefault != directed", graph.EdgeDefault)
	}
	if len(gml.Graphs) != 1 {
		t.Error("len(gml.graphs) != 1", len(gml.Graphs))
	}

	// check attributes processing
	if len(graph.parent.Keys) != 4 {
		t.Error("len(gr.parent.keys) != 4", len(graph.parent.Keys))
	}
	if len(graph.Data) != 4 {
		t.Error("len(node.data) != 4", len(graph.Data))
	}

	// check attributes
	checkAttributes(attributes, graph.Data, KeyForGraph, graph.parent, t)

	// test error
	graph, err = gml.AddGraph(description, EdgeDirectionDefault, nil)
	if err == nil {
		t.Error("error must be raised when default edge direction not provided")
	}
}

// tests if default value of common key will be used for empty attributes
func TestGraphML_RegisterKeyDefaultValue(t *testing.T) {
	description := "graphml"
	gml := NewGraphML(description)

	commonKeyName := "weight"
	keyDescr := "the weight function"
	keyDefault := 100.2
	c_key, err := gml.RegisterKey(KeyForAll, commonKeyName, keyDescr, reflect.TypeOf(keyDefault).Kind(), keyDefault)
	if err != nil {
		t.Error(err)
		return
	}

	// register graph and test number of keys
	attributes := make(map[string]interface{})
	attributes[commonKeyName] = NotAValue // empty attribute - default value will be used
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

	// check attributes processing
	if len(graph.parent.Keys) != 4 {
		// it is expected tha GraphML has 4 keys: one common and three specific to the Graph element
		t.Error("len(gr.parent.keys) != 4", len(graph.parent.Keys))
	}
	if len(graph.Data) != 4 {
		// it is expected that Graph element has 4 data elements
		t.Error("len(node.data) != 4", len(graph.Data))
	}

	// check if default value of common key was used
	for _, d := range graph.Data {
		if d.ID == c_key.ID && d.Value != c_key.DefaultValue {
			t.Error("d.Value != c_key.DefaultValue", d.Value, c_key.DefaultValue)
		}
	}

	// check attribute without value an without default value
	commonKeyName2 := "height"
	_, err = gml.RegisterKey(KeyForAll, commonKeyName2, keyDescr, reflect.TypeOf(keyDefault).Kind(), nil)
	if err != nil {
		t.Error(err)
		return
	}

	// add node with empty attribute which has no default value
	attributes[commonKeyName2] = NotAValue
	n, err := graph.AddNode(attributes, "test node with empty attribute without default value")
	// no node should be created
	if n != nil {
		t.Error("node should not be added due to attribute error")
	}
	// error must be raised
	if err == nil {
		t.Error("Error must be raised for empty attribute and no default value")
		return
	}

}

func TestGraphML_RegisterKeyForAll(t *testing.T) {
	description := "graphml"
	gml := NewGraphML(description)

	commonKeyName := "weight"
	keyDescr := "the weight function"
	keyDefault := 100.0
	_, err := gml.RegisterKey(KeyForAll, commonKeyName, keyDescr, reflect.TypeOf(keyDefault).Kind(), keyDefault)
	if err != nil {
		t.Error(err)
		return
	}

	// register graph and test number of keys
	attributes := make(map[string]interface{})
	attributes[commonKeyName] = 100.1
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

	// check attributes processing
	if len(graph.parent.Keys) != 4 {
		// it is expected tha GraphML has 4 keys: one common and three specific to the Graph element
		t.Error("len(gr.parent.keys) != 4", len(graph.parent.Keys))
	}
	if len(graph.Data) != 4 {
		// it is expected that Graph element has 4 data elements
		t.Error("len(node.data) != 4", len(graph.Data))
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
	if len(gml.Keys) != 1 {
		t.Error("len(gml.keys) != 1", len(gml.Keys))
		return
	}

	if gml.Keys[0].Name != keyName {
		t.Error("gml.keys[0].name != keyName", gml.Keys[0].Name)
	}
	if gml.Keys[0].Description != keyDescr {
		t.Error("gml.keys[0].Description != keyDescr", gml.Keys[0].Description)
	}
	if gml.Keys[0].KeyType != "double" {
		t.Error("gml.keys[0].keyType != double", gml.Keys[0].KeyType)
	}
	if val, err := strconv.ParseFloat(gml.Keys[0].DefaultValue, 64); err != nil {
		t.Error(err)
	} else if keyDefault != val {
		t.Error("keyDefault != val", keyDefault, val)
	}
	if gml.Keys[0].Target != KeyForNode {
		t.Error("gml.keys[0].keyFor != KeyForNode", gml.Keys[0].Target)
	}
	if gml.Keys[0].ID != "d0" {
		t.Error("gml.keys[0].ID != d0", gml.Keys[0].ID)
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
	if len(gml.Keys) != 2 {
		t.Error("len(gml.keys) != 2", len(gml.Keys))
		return
	}
	if gml.Keys[1].ID != "d1" {
		t.Error("gml.keys[1].ID != d1", gml.Keys[1].ID)
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
	if len(gr.Nodes) != 1 {
		t.Error("len(gr.nodes) != 1", len(gr.Nodes))
	}
	if len(gr.parent.Keys) != 4 {
		t.Error("len(gr.parent.keys) != 4", len(gr.parent.Keys))
	}
	if len(node.Data) != 4 {
		t.Error("len(node.data) != 4", len(node.Data))
	}

	// check attributes
	checkAttributes(attributes, node.Data, KeyForNode, gr.parent, t)
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
	if len(gr.Edges) != 1 {
		t.Error("len(gr.edges) != 1", len(gr.Edges))
	}
	if len(gr.edgesMap) != 1 {
		t.Error("len(gr.edgesMap) != 1", len(gr.edgesMap))
	}
	if edge.Description != description {
		t.Error("edge.Description != description", edge.Description)
	}
	if edge.Source != n1.ID {
		t.Error("edge.source != n1.ID ", edge.Source, n1.ID)
	}
	if edge.Target != n2.ID {
		t.Error("edge.target != n2.ID", edge.Target, n2.ID)
	}
	if len(edge.Directed) != 0 {
		t.Error("len(edge.directed) != 0", len(edge.Directed))
	}
	if _, ok := gr.edgesMap[edgeIdentifier(n1.ID, n2.ID)]; !ok {
		t.Error("edge not found in edges map")
	}

	// check attributes
	checkAttributes(attributes, edge.Data, KeyForEdge, gr.parent, t)

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
				if data.Key == key.ID {
					str_val := fmt.Sprint(val)
					if data.Value != str_val {
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
