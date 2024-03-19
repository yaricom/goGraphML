package graphml

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"reflect"
	"strconv"
	"testing"
)

func TestNewGraphML(t *testing.T) {
	description := "test"
	gml := NewGraphML(description)
	require.NotNil(t, gml)
	assert.Equal(t, description, gml.Description)
}

func TestNewGraphMLWithAttributes(t *testing.T) {
	description := "test"

	attributes := map[string]interface{}{
		"double":  10.2,
		"bool":    false,
		"integer": 120,
		"string":  "string data",
	}

	gml, err := NewGraphMLWithAttributes(description, attributes)
	require.NoError(t, err)
	require.NotNil(t, gml)
	assert.Equal(t, description, gml.Description)

	// check attributes
	attr, err := gml.GetAttributes()
	require.NoError(t, err)
	assert.Equal(t, attributes, attr)
}

func TestGraphML_Decode_keyTargetDefault(t *testing.T) {
	graphFile, err := os.Open("../data/test_graph_default_key_target.xml")
	require.NoError(t, err, "failed to open file")
	// decode
	gml := NewGraphML("")
	err = gml.Decode(graphFile)
	require.NoError(t, err, "failed to decode")

	// test results
	attributes := map[string]interface{}{
		"test-key": "test data",
	}

	// check Graph element
	//
	require.Len(t, gml.Graphs, 1, "wrong graphs number")
	graph := gml.Graphs[0]

	// check Node elements
	//
	require.Len(t, graph.Nodes, 1, "wrong nodes number")
	for i, n := range graph.Nodes {
		id := fmt.Sprintf("n%d", i)
		assert.Equal(t, id, n.ID, "wrong node ID at: %d", i)
		desc := fmt.Sprintf("test node #%d", i+1)
		assert.Equal(t, desc, n.Description, "wrong node description at: %d", i)
		// check GetAttributes
		attrs, err := n.GetAttributes()
		require.NoError(t, err, "failed to get attributes")
		assert.Equal(t, attributes, attrs)
	}

	// check Key target was set properly
	//
	key := gml.GetKey("test-key", KeyForAll)
	require.NotNil(t, key, "key expected")
	assert.Equal(t, KeyForAll, key.Target)
}

func TestGraphML_Decode_keyTypeDefault(t *testing.T) {
	graphFile, err := os.Open("../data/test_graph_default_key_type.xml")
	require.NoError(t, err, "failed to open file")
	// decode
	gml := NewGraphML("")
	err = gml.Decode(graphFile)
	require.NoError(t, err, "failed to decode")

	// test results
	attributes := map[string]interface{}{
		"integer-key": 10,
		"test-key":    "test data",
		"color":       "yellow",
	}

	// check Graph element
	//
	require.Len(t, gml.Graphs, 1, "wrong graphs number")
	graph := gml.Graphs[0]

	// check Node elements
	//
	require.Len(t, graph.Nodes, 1, "wrong nodes number")
	for i, n := range graph.Nodes {
		id := fmt.Sprintf("n%d", i)
		assert.Equal(t, id, n.ID, "wrong node ID at: %d", i)
		desc := fmt.Sprintf("test node #%d", i+1)
		assert.Equal(t, desc, n.Description, "wrong node description at: %d", i)
		// check GetAttributes
		attrs, err := n.GetAttributes()
		require.NoError(t, err, "failed to get attributes")
		assert.Equal(t, attributes, attrs)
	}

	// check Key types was set properly
	//
	key := gml.GetKey("test-key", KeyForNode)
	require.NotNil(t, key, "key expected")
	assert.Equal(t, key.KeyType, StringType)

	key = gml.GetKey("integer-key", KeyForNode)
	require.NotNil(t, key, "key expected")
	assert.Equal(t, key.KeyType, IntType)

	key = gml.GetKey("color", KeyForNode)
	require.NotNil(t, key, "key expected")
	assert.Equal(t, key.KeyType, StringType)
}

