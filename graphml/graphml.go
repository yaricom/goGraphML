// Package graphml implements marshaling and unmarshaling of GraphML XML documents.
package graphml

import (
	"encoding/xml"
	"fmt"
	"errors"
	"reflect"
	"io"
)

// The the elements where data-function can be attached
type KeyForElement string

const (
	KeyForGraphML KeyForElement = "graphml"
	KeyForGraph KeyForElement = "graph"
	KeyForNode KeyForElement = "node"
	KeyForEdge KeyForElement = "edge"
	KeyForALL KeyForElement = "all"
)

// The edge direction
type EdgeDirection int

const (
	EdgeDirectionDefault EdgeDirection = iota
	EdgeDirectionDirected
	EdgeDirectionUndirected
)

// The root element
type GraphML struct {
	// The name of root element
	XMLName          xml.Name      `xml:"graphml"`
	// The name space definitions
	xmlns            string        `xml:"xmlns,attr"`

	// Provides human readable description
	Description      string        `xml:"desc,omitempty"`
	// The custom keys describing data-functions used in this or other elements
	keys             []*Key         `xml:"key,omitempty"`
	// The data associated with root element
	data             []*Data        `xml:"data,omitempty"`
	// The graph objects encapsulated
	graphs           []*Graph       `xml:"graph,omitempty"`

	// The map to look for keys by their standard identifiers (see keyIdentifier(name string, target KeyForElement))
	keysByIdentifier map[string]*Key
}

// Description: In GraphML there may be data-functions attached to graphs, nodes, ports, edges, hyperedges and endpoint
// and to the whole collection of graphs described by the content of <graphml>. These functions are declared by <key>
// elements (children of <graphml>) and defined by <data> elements. Occurrence: <graphml>.
type Key struct {
	// The ID of this key element (in form dX, where X denotes the number of occurrences of the key element before the current one)
	ID           string        `xml:"id,attr"`
	// The name of element this key is for (graphml|graph|node|edge|hyperedge|port|endpoint|all)
	target       KeyForElement `xml:"for,attr"`
	// The name of data-function associated with this key
	name         string        `xml:"attr.name,attr"`
	// The type of input to the data-function associated with this key. (Allowed values: boolean, int, long, float, double, string)
	keyType      string        `xml:"attr.type,attr"`
	// Provides human readable description
	Description  string        `xml:"desc,omitempty"`
	// The default value
	defaultValue string        `xml:"default,omitempty"`
}

// In GraphML there may be data-functions attached to graphs, nodes, ports, edges, hyperedges and endpoint and to the
// whole collection of graphs described by the content of <graphml>. These functions are declared by <key> elements
// (children of <graphml>) and defined by <data> elements. Occurrence: <graphml>, <graph>, <node>, <port>, <edge>,
// <hyperedge>, and <endpoint>.
type Data struct {
	// The ID of this data element (in form dX, where X denotes the number of occurrences of the data element before the current one)
	ID    string              `xml:"id,attr,omitempty"`
	// The ID of <key> element for this data element
	key   string              `xml:"key,attr"`

	// The data value associated with this elment
	value string              `xml:",chardata"`
}

// Describes one graph in this document. Occurrence: <graphml>, <node>, <edge>, <hyperedge>.
type Graph struct {
	// The ID of this graph element (in form gX, where X denotes the number of occurrences of the graph element before the current one)
	ID             string        `xml:"id,attr"`
	// The default edge direction (directed|undirected)
	edgeDefault    string        `xml:"edgedefault,attr"`

	// Provides human readable description
	Description    string        `xml:"desc,omitempty"`
	// The nodes associated with this graph
	nodes          []*Node        `xml:"node,omitempty"`
	// The edges associated with this graph and connecting nodes
	edges          []*Edge        `xml:"edge,omitempty"`

	// The data associated with this node
	data           []*Data        `xml:"data,omitempty"`

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
	data        []*Data        `xml:"data,omitempty"`
}

// Describes an edge in the <graph> which contains this <edge>. Occurrence: <graph>.
type Edge struct {
	// The ID of this edge element (in form eX, where X is the number of edge elements before this one)
	ID          string           `xml:"id,attr"`
	// The source node ID
	source      string           `xml:"source,attr"`
	// The target node ID
	target      string           `xml:"target,attr"`
	// The direction type of this edge (true - directed, false - undirected)
	directed    string           `xml:"directed,attr"`

	// Provides human readable description
	Description string           `xml:"desc,omitempty"`
	// The data associated with this edge
	data        []*Data          `xml:"data,omitempty"`
}

// Creates new GraphML instance
func NewGraphML(description string) *GraphML {
	gml := GraphML{
		Description:description,
		keys:make([]*Key, 0),
		data:make([]*Data, 0),
		graphs:make([]*Graph, 0),
		xmlns:"http://graphml.graphdrawing.org/xmlns",
		keysByIdentifier:make(map[string]*Key),
	}
	return &gml
}

