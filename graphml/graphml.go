// Package graphml implements marshaling and unmarshaling of GraphML XML documents.
package graphml

import (
	"encoding/xml"
	"fmt"
	"errors"
	"reflect"
	"io"
	"strconv"
)

// The Not value of data attribute to substitute with default one if present
var NotAValue interface{} = nil

// The elements where data-function can be attached
type KeyForElement string

const (
	// the data-function is for root GraphML element only
	KeyForGraphML KeyForElement = "graphml"
	// the data-function is for Graph element only
	KeyForGraph KeyForElement = "graph"
	// the data-function is for Node element only
	KeyForNode KeyForElement = "node"
	// the data-function is for Edge element only
	KeyForEdge KeyForElement = "edge"
	// the data-function is for all elements
	KeyForAll KeyForElement = "all"
)

// The GraphML data types
type GraphMLDataType string

const (
	// boolean (reflect.Bool)
	BooleanType GraphMLDataType = "boolean"
	// single integer precision (reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint8, reflect.Uint16)
	IntType GraphMLDataType = "int"
	// double integer precision (reflect.Int64, reflect.Uint32)
	LongType GraphMLDataType = "long"
	// single float precision (reflect.Float32)
	FloatType GraphMLDataType = "float"
	// double float precision (reflect.Float64)
	DoubleType GraphMLDataType = "double"
	// string value (reflect.String)
	StringType GraphMLDataType = "string"
)

// The edge direction
type EdgeDirection int

const (
	// edge direction not specified
	EdgeDirectionDefault EdgeDirection = iota
	// edge is directed
	EdgeDirectionDirected
	// edge is undirected
	EdgeDirectionUndirected
)

// The root element
type GraphML struct {
	// The name of root element
	XMLName            xml.Name      `xml:"graphml"`
	// The name space definitions
	XmlNS              string        `xml:"xmlns,attr"`
	// The XML schema definition
	XmlNS_XSI          string        `xml:"xmlns:xsi,attr"`
	XSI_schemaLocation string        `xml:"xsi:schemaLocation,attr"`

	// Provides human readable description
	Description        string        `xml:"desc,omitempty"`
	// The custom keys describing data-functions used in this or other elements
	Keys               []*Key         `xml:"key,omitempty"`
	// The data associated with root element
	Data               []*Data        `xml:"data,omitempty"`
	// The graph objects encapsulated
	Graphs             []*Graph       `xml:"graph,omitempty"`

	// The map to look for keys by their standard identifiers (see keyIdentifier(name string, target KeyForElement))
	keysByIdentifier   map[string]*Key
	// The map to look for keys by their IDs. Useful for fast reverse mapping of Data -> Key -> Attribute Name/Type
	keysById           map[string]*Key
}

// Description: In GraphML there may be data-functions attached to graphs, nodes, ports, edges, hyperedges and endpoint
// and to the whole collection of graphs described by the content of <graphml>. These functions are declared by <key>
// elements (children of <graphml>) and defined by <data> elements. Occurrence: <graphml>.
type Key struct {
	// The ID of this key element (in form dX, where X denotes the number of occurrences of the key element before the current one)
	ID           string        `xml:"id,attr"`
	// The name of element this key is for (graphml|graph|node|edge|hyperedge|port|endpoint|all)
	Target       KeyForElement `xml:"for,attr"`
	// The name of data-function associated with this key
	Name         string        `xml:"attr.name,attr"`
	// The type of input to the data-function associated with this key. (Allowed values: boolean, int, long, float, double, string)
	KeyType      string        `xml:"attr.type,attr"`
	// Provides human readable description
	Description  string        `xml:"desc,omitempty"`
	// The default value
	DefaultValue string        `xml:"default,omitempty"`
}