func TestGraphML_Decode_emptyStringDefault(t *testing.T) {
	graphFile, err := os.Open("../data/test_graph_default_empty_string.xml")
	require.NoError(t, err, "failed to open file")
	// decode
	gml := NewGraphML("")
	err = gml.Decode(graphFile)
	require.NoError(t, err, "failed to decode")

	// test results
	attributes := map[string]interface{}{
		"int-key":      1,
		"string-key":   "",
		"string-key-1": "test",
	}

	// check Graph element
	//
	require.Len(t, gml.Graphs, 1, "wrong graphs number")
	graph := gml.Graphs[0]

	// check Node elements
	//
	require.Len(t, graph.Nodes, 1, "wrong nodes number")
	for i, n := range graph.Nodes {
		id := fmt.Sprintf("n%d", i)
		assert.Equal(t, id, n.ID, "wrong node ID at: %d", i)
		desc := fmt.Sprintf("test node #%d", i+1)
		assert.Equal(t, desc, n.Description, "wrong node description at: %d", i)
		// check GetAttributes
		attrs, err := n.GetAttributes()
		require.NoError(t, err, "failed to get attributes")
		assert.Equal(t, attributes, attrs)
	}

	// check Key types was set properly
	//
	key := gml.GetKey("string-key", KeyForNode)
	require.NotNil(t, key, "key expected")
	assert.Equal(t, key.KeyType, StringType)
}

func TestGraphML_Decode(t *testing.T) {
	graphFile, err := os.Open("../data/test_graph.xml")
	require.NoError(t, err, "failed to open file")
	// decode
	gml := NewGraphML("")
	err = gml.Decode(graphFile)
	require.NoError(t, err, "failed to decode")

	// test results
	attributes := map[string]interface{}{
		"bool":    false,
		"integer": 120,
		"string":  "string data",
	}

	// check Graph element
	//
	require.Len(t, gml.Graphs, 1, "wrong graphs number")
	graph := gml.Graphs[0]
	assert.Equal(t, edgeDirectionDirected, graph.EdgeDefault, "wrong edge default")
	assert.Equal(t, "g0", graph.ID, "wrong graph ID")
	assert.Equal(t, EdgeDirectionDirected, graph.edgesDirection, "wrong edges direction")
	assert.Len(t, graph.edgesMap, 1, "wrong size of edges map")

	// check attributes
	attributes["double"] = 3.14
	attrs, err := graph.GetAttributes()
	require.NoError(t, err)
	assert.Equal(t, attributes, attrs)

	attributes["double"] = 10.2
	// check Node elements
	//
	require.Len(t, graph.Nodes, 2, "wrong nodes number")
	for i, n := range graph.Nodes {
		id := fmt.Sprintf("n%d", i)
		assert.Equal(t, id, n.ID, "wrong node ID at: %d", i)
		desc := fmt.Sprintf("test node #%d", i+1)
		assert.Equal(t, desc, n.Description, "wrong node description at: %d", i)
		// check attributes
		attrs, err := n.GetAttributes()
		require.NoError(t, err, "failed to get attributes")
		assert.Equal(t, attributes, attrs)
	}

	// check Edge elements
	//
	require.Len(t, graph.Edges, 1, "wrong edges number")
	for i, e := range graph.Edges {
		id := fmt.Sprintf("e%d", i)
		assert.Equal(t, id, e.ID, "wrong edge ID at: %d", i)
		assert.Equal(t, graph.Nodes[0].ID, e.Source, "wrong edge source: %v", e)
		assert.Equal(t, graph.Nodes[1].ID, e.Target, "wrong edge target: %v", e)
		assert.Equal(t, "test edge", e.Description, "wrong edge description: %v", e)
		// check attributes
		attrs, err := e.GetAttributes()
		assert.NoError(t, err, "failed to get attributes")
		assert.Equal(t, attributes, attrs)
	}
}

