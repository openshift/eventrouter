# Eventrouter

This repository contains a simple event router for the [Kubernetes][kubernetes] project. The event router serves as an active watcher of _event_ resource in the kubernetes system, which takes those events and serializes them to standard output as JSON.  This is useful for a number of different purposes, but most notably long term behavioral analysis of your 
workloads running on your kubernetes cluster. 

## Goals

This project has several objectives, which include: 

* Persist events for longer period of time to allow for system debugging
* Allows operators to forward events to other system(s) for archiving/ML/introspection/etc. 
* It should be relatively low overhead
* Minimal dependency footprint: only the `stdout` sink is compiled in, reducing CVE exposure

### NOTE:

By default, eventrouter is configured to leverage existing EFK stacks by outputting wrapped json object which are easy to index in elastic search. 

## Sinks

The event router writes every event to standard output as a JSON object. There is
no sink configuration; to forward events to other systems (Kafka, HTTP, syslog,
...), collect the output with a log collector such as
[Vector](https://vector.dev/), which supports a wide range of outputs.

Historically this project also shipped `glog`, `http`, and `kafka` sinks. They
were never supported and have been removed ([LOG-10393](https://redhat.atlassian.net/browse/LOG-10393))
to shrink the dependency and CVE surface. A stale `sink` key in an existing
configuration is ignored.

## Non-Goals: 

* This service does not provide a querable extension, that is a responsibility of the log collector
* This service does not serve as a storage layer, that is also the responsibility of the log collector

## Running Eventrouter 
Standup: 
```
$ kubectl create -f https://raw.githubusercontent.com/openshift/eventrouter/master/yaml/eventrouter.yaml
```
Teardown: 
```
$ kubectl delete -f https://raw.githubusercontent.com/openshift/eventrouter/master/yaml/eventrouter.yaml
```

### Inspecting the output 
```
$ kubectl logs -f deployment/eventrouter -n kube-system 
``` 

Watch events roll through the system and hopefully stream into your ES cluster for mining, Hooray!

[kubernetes]: https://github.com/kubernetes/kubernetes/ "Kubernetes"
