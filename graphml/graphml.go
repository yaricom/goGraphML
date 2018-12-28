// Package graphml implements marshaling and unmarshaling of GraphML XML documents.
package graphml

import "encoding/xml"

// The root element
type GraphML struct {
	XMLName xml.Name      `xml:"graphml"`
	Xmlns   string        `xml:"xmlns,attr"`
	// Provides human readable description
	Desc    string        `xml:"desc,omitempty"`
	Key     *Key          `xml:"key,omitempty"`
	Data    *Data         `xml:"data,omitempty"`
	Graph   *Graph        `xml:"graph,omitempty"`
}

// Description: In GraphML there may be data-functions attached to graphs, nodes, ports, edges, hyperedges and endpoint
// and to the whole collection of graphs described by the content of <graphml>. These functions are declared by <key>
// elements (children of <graphml>) and defined by <data> elements. Occurrence: <graphml>.
type Key struct {
	// Provides human readable description
	Desc    string        `xml:"desc,omitempty"`
	// The default value
	Default interface{}   `xml:"default,omitempty"`
	// The ID of this key element to be linked against
	ID      string        `xml:"id,attr"`
	// The name of element this key is for (graphml|graph|node|edge|hyperedge|port|endpoint|all)
	For     string        `xml:"for,attr"`
	// The name of data-function associated with this key
	Name    string        `xml:"attr.name,attr"`
	// The type of input to the data-function associated with this key
	Type    string        `xml:"attr.type,attr"`
}

// In GraphML there may be data-functions attached to graphs, nodes, ports, edges, hyperedges and endpoint and to the
// whole collection of graphs described by the content of <graphml>. These functions are declared by <key> elements
// (children of <graphml>) and defined by <data> elements. Occurence: <graphml>, <graph>, <node>, <port>, <edge>,
// <hyperedge>, and <endpoint>.
type Data struct {
	// The data value associated with this elment
	Value interface{}        `xml:",chardata"`
}

// Describes one graph in this document. Occurrence: <graphml>, <node>, <edge>, <hyperedge>.
type Graph struct {

}