func TestGraphML_Encode(t *testing.T) {
	// build GraphML
	description := "test graph"
	gml := NewGraphML("TestGraphML_Encode")

	// register common data-function for all elements
	keyForAllName := "attr_double"
	_, err := gml.RegisterKey(KeyForAll, keyForAllName, "common double data-function", reflect.Float64, 10.2)
	require.NoError(t, err, "failed to register key")

	attributes := map[string]interface{}{
		keyForAllName:  NotAValue,
		"attr_bool":    false,
		"attr_integer": 120,
		"attr_string":  "string data",
	}

	// add graph
	graph, err := gml.AddGraph(description, EdgeDirectionDirected, attributes)
	require.NoError(t, err, "failed to add graph")

	// add nodes
	description = "test node #1"
	n1, err := graph.AddNode(attributes, description)
	require.NoError(t, err, "failed to add node: %s", description)

	description = "test node #2"
	n2, err := graph.AddNode(attributes, description)
	require.NoError(t, err, "failed to add node: %s", description)

	// add edge
	description = "test edge"
	_, err = graph.AddEdge(n1, n2, attributes, EdgeDirectionDefault, description)
	require.NoError(t, err, "failed to add edge: %s", description)

	// encode
	outBuf := bytes.NewBufferString("")
	err = gml.Encode(outBuf, false)

	// check results
	const resString = "<graphml xmlns=\"http://graphml.graphdrawing.org/xmlns\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://graphml.graphdrawing.org/xmlns http://graphml.graphdrawing.org/xmlns/1.0/graphml.xsd\"><desc>TestGraphML_Encode</desc><key id=\"d0\" for=\"all\" attr.name=\"attr_double\" attr.type=\"double\"><desc>common double data-function</desc><default>10.2</default></key><key id=\"d1\" for=\"graph\" attr.name=\"attr_bool\" attr.type=\"boolean\"></key><key id=\"d2\" for=\"graph\" attr.name=\"attr_integer\" attr.type=\"int\"></key><key id=\"d3\" for=\"graph\" attr.name=\"attr_string\" attr.type=\"string\"></key><key id=\"d4\" for=\"node\" attr.name=\"attr_bool\" attr.type=\"boolean\"></key><key id=\"d5\" for=\"node\" attr.name=\"attr_integer\" attr.type=\"int\"></key><key id=\"d6\" for=\"node\" attr.name=\"attr_string\" attr.type=\"string\"></key><key id=\"d7\" for=\"edge\" attr.name=\"attr_bool\" attr.type=\"boolean\"></key><key id=\"d8\" for=\"edge\" attr.name=\"attr_integer\" attr.type=\"int\"></key><key id=\"d9\" for=\"edge\" attr.name=\"attr_string\" attr.type=\"string\"></key><graph id=\"g0\" edgedefault=\"directed\"><desc>test graph</desc><node id=\"n0\"><desc>test node #1</desc><data key=\"d4\">false</data><data key=\"d0\">10.2</data><data key=\"d5\">120</data><data key=\"d6\">string data</data></node><node id=\"n1\"><desc>test node #2</desc><data key=\"d4\">false</data><data key=\"d0\">10.2</data><data key=\"d5\">120</data><data key=\"d6\">string data</data></node><edge id=\"e0\" source=\"n0\" target=\"n1\"><desc>test edge</desc><data key=\"d7\">false</data><data key=\"d0\">10.2</data><data key=\"d8\">120</data><data key=\"d9\">string data</data></edge><data key=\"d1\">false</data><data key=\"d0\">10.2</data><data key=\"d2\">120</data><data key=\"d3\">string data</data></graph></graphml>"
	assert.Equal(t, resString, outBuf.String())
}

func TestGraphML_Encode_EmptyFor(t *testing.T) {
	// build GraphML
	gml := NewGraphML("TestGraphML_Encode_EmptyFor")

	// register common data-function for all elements
	keyForAllName := "keyForAll"
	_, err := gml.RegisterKey("", keyForAllName, "", reflect.String, nil)
	require.NoError(t, err, "failed to register key")

	// encode
	outBuf := &bytes.Buffer{}
	err = gml.Encode(outBuf, false)

	// check results
	assert.NotContains(t, outBuf.String(), "for=\"\"", "a key with an empty target should omit the for attribute")
}