// In GraphML there may be data-functions attached to graphs, nodes, ports, edges, hyperedges and endpoint and to the
// whole collection of graphs described by the content of <graphml>. These functions are declared by <key> elements
// (children of <graphml>) and defined by <data> elements. Occurrence: <graphml>, <graph>, <node>, <port>, <edge>,
// <hyperedge>, and <endpoint>.
type Data struct {
	// The ID of this data element (in form dX, where X denotes the number of occurrences of the data element before the current one)
	ID    string              `xml:"id,attr,omitempty"`
	// The ID of <key> element for this data element
	Key   string              `xml:"key,attr"`

	// The data value associated with this elment
	Value string              `xml:",chardata"`
}

// Describes one graph in this document. Occurrence: <graphml>, <node>, <edge>, <hyperedge>.
type Graph struct {
	// The ID of this graph element (in form gX, where X denotes the number of occurrences of the graph element before the current one)
	ID             string        `xml:"id,attr"`
	// The default edge direction (directed|undirected)
	EdgeDefault    string        `xml:"edgedefault,attr"`

	// Provides human readable description
	Description    string        `xml:"desc,omitempty"`
	// The nodes associated with this graph
	Nodes          []*Node        `xml:"node,omitempty"`
	// The edges associated with this graph and connecting nodes
	Edges          []*Edge        `xml:"edge,omitempty"`
	// The data associated with this node
	Data           []*Data        `xml:"data,omitempty"`

	// The parent GraphML
	parent         *GraphML
	// The map of edges by connected nodes
	edgesMap       map[string]*Edge
	// The default edge direction flag
	edgesDirection EdgeDirection
}

// Describes one node in the <graph> containing this <node>. Occurrence: <graph>.
type Node struct {
	// The ID of this node element (in form nX, where X denotes the number of occurrences of the node element before the current one)
	ID          string        `xml:"id,attr"`
	// Provides human readable description
	Description string        `xml:"desc,omitempty"`
	// The data associated with this node
	Data        []*Data        `xml:"data,omitempty"`

	// The reference to the parent graph for reverse mapping
	graph       *Graph
}

// Describes an edge in the <graph> which contains this <edge>. Occurrence: <graph>.
type Edge struct {
	// The ID of this edge element (in form eX, where X is the number of edge elements before this one)
	ID          string           `xml:"id,attr"`
	// The source node ID
	Source      string           `xml:"source,attr"`
	// The target node ID
	Target      string           `xml:"target,attr"`
	// The direction type of this edge (true - directed, false - undirected)
	Directed    string           `xml:"directed,attr,omitempty"`

	// Provides human readable description
	Description string           `xml:"desc,omitempty"`
	// The data associated with this edge
	Data        []*Data          `xml:"data,omitempty"`

	// The reference to the parent graph for reverse mapping
	graph       *Graph
}

// Creates new GraphML instance
func NewGraphML(description string) *GraphML {
	gml := GraphML{
		Description:description,
		Keys:make([]*Key, 0),
		Data:make([]*Data, 0),
		Graphs:make([]*Graph, 0),
		XmlNS:"http://graphml.graphdrawing.org/xmlns",
		XmlNS_XSI:"http://www.w3.org/2001/XMLSchema-instance",
		XSI_schemaLocation:"http://graphml.graphdrawing.org/xmlns http://graphml.graphdrawing.org/xmlns/1.0/graphml.xsd",
		keysByIdentifier:make(map[string]*Key),
		keysById:make(map[string]*Key),
	}
	return &gml
}

// Creates new GraphML instance with given attributes
func NewGraphMLWithAttributes(description string, attributes map[string]interface{}) (gml *GraphML, err error) {
	gml = NewGraphML(description)
	// add attributes
	if gml.Data, err = gml.createDataAttributes(attributes, KeyForGraphML); err != nil {
		return nil, err
	}
	return gml, nil
}

