# IPFS benchmark

Tool to benchmark IPFS propagation

will spin a ton of nodes worldwide, upload / pin files and benchmark how long it takes other nodes
to fetch files, the assumption here is that we can get propagation time very low, sub 1sec, if we have a network
of well-connected nodes even if a file is not in the pinset of every node.

# How to run

make sure you have aws credentials stored in your ~/.aws folder
make sure you have a local ELK running: https://github.com/deviantony/docker-elk