func TestGraphML_AddGraph(t *testing.T) {
	description := "test graph"
	gml := NewGraphML("")

	// test normal creation
	attributes := map[string]interface{}{
		"attr_double":  100.1,
		"attr_bool":    false,
		"attr_integer": 120,
		"attr_string":  "string data",
	}

	graph, err := gml.AddGraph(description, EdgeDirectionDirected, attributes)
	require.NoError(t, err, "failed to add graph")
	require.NotNil(t, graph)
	assert.Equal(t, description, graph.Description, "wrong description")
	assert.Equal(t, edgeDirectionDirected, graph.EdgeDefault, "wrong default direction")
	assert.Len(t, gml.Graphs, 1, "wrong graphs number")

	// check attributes processing
	assert.Len(t, graph.parent.Keys, 4, "wrong keys number")
	assert.Len(t, graph.Data, 4, "wrong root data functions number")

	// check attributes
	attrs, err := graph.GetAttributes()
	require.NoError(t, err, "failed to get attributes")
	assert.Equal(t, attributes, attrs)

	// test error
	graph, err = gml.AddGraph(description, EdgeDirectionDefault, nil)
	assert.EqualError(t, err, "default edge direction must be provided")
	assert.Nil(t, graph)
}

func TestGraph_GetEdge(t *testing.T) {
	description := "test graph"
	gml := NewGraphML("")

	graph, err := gml.AddGraph(description, EdgeDirectionDirected, nil)
	require.NoError(t, err, "failed to add graph")
	require.NotNil(t, graph)

	// add nodes
	description = "test node #1"
	n1, err := graph.AddNode(nil, description)
	require.NoError(t, err, "failed to add node: %s", description)

	description = "test node #2"
	n2, err := graph.AddNode(nil, description)
	require.NoError(t, err, "failed to add node: %s", description)

	// add edge
	description = "test edge"
	_, err = graph.AddEdge(n1, n2, nil, EdgeDirectionDefault, description)
	require.NoError(t, err, "failed to add edge: %s", description)

	// check existing edge
	edge := graph.GetEdge(n1.ID, n2.ID)
	assert.NotNil(t, edge, "edge is expected")

	// check non-existing edge
	edge = graph.GetEdge("n10", "n11")
	assert.Nil(t, edge, "edge is not expected")
}

func TestGraph_GetNode(t *testing.T) {
	description := "test graph"
	gml := NewGraphML("")

	graph, err := gml.AddGraph(description, EdgeDirectionDirected, nil)
	require.NoError(t, err, "failed to add graph")
	require.NotNil(t, graph)

	// add nodes
	description = "test node #1"
	n1, err := graph.AddNode(nil, description)
	require.NoError(t, err, "failed to add node: %s", description)

	description = "test node #2"
	n2, err := graph.AddNode(nil, description)
	require.NoError(t, err, "failed to add node: %s", description)

	// add edge
	description = "test edge"
	e1, err := graph.AddEdge(n1, n2, nil, EdgeDirectionDefault, description)
	require.NoError(t, err, "failed to add edge: %s", description)

	// check existing node
	node := graph.GetNode(n1.ID)
	assert.NotNil(t, node, "node is expected")
	assert.Same(t, n1, node, "node should be the same as n1")

	// check non-existing node
	node = graph.GetNode("n42")
	assert.Nil(t, node, "node is not expected")

	// check edge nodes
	node = e1.SourceNode()
	assert.NotNil(t, node, "node is expected")
	assert.Same(t, n1, node, "node should be the same as n1")

	node = e1.TargetNode()
	assert.NotNil(t, node, "node is expected")
	assert.Same(t, n2, node, "node should be the same as n2")
}

