#!/bin/bash

PWD=$(pwd)

docker run -ti -v $PWD:/work --name ubuntu --privileged ubuntu /bin/bash