// Encodes GraphML into provided Writer
func (gml *GraphML) Encode(w io.Writer) error {
	enc := xml.NewEncoder(w)
	err := enc.Encode(gml)
	if err == nil {
		err = enc.Flush()
	}
	return err
}

// Register data function with GraphML instance
func (gml *GraphML) RegisterKey(target KeyForElement, name, description string, keyType reflect.Kind, defaultValue interface{}) (key *Key, err error) {
	if key := gml.GetKey(name, target); key != nil {
		return nil, errors.New("key with given name already registered")
	}
	count := len(gml.keys)
	key = &Key{
		ID:fmt.Sprintf("d%d", count),
		target:target,
		name:name,
		Description:description,
	}
	// add key type (boolean, int, long, float, double, string)
	if key.keyType, err = typeNameForKind(keyType); err != nil {
		return nil, err
	}

	// store default value
	if defaultValue != nil {
		if key.defaultValue, err = stringValueIfSupported(defaultValue, key.keyType); err != nil {
			return nil, err
		}
	}

	// store key
	gml.addKey(key)

	return key, nil
}

// Looks for registered keys with specified name for a given target element. Returns Key or nil.
func (gml *GraphML) GetKey(name string, target KeyForElement) *Key {
	keyNameId := keyIdentifier(name, target)
	if key, ok := gml.keysByIdentifier[keyNameId]; ok {
		return key
	}
	return nil
}

// Creates new Graph and add it to the root GraphML
func (gml *GraphML) AddGraph(description string, edgeDefault EdgeDirection, attributes map[string]interface{}) (graph *Graph, err error) {
	count := len(gml.graphs)
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
		edgeDefault:edge_direction,
		Description:description,
		nodes:make([]*Node, 0),
		edges:make([]*Edge, 0),
		parent:gml,
		edgesMap:make(map[string]*Edge),
		edgesDirection:edgeDefault,
	}
	// add attributes
	if graph.data, err = gml.createDataAttributes(attributes, KeyForGraph); err != nil {
		return nil, err
	}

	// store graph in parent
	gml.graphs = append(gml.graphs, graph)
	return graph, nil
}

// Adds node to the graph with provided additional attributes and description
func (gr *Graph) AddNode(attributes map[string]interface{}, description string) (node *Node, err error) {
	count := len(gr.nodes)
	node = &Node{
		ID:fmt.Sprintf("n%d", count),
		Description:description,
		data:make([]*Data, 0),
	}
	// add attributes
	if node.data, err = gr.parent.createDataAttributes(attributes, KeyForNode); err != nil {
		return nil, err
	}

	// add node
	gr.nodes = append(gr.nodes, node)
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

	count := len(gr.edges)
	edge = &Edge{
		ID:fmt.Sprintf("e%d", count),
		source:source.ID,
		target:target.ID,
		Description:description,
	}
	switch edgeDirection {
	case EdgeDirectionDirected:
		edge.directed = "true"
	case EdgeDirectionUndirected:
		edge.directed = "false"
	}

	// add attributes
	if edge.data, err = gr.parent.createDataAttributes(attributes, KeyForEdge); err != nil {
		return nil, err
	}

	// add edge
	gr.edges = append(gr.edges, edge)
	gr.edgesMap[edgeIdentifier(source.ID, target.ID)] = edge

	return edge, nil
}

// appends given key
func (gml *GraphML) addKey(key *Key) {
	gml.keys = append(gml.keys, key)
	key_identifier := keyIdentifier(key.name, key.target)
	gml.keysByIdentifier[key_identifier] = key
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
		key:key.ID,
	}
	// add value
	if data.value, err = stringValueIfSupported(value, key.keyType); err != nil {
		return nil, err
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
		keyType = "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint8, reflect.Uint16:
		keyType = "int"
	case reflect.Int64, reflect.Uint32:
		keyType = "long"
	case reflect.Float32:
		keyType = "float"
	case reflect.Float64:
		keyType = "double"
	case reflect.String:
		keyType = "string"
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
	case "boolean":
		if reflect.TypeOf(value).Kind() != reflect.Bool {
			return res, errors.New("default value has wrong data type when boolean expected")
		}
	case "int", "long":
		if defTypeName, err := typeNameForKind(reflect.TypeOf(value).Kind()); err != nil {
			return res, err
		} else if !(defTypeName == "int" || defTypeName == "long") {
			return res, errors.New(
				fmt.Sprintf("default value has wrong data type when int/long expected: %s", defTypeName))
		}
	case "float", "double":
		if defTypeName, err := typeNameForKind(reflect.TypeOf(value).Kind()); err != nil {
			return res, err
		} else if !(defTypeName == "float" || defTypeName == "double") {
			return res, errors.New(
				fmt.Sprintf("default value has wrong data type when float/double expected: %s", defTypeName))
		}
	case "string":
		if defTypeName, err := typeNameForKind(reflect.TypeOf(value).Kind()); err != nil {
			return res, err
		} else if defTypeName != "string" {
			return res, errors.New(
				fmt.Sprintf("default value has wrong data type when string expected: %s", defTypeName))
		}
	}
	return fmt.Sprint(value), nil
}
