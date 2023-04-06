#!/bin/bash

docker run -d -v $(pwd)/tmp/results:/app/tmp/results --rm study/scanner:1.0