// Encodes GraphML into provided Writer. If withIndent set then each element begins on a new indented line.
func (gml *GraphML) Encode(w io.Writer, withIndent bool) error {
	enc := xml.NewEncoder(w)
	if withIndent {
		enc.Indent("  ", "    ")
	}
	err := enc.Encode(gml)
	if err == nil {
		err = enc.Flush()
	}
	return err
}

// Decodes GraphML from provided Reader
func (gml *GraphML) Decode(r io.Reader) error {
	dec := xml.NewDecoder(r)
	err := dec.Decode(gml)
	if err != nil {
		return err
	}

	// populate auxiliary data structure
	for _, key := range gml.Keys {
		gml.keysByIdentifier[keyIdentifier(key.Name, key.Target)] = key
		gml.keysById[key.ID] = key
	}

	for _, gr := range gml.Graphs {
		gr.parent = gml
		if gr.EdgeDefault == "directed" {
			gr.edgesDirection = EdgeDirectionDirected
		} else if gr.EdgeDefault == "undirected" {
			gr.edgesDirection = EdgeDirectionUndirected
		}
		// populate edges map
		gr.edgesMap = make(map[string]*Edge)
		for _, e := range gr.Edges {
			gr.edgesMap[edgeIdentifier(e.Source, e.Target)] = e
		}
	}

	return err
}

// Register data function with GraphML instance
func (gml *GraphML) RegisterKey(target KeyForElement, name, description string, keyType reflect.Kind, defaultValue interface{}) (key *Key, err error) {
	if key := gml.GetKey(name, target); key != nil {
		return nil, errors.New("key with given name already registered")
	}
	count := len(gml.Keys)
	key = &Key{
		ID:fmt.Sprintf("d%d", count),
		Target:target,
		Name:name,
		Description:description,
	}
	// add key type (boolean, int, long, float, double, string)
	if key.KeyType, err = typeNameForKind(keyType); err != nil {
		return nil, err
	}

	// store default value
	if defaultValue != nil {
		if key.DefaultValue, err = stringValueIfSupported(defaultValue, key.KeyType); err != nil {
			return nil, err
		}
	}

	// store key
	gml.addKey(key)

	return key, nil
}

// Looks for registered keys with specified name for a given target element. If specific target has no registered key then
// common target (KeyForAll) will be checked next. Returns Key (either specific or common) or nil.
func (gml *GraphML) GetKey(name string, target KeyForElement) *Key {
	if key, ok := gml.keysByIdentifier[keyIdentifier(name, target)]; ok {
		// found element specific data-function
		return key
	} else if key, ok = gml.keysByIdentifier[keyIdentifier(name, KeyForAll)]; ok {
		// found common data-function with given name
		return key
	}
	return nil
}

// Creates new Graph and add it to the root GraphML
func (gml *GraphML) AddGraph(description string, edgeDefault EdgeDirection, attributes map[string]interface{}) (graph *Graph, err error) {
	count := len(gml.Graphs)
	var edge_direction string
	switch edgeDefault {
	case EdgeDirectionDirected:
		edge_direction = "directed"
	case EdgeDirectionUndirected:
		edge_direction = "undirected"
	default:
		return nil, errors.New("default edge direction must provided")
	}

	graph = &Graph{
		ID:fmt.Sprintf("g%d", count),
		EdgeDefault:edge_direction,
		Description:description,
		Nodes:make([]*Node, 0),
		Edges:make([]*Edge, 0),
		parent:gml,
		edgesMap:make(map[string]*Edge),
		edgesDirection:edgeDefault,
	}
	// add attributes
	if graph.Data, err = gml.createDataAttributes(attributes, KeyForGraph); err != nil {
		return nil, err
	}

	// store graph in parent
	gml.Graphs = append(gml.Graphs, graph)
	return graph, nil
}