// tests if default value of common key will be used for empty attributes
func TestGraphML_RegisterKeyDefaultValue(t *testing.T) {
	description := "graphml"
	gml := NewGraphML(description)

	commonKeyName := "weight"
	keyDesc := "the weight function"
	keyDefault := 100.2
	commonKey, err := gml.RegisterKey(KeyForAll, commonKeyName, keyDesc, reflect.TypeOf(keyDefault).Kind(), keyDefault)
	require.NoError(t, err, "failed to register key: %s", commonKeyName)

	// register graph and test number of keys
	attributes := map[string]interface{}{
		commonKeyName:  NotAValue, // empty attribute - default value will be used
		"attr_bool":    false,
		"attr_integer": 120,
		"attr_string":  "string data",
	}

	graph, err := gml.AddGraph(description, EdgeDirectionDirected, attributes)
	require.NoError(t, err, "failed to add graph")
	require.NotNil(t, graph)

	// check attributes processing
	assert.Len(t, graph.parent.Keys, 4,
		"it is expected tha GraphML has 4 keys: one common and three specific to the Graph element")
	assert.Len(t, graph.Data, 4, "it is expected that Graph element has 4 data elements")

	// check if default value of common key was used
	for _, d := range graph.Data {
		if d.ID == commonKey.ID {
			assert.Equal(t, commonKey.DefaultValue, d.Value, "wrong default key value")
		}
	}

	// check attribute without value and without default value
	commonKeyName2 := "height"
	_, err = gml.RegisterKey(KeyForAll, commonKeyName2, keyDesc, reflect.TypeOf(keyDefault).Kind(), nil)
	require.NoError(t, err, "failed to register key: %s", commonKeyName2)

	// add node with empty attribute which has no default value
	attributes[commonKeyName2] = NotAValue
	n, err := graph.AddNode(attributes, "test node with empty attribute without default value")
	assert.Nil(t, n, "node should not be added due to attribute error")
	assert.EqualError(t, err, "empty attribute without default value: height",
		"error must be raised for empty attribute and no default value")
}

func TestGraphML_RegisterKeyForAll(t *testing.T) {
	description := "graphml"
	gml := NewGraphML(description)

	commonKeyName := "weight"
	keyDesc := "the weight function"
	keyDefault := 100.0
	_, err := gml.RegisterKey(KeyForAll, commonKeyName, keyDesc, reflect.TypeOf(keyDefault).Kind(), keyDefault)
	require.NoError(t, err, "failed to register key: %s", commonKeyName)

	// register graph and test number of keys
	attributes := map[string]interface{}{
		commonKeyName:  100.1,
		"attr_bool":    false,
		"attr_integer": 120,
		"attr_string":  "string data",
	}

	graph, err := gml.AddGraph(description, EdgeDirectionDirected, attributes)
	require.NoError(t, err, "failed to add graph")
	require.NotNil(t, graph)

	// check attributes processing
	assert.Len(t, graph.parent.Keys, 4,
		"it is expected tha GraphML has 4 keys: one common and three specific to the Graph element")
	assert.Len(t, graph.Data, 4, "it is expected that Graph element has 4 data elements")
}

func TestGraphML_RegisterKey(t *testing.T) {
	description := "graphml"
	gml := NewGraphML(description)

	keyName := "weight"
	keyDesc := "the weight function"
	keyDefault := 100.0
	_, err := gml.RegisterKey(KeyForNode, keyName, keyDesc, reflect.TypeOf(keyDefault).Kind(), keyDefault)
	require.NoError(t, err, "failed to register key: %s", keyName)
	require.Len(t, gml.Keys, 1)

	// check key attributes
	key := gml.Keys[0]
	assert.Equal(t, keyName, key.Name)
	assert.Equal(t, keyDesc, key.Description)
	assert.Equal(t, DoubleType, key.KeyType)
	val, err := strconv.ParseFloat(key.DefaultValue, 64)
	require.NoError(t, err, "failed to parse key value")
	assert.Equal(t, keyDefault, val)
	assert.Equal(t, KeyForNode, key.Target)
	assert.Equal(t, "d0", key.ID)

	// register key with the same standard identifier and check that error raised
	_, err = gml.RegisterKey(KeyForNode, keyName, keyDesc, reflect.TypeOf(keyDefault).Kind(), keyDefault+100)
	assert.EqualError(t, err, fmt.Sprintf("key with given name already registered: %s", keyName),
		"error should be raised when attempting to register Key with already existing standard identifier")

	// register another key and check ID
	keyName = "height"
	_, err = gml.RegisterKey(KeyForNode, keyName, keyDesc, reflect.TypeOf(keyDefault).Kind(), keyDefault+100)
	require.NoError(t, err, "failed to register key with name: %s", keyName)
	require.Len(t, gml.Keys, 2, "wrong keys number")
	assert.Equal(t, "d1", gml.Keys[1].ID)
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
	attributes := map[string]interface{}{
		"attr_double":  100.1,
		"attr_bool":    false,
		"attr_integer": 120,
		"attr_string":  "string data",
	}

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
	attrs, err := node.GetAttributes()
	require.NoError(t, err, "failed to get attributes")
	assert.Equal(t, attributes, attrs)
}

