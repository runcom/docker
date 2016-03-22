<!--[metadata]>
+++
title = "Managing Docker object labels"
description = "Description of labels, which are used to manage metadata on Docker objects."
keywords = ["Usage, user guide, labels, metadata, docker, documentation, examples, annotating"]
[menu.main]
parent = "engine_guide"
weight=100
+++
<![end-metadata]-->

# About labels

Labels are a mechanism for applying metadata to Docker objects, including:

- Images
- Containers
- Local daemons
- Volumes
- Networks
- Swarm nodes
- Swarm services

You can use labels to organize your images, record licensing information, annotate
relationships between containers, volumes, and networks, or in any way that makes
sense for your business or application.

# Label keys and values

A label is a key-value pair, stored as a string. You can specify multiple labels
for an object, but each key-value pair must be unique within an object. If the
same key is given multiple values, the most-recently-written value overwrites
all previous values.

## Key format recommendations

A label _key_ is the left-hand side of the key-value pair. Keys are alphanumeric
strings which may contain periods (`.`) and hyphens (`-`). Most Docker users use
images created by other organizations, and the following guidelines help to
prevent inadvertent duplication of labels across objects, especially if you plan
to use labels as a mechanism for automation.

- Authors of third-party tools should prefix each label key with the
  reverse DNS notation of a domain they own, such as `com.example.some-label`.

- Do not use a domain in your label key without the domain owner's permission.

- The `com.docker.*`, `io.docker.*` and `org.dockerproject.*` namespaces are
  reserved for Docker's internal use.

- Keys should only consist of lower-cased alphanumeric characters,
  dots and dashes (for example, `[a-z0-9-.]`).

- Keys should start *and* end with an alpha numeric character.

- Keys may not contain consecutive dots or dashes.

- Keys *without* namespace (dots) are reserved for CLI use. This allows end-
  users to add metadata to their containers and images without having to type
  cumbersome namespaces on the command-line.


These are simply guidelines and Docker does not *enforce* them. However, for
the benefit of the community, you *should* use namespaces for your label keys.


## Store structured data in labels

Label values can contain any data type as long as it can be represented as a
string. For example, consider this JSON document:


    {
        "Description": "A containerized foobar",
        "Usage": "docker run --rm example/foobar [args]",
        "License": "GPL",
        "Version": "0.0.1-beta",
        "aBoolean": true,
        "aNumber" : 0.01234,
        "aNestedArray": ["a", "b", "c"]
    }

You can store this struct in a label by serializing it to a string first:

    LABEL com.example.image-specs="{\"Description\":\"A containerized foobar\",\"Usage\":\"docker run --rm example\\/foobar [args]\",\"License\":\"GPL\",\"Version\":\"0.0.1-beta\",\"aBoolean\":true,\"aNumber\":0.01234,\"aNestedArray\":[\"a\",\"b\",\"c\"]}"

While it is *possible* to store structured data in label values, Docker treats
this data as a 'regular' string. This means that Docker doesn't offer ways to
query (filter) based on nested properties. If your tool needs to filter on
nested properties, the tool itself needs to implement this functionality.


## Add labels to images

To add labels to an image, use the `LABEL` instruction in your Dockerfile:


    LABEL [<namespace>.]<key>=<value> ...

The `LABEL` instruction adds a label to your image. A `LABEL` consists of a `<key>`
and a `<value>`.
Use an empty string for labels  that don't have a `<value>`,
Use surrounding quotes or backslashes for labels that contain
white space characters in the `<value>`:

    LABEL vendor=ACME\ Incorporated
    LABEL com.example.version.is-beta=
    LABEL com.example.version.is-production=""
    LABEL com.example.version="0.0.1-beta"
    LABEL com.example.release-date="2015-02-12"

The `LABEL` instruction also supports setting multiple `<key>` / `<value>` pairs
in a single instruction:

    LABEL com.example.version="0.0.1-beta" com.example.release-date="2015-02-12"

Long lines can be split up by using a backslash (`\`) as continuation marker:

    LABEL vendor=ACME\ Incorporated \
          com.example.is-beta= \
          com.example.is-production="" \
          com.example.version="0.0.1-beta" \
          com.example.release-date="2015-02-12"

Docker recommends you add multiple labels in a single `LABEL` instruction. Using
individual instructions for each label can result in an inefficient image. This
is because each `LABEL` instruction in a Dockerfile produces a new IMAGE layer.

You can view the labels via the `docker inspect` command:

    $ docker inspect 4fa6e0f0c678

    ...
    "Labels": {
        "vendor": "ACME Incorporated",
        "com.example.is-beta": "",
        "com.example.is-production": "",
        "com.example.version": "0.0.1-beta",
        "com.example.release-date": "2015-02-12"
    }
    ...

    # Inspect labels on container
    $ docker inspect -f "{{json .Config.Labels }}" 4fa6e0f0c678

    {"Vendor":"ACME Incorporated","com.example.is-beta":"", "com.example.is-production":"", "com.example.version":"0.0.1-beta","com.example.release-date":"2015-02-12"}

    # Inspect labels on images
    $ docker inspect -f "{{json .ContainerConfig.Labels }}" myimage


## Query labels

Besides storing metadata, you can filter images and containers by label. To list all
running containers that have the `com.example.is-beta` label:

    # List all running containers that have a `com.example.is-beta` label
    $ docker ps --filter "label=com.example.is-beta"

List all running containers with the label `color` that have a value `blue`:

    $ docker ps --filter "label=color=blue"

List all images with the label `vendor` that have the value `ACME`:

    $ docker images --filter "label=vendor=ACME"


## Container labels

    docker run \
       -d \
       --label com.example.group="webservers" \
       --label com.example.environment="production" \
       busybox \
       top

Please refer to the [Query labels](#query-labels) section above for information
on how to query labels set on a container.


## Daemon labels

    dockerd \
      --dns 8.8.8.8 \
      --dns 8.8.4.4 \
      -H unix:///var/run/docker.sock \
      --label com.example.environment="production" \
      --label com.example.storage="ssd"

These labels appear as part of the `docker info` output for the daemon:

    $ docker -D info

    Containers: 12
     Running: 5
     Paused: 2
     Stopped: 5
    Images: 672
    Server Version: 1.9.0
    Storage Driver: aufs
     Root Dir: /var/lib/docker/aufs
     Backing Filesystem: extfs
     Dirs: 697
     Dirperm1 Supported: true
    Execution Driver: native-0.2
    Logging Driver: json-file
    Kernel Version: 3.19.0-22-generic
    Operating System: Ubuntu 15.04
    Number of Hooks: 2
    CPUs: 24
    Total Memory: 62.86 GiB
    Name: docker
    ID: I54V:OLXT:HVMM:TPKO:JPHQ:CQCD:JNLC:O3BZ:4ZVJ:43XJ:PFHZ:6N2S
    Debug mode (server): true
     File Descriptors: 59
     Goroutines: 159
     System Time: 2015-09-23T14:04:20.699842089+08:00
     EventsListeners: 0
     Init SHA1:
     Init Path: /usr/bin/docker
     Docker Root Dir: /var/lib/docker
     Http Proxy: http://test:test@localhost:8080
     Https Proxy: https://test:test@localhost:8080
    WARNING: No swap limit support
    Username: svendowideit
    Registry: [https://index.docker.io/v1/]
    Labels:
     com.example.environment=production
     com.example.storage=ssd
