#  ----------------------------------------------------------------------------------
#  Copyright 2018 Dell Technologies, Inc.
#  Copyright 2018 Cavium
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
#
#  ----------------------------------------------------------------------------------

ARG BUILDER_BASE=golang:latest
FROM ${BUILDER_BASE} AS builder

WORKDIR /edge

# gitlab
COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build make cmd/mqtt-broker/mqtt-broker

#Next image - Copy built Go binary into new workspace
FROM alpine:3.16

RUN --mount=type=cache,target=/var/cache/apk apk add --update --no-cache dumb-init

EXPOSE 58090

WORKDIR /
COPY --from=builder /edge/cmd/mqtt-broker/mqtt-broker /bin/
COPY --from=builder /edge/cmd/mqtt-broker/res/configuration.yml.dist /etc/emqtt-broker/res/configuration.yml

#RUN mkdir -p  /logs/mqtt-broker

CMD ["/bin/sh", "-c", "/bin/mqtt-broker start -c=/etc/emqtt-broker/res/configuration.yml"]