func TestGraph_RemoveKey(t *testing.T) {
	gmlattrs := map[string]interface{}{
		"k4": 8000,
	}
	gml, err := NewGraphMLWithAttributes("", gmlattrs)
	require.NoError(t, err, "failed to create GraphML")
	k1name := "k1"
	k1target := KeyForAll
	k1, err := gml.RegisterKey(k1target, k1name, "", reflect.Int, 0)
	require.NoError(t, err, "failed to register key: %s", k1name)
	k2name := "k2"
	k2target := KeyForNode
	k2, err := gml.RegisterKey(k2target, k2name, "", reflect.Int, 0)
	require.NoError(t, err, "failed to register key: %s", k2name)
	k3name := "k3"
	k3target := KeyForEdge
	k3, err := gml.RegisterKey(k3target, k3name, "", reflect.Int, 0)
	require.NoError(t, err, "failed to register key: %s", k3name)
	require.Len(t, gml.Keys, 4)

	grattrs := map[string]interface{}{
		"k1": 999,
	}
	gr, err := gml.AddGraph("test graph", EdgeDirectionDirected, grattrs)
	require.NoError(t, err, "failed to add graph")

	// add elements
	n1attrs := map[string]interface{}{
		"k1": 100,
		"k2": 10,
	}
	n1, err := gr.AddNode(n1attrs, "test node 1")
	require.NoError(t, err, "failed to add node 1")
	n2attrs := map[string]interface{}{
		"k1": 200,
		"k2": 20,
	}
	n2, err := gr.AddNode(n2attrs, "test node 2")
	require.NoError(t, err, "failed to add node 2")
	e1attrs := map[string]interface{}{
		"k1": 300,
		"k3": 3,
	}
	e1, err := gr.AddEdge(n1, n2, e1attrs, EdgeDirectionDefault, "test edge")
	require.NoError(t, err, "failed to add edge")

	// try removing k1
	err = gml.RemoveKey(k1)
	require.NoError(t, err, "failed to remove key 1")
	key := gml.GetKey(k1name, k1target)
	assert.Nil(t, key)
	attrs, _ := gr.GetAttributes()
	assert.NotContains(t, attrs, "k1")
	attrs, _ = n1.GetAttributes()
	assert.NotContains(t, attrs, "k1")
	attrs, _ = n2.GetAttributes()
	assert.NotContains(t, attrs, "k1")
	attrs, _ = e1.GetAttributes()
	assert.NotContains(t, attrs, "k1")

	// try removing k2
	err = gml.RemoveKey(k2)
	require.NoError(t, err, "failed to remove key 2")
	key = gml.GetKey(k2name, k2target)
	assert.Nil(t, key)
	attrs, _ = n1.GetAttributes()
	assert.NotContains(t, attrs, "k2")
	attrs, _ = n2.GetAttributes()
	assert.NotContains(t, attrs, "k2")

	// try removing k3
	err = gml.RemoveKey(k3)
	require.NoError(t, err, "failed to remove key 3")
	key = gml.GetKey(k3name, k3target)
	assert.Nil(t, key)
	attrs, _ = e1.GetAttributes()
	assert.NotContains(t, attrs, "k3")

	// try removing k4
	k4name := "k4"
	k4target := KeyForGraphML
	k4 := gml.GetKey(k4name, k4target)
	assert.NotNil(t, k4)
	err = gml.RemoveKey(k4)
	require.NoError(t, err, "failed to remove key 4")
	key = gml.GetKey(k4name, k4target)
	assert.Nil(t, key)
	attrs, _ = gml.GetAttributes()
	assert.NotContains(t, attrs, "k4")

	// try removing k1 once again
	err = gml.RemoveKey(k1)
	assert.Error(t, err)
}

