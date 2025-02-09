#  ----------------------------------------------------------------------------------
#  Copyright 2018 hummingbird, J.
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

ARG BUILDER_BASE=registry.cn-shanghai.aliyuncs.com/winc-edgex/edgex:1.18.8-alpine3.16

FROM ${BUILDER_BASE} AS builder

WORKDIR /edge

COPY . .

RUN make cmd/hummingbird-core/hummingbird-core

#Next image - Copy built Go binary into new workspace
FROM registry.cn-shanghai.aliyuncs.com/winc-edgex/alpine:3.12

ENV APP_PORT=58081
ENV WEB_PORT=3000
EXPOSE $APP_PORT
EXPOSE $WEB_PORT

WORKDIR /var/bin/


COPY --from=builder /edge/cmd/hummingbird-core/hummingbird-core /bin/
COPY --from=builder /edge/cmd/hummingbird-ui/build /var/bin/cmd/hummingbird-ui/build

CMD ["/bin/sh", "-c", "/bin/hummingbird-core --confdir=/etc/hummingbird-core/configuration.toml"]
