# MapReducer

This project contains a classical MapReduce distributed system, based on the paper:
J. Dean, S. Ghemawat - "MapReduce: Simplified Data Processing on Large Clusters" [2004]

The following implementation uses Go Programming Language and it has the scope of being an API-style app which perform the desired task. 

Use "go build -buildmode=plugin -o <so_file> <go_file>" to build a plugin and pass it as a parameter to run_worker.go