func TestGraph_RemoveKeyByName(t *testing.T) {
	gmlattrs := map[string]interface{}{
		"k4": 8000,
	}
	gml, err := NewGraphMLWithAttributes("", gmlattrs)
	require.NoError(t, err, "failed to create GraphML")
	k1name := "k1"
	k1target := KeyForAll
	_, err = gml.RegisterKey(KeyForAll, "k1", "", reflect.Int, 0)
	require.NoError(t, err, "failed to register key: %s", k1name)

	grattrs := map[string]interface{}{
		"k1": 999,
	}
	gr, err := gml.AddGraph("test graph", EdgeDirectionDirected, grattrs)
	require.NoError(t, err, "failed to add graph")

	// add elements
	n1attrs := map[string]interface{}{
		"k1": 100,
		"k2": 10,
	}
	n1, err := gr.AddNode(n1attrs, "test node 1")
	require.NoError(t, err, "failed to add node 1")
	n2attrs := map[string]interface{}{
		"k1": 200,
		"k2": 20,
	}
	n2, err := gr.AddNode(n2attrs, "test node 2")
	require.NoError(t, err, "failed to add node 2")
	e1attrs := map[string]interface{}{
		"k1": 300,
		"k3": 3,
	}
	e1, err := gr.AddEdge(n1, n2, e1attrs, EdgeDirectionDefault, "test edge")
	require.NoError(t, err, "failed to add edge")

	// try removing k1
	err = gml.RemoveKeyByName(k1target, k1name)
	require.NoError(t, err, "failed to remove key 1")
	key := gml.GetKey(k1name, k1target)
	assert.Nil(t, key)
	attrs, _ := gr.GetAttributes()
	assert.NotContains(t, attrs, "k1")
	attrs, _ = n1.GetAttributes()
	assert.NotContains(t, attrs, "k1")
	attrs, _ = n2.GetAttributes()
	assert.NotContains(t, attrs, "k1")
	attrs, _ = e1.GetAttributes()
	assert.NotContains(t, attrs, "k1")

	// try removing not existing key
	err = gml.RemoveKeyByName(KeyForAll, "not existing")
	assert.Error(t, err, "key no found")
}

func TestNode_GetAttributes(t *testing.T) {
	description := "test graph"
	gml := NewGraphML("")
	gr, err := gml.AddGraph(description, EdgeDirectionDirected, nil)
	require.NoError(t, err, "failed to add graph")

	// add node
	attributes := map[string]interface{}{
		"attr_double":  100.1,
		"attr_bool":    false,
		"attr_integer": 120,
		"attr_string":  "string data",
	}

	description = "test node"
	node, err := gr.AddNode(attributes, description)
	require.NoError(t, err, "failed to add node: %s", description)

	// get attributes and test
	nAttr, err := node.GetAttributes()
	require.NoError(t, err, "failed to get attributes")
	require.Len(t, nAttr, len(attributes))
	assert.Equal(t, attributes, nAttr)
}

