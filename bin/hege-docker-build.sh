#!/bin/bash
# Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

set -euxo pipefail

for T in dependencies runtime demo ; do
  docker build --target="$T" --tag="jfsmig/hegemonie-$T" .
  docker push "jfsmig/hegemonie-$T:latest"
done
