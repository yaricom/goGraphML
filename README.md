[![Build Status](https://travis-ci.org/yaricom/goGraphML.svg?branch=master)](https://travis-ci.org/yaricom/goGraphML) [![GoDoc](https://godoc.org/github.com/yaricom/goGraphML/neat?status.svg)](https://godoc.org/github.com/yaricom/goGraphML/graphml)

The GraphML support for GO language

## Overview

This repository includes implementation of [GraphML][1] specification to represent directed/undirected graphs with data-functions
attached to any element of the resulting graph. Concept of data-functions provides great flexibility in managing additional
data which can be associated with Nodes, Edges, Graphs, etc.

## Installation and Dependencies

The source code has no external dependencies except standard XML processing of GO platform. To install package into local
environment run following command in terminal:

```bash

go get -t github.com/yaricom/goGraphML

```

## Usage

Current realisation provides implementation of basic subset of GraphML specification which allows to build Graphs which
consist of Nodes and Edges with additional attributes (data-functions) associated. The root element can maintain collection
of graph elements, which allows to hold multiple Graphs in one object.

### Declaring Root GraphML

The root GraphML is the container which maintains collection of Graph elements as well as a collection of custom data-functions
definitions. The new GraphML element can be created as following:

```GO

    gml := NewGraphML("neural network solvers")

```
where:

* "neural network solvers" - is the human readable description associated with root element (optional)


### Register Custom Data-Function

The custom data-function representing particular data attribute can be registered with root element using designated
method or can be registered automatically when adding Graph, Node, Edge with specific attributes. The data-function can
be associated default value, which will can used by all elements referring this function key without providing data value.

With designated method it can be registered as following:

```GO

    key, err := gml.RegisterKey(KeyForNode, "weight", "the weight of link", reflect.Float64, 1.0)

```
where:

* KeyForNode specifies that data-function must be applied only for Nodes (see KeyForElement constants)
* "weight" - is the name of data attribute
* "the weight of link" - is the human readable description of this data-function (optional)
* reflect.Float64 - is the Kind of data (type) for accepted function value (see GraphMLDataType constants)
* 1.0 - the default value for data-function

The automatic data-function registration for any element with Add* method from provided attributes will not result in creation
of new Key definition if it is already defined in the element scope or in ALL elements scope. The existing Key elements will
be evaluated in order (see GetKey()):

* first it will be looked for Key with given name and targeting specific element (graphml|graph|node|edge)
* if not found then it will be looked for Key with given name and targeting ALL elements

If above lookup failed the new Key will be registered for given name and targeting specific element.

### Declaring a Graph

The new Graph can be added with associated attributes as following:

```GO
    attributes := make(map[string]interface{})
    attributes["default_weight"] = 1.1
    attributes["acyclic"] = false
    attributes["max_depth"] = 10

    graph, err := gml.AddGraph("the graph", EdgeDirectionDirected, attributes)

```
where:

* "the graph" - is the human readable description for the graph
* EdgeDirectionDirected specified that graph edges by default is directed (see EdgeDirection constants)
* attributes - the data attributes to be associated with this Graph element

### Declaring a Node

The Node elements can be added to the Graph as following:

```GO
    attributes := make(map[string]interface{})
    attributes["X"] = 0.1
    attributes["Y"] = 1.0
    attributes["NodeNeuronType"] = network.InputNeuron
    attributes["NodeActivationType"] = network.NullActivation


    node, err := gr.AddNode(attributes, "the input node")

```
where:

* attributes - the data attributes to be associated with this Node element
* "the input node" - is the human readable description (optional)

### Declaring an Edge

The Edge elements can be added to the Graph as following:

```GO
    attributes := make(map[string]interface{})
    attributes["weight"] = -1.1
    attributes["sourceId"] = 1
    attributes["targetId"] = 3

    edge, err := gr.AddEdge(n1, n2, attributes, EdgeDirectionDefault, "the first level")

```
where:

* n1 - the source Node element reference
* n2 - the target Node element reference
* EdgeDirectionDefault - the Edge direction specification which will override Graph direction if not EdgeDirectionDefault
* "the first level edge" - is the human readable description (optional)


### The GraphML Serialization

The collected GraphML data can be serialized into well defined XML format (see [GraphML specification][1]) using following
command:

```GO

    err = gml.Encode(writer, false)

```
where:

* writer - is an io.Writer to receive serialized data
* false - is a flag to indicate whether XML should be generated with indents to improve readability (true) or without to
have more compact representation (false)

The GraphML can also be read from serialized representation using following command:

```GO

    err = gml.Decode(reader)

```
where:

* reader - is an io.Reader to read data from

## Limitations

The current version does not implement the following parts of GraphML specification:

* Nested Graphs
* Hyper-Edges
* Ports

## References:

1. The original [GraphML specification][1]


[1]:http://graphml.graphdrawing.org/specification.html