func TestGraph_AddEdge(t *testing.T) {
	description := "test graph"
	gml := NewGraphML("")
	gr, err := gml.AddGraph(description, EdgeDirectionDirected, nil)
	require.NoError(t, err, "failed to add graph")

	// Add nodes
	n1, err := gr.AddNode(nil, "#1")
	require.NoError(t, err)
	n2, err := gr.AddNode(nil, "#2")
	require.NoError(t, err)

	// Add graph
	attributes := map[string]interface{}{
		"attr_double":  100.1,
		"attr_bool":    false,
		"attr_integer": 120,
		"attr_string":  "string data",
	}

	description = "test edge"
	edge, err := gr.AddEdge(n1, n2, attributes, EdgeDirectionDefault, description)
	require.NoError(t, err, "failed to add edge: %s", description)

	// test results
	assert.Len(t, gr.Edges, 1)
	assert.Len(t, gr.edgesMap, 1)
	assert.Equal(t, description, edge.Description)
	assert.Equal(t, n1.ID, edge.Source)
	assert.Equal(t, n2.ID, edge.Target)
	assert.Empty(t, edge.Directed, "directed should be empty")

	assert.Contains(t, gr.edgesMap, edgeIdentifier(n1.ID, n2.ID), "edge not found in edges map")

	// check attributes
	attrs, err := edge.GetAttributes()
	require.NoError(t, err, "failed to get attributes")
	assert.Equal(t, attributes, attrs)

	// check error by adding the same node
	edge, err = gr.AddEdge(n1, n2, attributes, EdgeDirectionDefault, description)
	require.EqualError(t, err, "edge already added to the graph")

	// check no error for directed node backward
	edge, err = gr.AddEdge(n2, n1, attributes, EdgeDirectionDefault, description)
	assert.NoError(t, err)
}

func TestEdge_GetAttributes(t *testing.T) {
	description := "test graph"
	gml := NewGraphML("")
	gr, err := gml.AddGraph(description, EdgeDirectionDirected, nil)
	require.NoError(t, err, "failed to add graph")

	// Add nodes
	n1, err := gr.AddNode(nil, "#1")
	require.NoError(t, err)
	n2, err := gr.AddNode(nil, "#2")
	require.NoError(t, err)

	// Add graph
	attributes := map[string]interface{}{
		"attr_double":  100.1,
		"attr_bool":    false,
		"attr_integer": 120,
		"attr_string":  "string data",
	}

	description = "test edge"
	edge, err := gr.AddEdge(n1, n2, attributes, EdgeDirectionDefault, description)
	require.NoError(t, err, "failed to add edge: %s", description)

	// get attributes and check results
	attrs, err := edge.GetAttributes()
	require.NoError(t, err, "failed to get attributes")
	require.NotNil(t, attrs)
	assert.Len(t, attrs, len(attributes))
	for k, v := range attrs {
		assert.Equal(t, attributes[k], v, "wrong attribute for: %s", k)
	}
}

func TestGraphML_stringValueIfSupported(t *testing.T) {
	res, err := stringValueIfSupported(true, BooleanType)
	require.NoError(t, err)
	bRes, err := strconv.ParseBool(res)
	require.NoError(t, err)
	assert.True(t, bRes)

	testInt := 42
	res, err = stringValueIfSupported(testInt, "int")
	require.NoError(t, err)
	iRes, err := strconv.ParseInt(res, 10, 32)
	require.NoError(t, err)
	assert.EqualValues(t, testInt, iRes)

	testLong := int64(12993888475775)
	res, err = stringValueIfSupported(testLong, "long")
	require.NoError(t, err)
	lRes, err := strconv.ParseInt(res, 10, 64)
	require.NoError(t, err)
	assert.Equal(t, testLong, lRes)

	testFloat := float32(0.5)
	res, err = stringValueIfSupported(testFloat, "float")
	require.NoError(t, err)
	fRes, err := strconv.ParseFloat(res, 32)
	require.NoError(t, err)
	assert.EqualValues(t, testFloat, fRes)

	testDouble := 10000.552
	res, err = stringValueIfSupported(testDouble, "double")
	require.NoError(t, err)
	dRes, err := strconv.ParseFloat(res, 64)
	require.NoError(t, err)
	assert.Equal(t, testDouble, dRes)

	testString := "test string"
	res, err = stringValueIfSupported(testString, "string")
	require.NoError(t, err)
	assert.Equal(t, testString, res)
}
