// Package graphml implements marshaling and unmarshaling of GraphML XML documents.
package graphml

import (
	"encoding/xml"
	"fmt"
	"errors"
	"reflect"
	"strconv"
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
	XMLName     xml.Name      `xml:"graphml"`
	// The name space definitions
	xmlns       string        `xml:"xmlns,attr"`

	// Provides human readable description
	Description string        `xml:"desc,omitempty"`
	// The custom keys describing data-functions used in this or other elements
	keys        []*Key         `xml:"key,omitempty"`
	// The data associated with root element
	datas       []*Data        `xml:"data,omitempty"`
	// The graph objects encapsulated
	graphs      []*Graph       `xml:"graph,omitempty"`

	// The map to define keys by data-function name
	keysByName  map[string]*Key
}

// Description: In GraphML there may be data-functions attached to graphs, nodes, ports, edges, hyperedges and endpoint
// and to the whole collection of graphs described by the content of <graphml>. These functions are declared by <key>
// elements (children of <graphml>) and defined by <data> elements. Occurrence: <graphml>.
type Key struct {
	// The ID of this key element (in form dX, where X denotes the number of occurrences of the key element before the current one)
	ID           string        `xml:"id,attr"`
	// The name of element this key is for (graphml|graph|node|edge|hyperedge|port|endpoint|all)
	keyFor       KeyForElement        `xml:"for,attr"`
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
	ID          string        `xml:"id,attr"`
	// The default edge direction (directed|undirected)
	edgeDefault string        `xml:"edgedefault,attr"`

	// Provides human readable description
	Description string        `xml:"desc,omitempty"`
	// The nodes associated with this graph
	nodes       []*Node        `xml:"node,omitempty"`
	// The edges associated with this graph and connecting nodes
	edges       []*Edge        `xml:"edge,omitempty"`

	// The parent GraphML
	parent      *GraphML
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
	ID       string           `xml:"id,attr"`
	// The source node ID
	Source   string           `xml:"source,attr"`
	// The target node ID
	Target   string           `xml:"target,attr"`
	// The direction type of this edge (true - directed, false - undirected)
	Directed string           `xml:"directed,attr"`

	// The data associated with this edge
	Data     []*Data           `xml:"data,omitempty"`
}

// Creates new GraphML instance
func NewGraphML(description string) *GraphML {
	gml := GraphML{
		Description:description,
		keys:make([]*Key, 0),
		datas:make([]*Data, 0),
		graphs:make([]*Graph, 0),
		xmlns:"http://graphml.graphdrawing.org/xmlns",
		keysByName:make(map[string]*Key),
	}
	return &gml
}

// Register data function with GraphML instance
func (gml *GraphML) RegisterKey(forElem KeyForElement, name, description string, keyType reflect.Kind, defaultValue interface{}) (key *Key, err error) {
	if _, ok := gml.keysByName[name]; ok {
		return nil, errors.New("key with given name already registered")
	}
	count := len(gml.keys)
	key = &Key{
		ID:fmt.Sprintf("d%d", count),
		keyFor:forElem,
		name:name,
		Description:description,
	}
	// add key type (boolean, int, long, float, double, string)
	key.keyType, err = typeNameForKind(keyType)
	if err != nil {
		return nil, err
	}
	// store default value
	if defaultValue != nil {
		key.defaultValue, err = stringValueIfSupported(defaultValue, key.keyType)
		if err != nil {
			return nil, err
		}
	}

	// store key
	gml.addKey(key)

	return key, nil
}

// Creates new Graph withing GraphML
func (gml *GraphML) CreateGraph(description string, edgeDefault EdgeDirection) (*Graph, error) {
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

	gr := &Graph{
		ID:fmt.Sprintf("g%d", count),
		edgeDefault:edge_direction,
		Description:description,
		nodes:make([]*Node, 0),
		edges:make([]*Edge, 0),
		parent:gml,
	}
	// store graph in parent
	gml.graphs = append(gml.graphs, gr)
	return gr, nil
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
	var key_func *Key
	ok := false
	for key, val := range attributes {
		if key_func, ok = gr.parent.keysByName[key]; !ok {
			if key_func, err = gr.parent.RegisterKey(KeyForNode, key, "", reflect.TypeOf(val).Kind(), nil); err != nil {
				return nil, err
			}
		}
		if err := node.addDataWithKey(key, val, key_func); err != nil {
			return nil, err
		}
	}
	// add node
	gr.nodes = append(gr.nodes, node)

	return node, nil
}

// adds data attribute to the node
func (n *Node) addDataWithKey(name string, value interface{}, key *Key) (err error) {
	data := &Data{
		key:key.ID,
	}
	// add value
	data.value, err = stringValueIfSupported(value, key.keyType)
	if err != nil {
		return err
	}

	n.data = append(n.data, data)
	return nil
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
	switch keyType {
	case "boolean":
		if reflect.TypeOf(value).Kind() != reflect.Bool {
			return res, errors.New("default value has wrong data type when boolean expected")
		} else {
			res = strconv.FormatBool(value.(bool))
		}
	case "int", "long":
		if defTypeName, err := typeNameForKind(reflect.TypeOf(value).Kind()); err != nil {
			return res, err
		} else if defTypeName == "int" || defTypeName == "long" {
			res = fmt.Sprintf("%d", value)
		} else {
			return res, errors.New(fmt.Sprintf("default value has wrong data type when int/long expected: %s", defTypeName))
		}
	case "float", "double":
		if defTypeName, err := typeNameForKind(reflect.TypeOf(value).Kind()); err != nil {
			return res, err
		} else if defTypeName == "float" {
			res = strconv.FormatFloat(float64(value.(float32)), 'f', -1, 32)
		} else if defTypeName == "double" {
			res = strconv.FormatFloat(value.(float64), 'f', -1, 64)
		} else {
			return res, errors.New("default value has wrong data type when float/double expected")
		}
	case "string":
		if defTypeName, err := typeNameForKind(reflect.TypeOf(value).Kind()); err != nil {
			return res, err
		} else if defTypeName == "string" {
			res = value.(string)
		} else {
			return res, errors.New("default value has wrong data type when string expected")
		}
	}
	return res, nil
}

// appends given key
func (gml *GraphML) addKey(key *Key) {
	gml.keys = append(gml.keys, key)
	gml.keysByName[key.name] = key
}