// Adds node to the graph with provided additional attributes and description
func (gr *Graph) AddNode(attributes map[string]interface{}, description string) (node *Node, err error) {
	count := len(gr.Nodes)
	node = &Node{
		ID:fmt.Sprintf("n%d", count),
		Description:description,
		Data:make([]*Data, 0),
	}
	// add attributes
	if node.Data, err = gr.parent.createDataAttributes(attributes, KeyForNode); err != nil {
		return nil, err
	}

	// add node
	node.graph = gr
	gr.Nodes = append(gr.Nodes, node)
	return node, nil
}

// Adds edge to the graph which connects two its nodes with provided additional attributes and description
func (gr *Graph) AddEdge(source, target *Node, attributes map[string]interface{}, edgeDirection EdgeDirection, description string) (edge *Edge, err error) {
	// test if edge already exists
	edge_identification := edgeIdentifier(source.ID, target.ID)
	exists := false
	if _, exists = gr.edgesMap[edge_identification]; !exists && (edgeDirection == EdgeDirectionUndirected || gr.edgesDirection == EdgeDirectionUndirected) {
		// check other direction for undirected edge or graph types
		edge_identification = edgeIdentifier(target.ID, source.ID)
		_, exists = gr.edgesMap[edge_identification]
	}
	if exists {
		return nil, errors.New("edge already added to the graph")
	}

	count := len(gr.Edges)
	edge = &Edge{
		ID:fmt.Sprintf("e%d", count),
		Source:source.ID,
		Target:target.ID,
		Description:description,
	}
	switch edgeDirection {
	case EdgeDirectionDirected:
		edge.Directed = "true"
	case EdgeDirectionUndirected:
		edge.Directed = "false"
	}

	// add attributes
	if edge.Data, err = gr.parent.createDataAttributes(attributes, KeyForEdge); err != nil {
		return nil, err
	}

	// add edge
	edge.graph = gr
	gr.Edges = append(gr.Edges, edge)
	gr.edgesMap[edgeIdentifier(source.ID, target.ID)] = edge

	return edge, nil
}

// method to test if edge exists between given nodes. If edge exists it will be returned, otherwise nil returned
func (gr *Graph) GetEdge(sourceId, targetId string) *Edge {
	edge_identification := edgeIdentifier(sourceId, targetId)
	if edge, ok := gr.edgesMap[edge_identification]; ok {
		return edge
	}
	return nil
}

// returns data attributes map associated with Node
func (n *Node) GetAttributes() (map[string]interface{}, error) {
	return attributesForData(n.Data, n.graph.parent)
}

// returns data attributes map associated with Edge
func (e *Edge) GetAttributes() (map[string]interface{}, error) {
	return attributesForData(e.Data, e.graph.parent)
}

// builds attributes map for specified data array
func attributesForData(data []*Data, gml *GraphML) (map[string]interface{}, error) {
	attr := make(map[string]interface{})
	for _, d := range data {
		if key, ok := gml.keysById[d.Key]; !ok {
			return nil, errors.New(fmt.Sprintf("failed to find attribute name/type by id: %s", d.Key))
		} else if value, err := valueByType(d.Value, key.KeyType); err != nil {
			return nil, err
		} else {
			attr[key.Name] = value
		}
	}
	return attr, nil
}

// appends given key
func (gml *GraphML) addKey(key *Key) {
	gml.Keys = append(gml.Keys, key)
	key_identifier := keyIdentifier(key.Name, key.Target)
	gml.keysByIdentifier[key_identifier] = key
	gml.keysById[key.ID] = key
}

// Creates data-functions from given attributes and appends definitions of created functions to the provided data list.
func (gml *GraphML) createDataAttributes(attributes map[string]interface{}, target KeyForElement) (data []*Data, err error) {
	data = make([]*Data, len(attributes))
	count := 0
	for key, val := range attributes {
		key_func := gml.GetKey(key, target)
		if key_func == nil {
			// register new Key
			if key_func, err = gml.RegisterKey(target, key, "", reflect.TypeOf(val).Kind(), nil); err != nil {
				// failed
				return nil, err
			}
		}
		if d, err := createDataWithKey(val, key_func); err != nil {
			// failed
			return nil, err
		} else {
			data[count] = d
		}
		count++
	}
	return data, nil
}

