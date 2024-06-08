# Stack Overflow Trends in GO

JSON API built on top of a Neo4j DB containing a graph made of Stack Overflow tags using [this](https://www.kaggle.com/datasets/stackoverflow/stack-overflow-tag-network/data) dataset.

## Usage

Requires `Docker version 26.1.1, build 4cf5afa` or later

1. extract data into `data/` (see `data/SOURCE.md`)
2. run `docker compose up`

## Original assignment topic selection

*I selected the dataset from https://www.kaggle.com/datasets/stackoverflow/stack-overflow-tag-network/data*

*The data contains tags from StackOverflow and their usage, as well as the relationships between tags based on how often they appear together in Developer Stories.*

*An application with a database built from this data could, for example, search for the most common tags (relevance, centrality - trends), the most frequent tag combinations through some clustering, etc. It could be some kind of "recommendation of relevant technologies" based on a prompt that specifies some "vagueness" of the search - e.g., I want something that can be used with React, but it doesn't have to be as closely related as Redux, perhaps some more versatile UI library (which would presumably be further away in the graph).*

*I would write the application in GO with the official Neo4j library, most likely in the form of a simple REST API.*

## Related links

https://medium.com/@matthewghannoum/simple-graph-database-setup-with-neo4j-and-docker-compose-061253593b5a

https://github.com/neo4j/neo4j-go-driver

https://neo4j.com/docs/go-manual/current/