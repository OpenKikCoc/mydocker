#!/bin/bash

./bin/docker run --ti busybox /bin/sh

# ./bin/docker run --ti -m 100m -- stress --vm-bytes 200m --vm-keep -m 1