// Creates data object with specified name, value and for provided Key
func createDataWithKey(value interface{}, key *Key) (data *Data, err error) {
	data = &Data{
		Key:key.ID,
	}
	// add value
	if value != NotAValue {
		if data.Value, err = stringValueIfSupported(value, key.KeyType); err == nil {
			return data, nil
		}
	} else if key.Target == KeyForAll && len(key.DefaultValue) > 0 {
		// use default value
		data.Value = key.DefaultValue
	} else {
		// raise error
		return nil, errors.New(fmt.Sprintf("empty attribute without default value: %s", key.Name))
	}
	return data, nil
}

// returns standard edge identifier based on provided iDs of connected nodes
func edgeIdentifier(soure, target string) string {
	return fmt.Sprintf("%s<->%s", soure, target)
}

// returns standard key identifier based on provided name and target
func keyIdentifier(name string, target KeyForElement) string {
	return fmt.Sprintf("%s_for_%s", name, target)
}

// Returns type name for a given kind
func typeNameForKind(kind reflect.Kind) (string, error) {
	var keyType string
	switch kind {
	case reflect.Bool:
		keyType = BooleanType
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint8, reflect.Uint16:
		keyType = IntType
	case reflect.Int64, reflect.Uint32:
		keyType = LongType
	case reflect.Float32:
		keyType = FloatType
	case reflect.Float64:
		keyType = DoubleType
	case reflect.String:
		keyType = StringType
	default:
		return "unsupported", errors.New("usupported data type for key")
	}
	return keyType, nil
}

// Converts provided value to string if it's supported by this keyType
func stringValueIfSupported(value interface{}, keyType string) (string, error) {
	res := "unsupported"
	// check that key and value types compatible
	switch keyType {
	case BooleanType:
		if reflect.TypeOf(value).Kind() != reflect.Bool {
			return res, errors.New("default value has wrong data type when boolean expected")
		}
	case IntType, LongType:
		if defTypeName, err := typeNameForKind(reflect.TypeOf(value).Kind()); err != nil {
			return res, err
		} else if !(defTypeName == IntType || defTypeName == LongType) {
			return res, errors.New(
				fmt.Sprintf("default value has wrong data type when int/long expected: %s", defTypeName))
		}
	case FloatType, DoubleType:
		if defTypeName, err := typeNameForKind(reflect.TypeOf(value).Kind()); err != nil {
			return res, err
		} else if !(defTypeName == FloatType || defTypeName == DoubleType) {
			return res, errors.New(
				fmt.Sprintf("default value has wrong data type when float/double expected: %s", defTypeName))
		}
	case StringType:
		if defTypeName, err := typeNameForKind(reflect.TypeOf(value).Kind()); err != nil {
			return res, err
		} else if defTypeName != StringType {
			return res, errors.New(
				fmt.Sprintf("default value has wrong data type when string expected: %s", defTypeName))
		}
	}
	return fmt.Sprint(value), nil
}

// Converts provided string value to the specified data type
func valueByType(val string, keyType string) (interface{}, error) {
	switch keyType {
	case BooleanType:
		return strconv.ParseBool(val)
	case IntType, LongType:
		if i_val, err := strconv.ParseInt(val, 10, 64); err != nil {
			return nil, err
		} else if keyType == IntType {
			return int(i_val), nil
		} else {
			return i_val, nil
		}
	case FloatType, DoubleType:
		if f_val, err := strconv.ParseFloat(val, 64); err != nil {
			return nil, err
		} else if keyType == FloatType {
			return float32(f_val), nil
		} else {
			return f_val, nil
		}
	case StringType:
		return val, nil
	default:
		return nil, errors.New(fmt.Sprintf("unsupported value type: %s", keyType))
	}
}
