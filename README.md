# kubot

Parallel robot execution on Kubernetes workloads

## Table of Contents

- [Introduction](#introduction)
- [Getting Started](#getting-started)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)

## Introduction

This project provides a method for parallelizing the executions of robot scripts over Kubernetes workloads.

## Getting Started

To get started with this project, you will need to have access to a Kubernetes cluster as well as the following
prerequisites:

### Prerequisites

- kubectl

## Installation

```
TODO implement here
```

## Usage

```
kubot --workspace . \
      --image docker.io/kubot:basic \
      --namespace [your-namespace] \
      --selector [robot-selector]
      --output [output-directory]
```

## Configuration

Here is the configuration parameters you can use;

| Name          | Description          | Default |
|---------------|----------------------|---------|
| POD_CPU_LIMIT | CPU limit per pod    | 10      | 
| POD_MEM_LIMIT | Memory limit per Pod | 128 Mb  |

## Contributing

To contribute to this project, please fork the repository and submit a pull request. All contributions are welcome!
