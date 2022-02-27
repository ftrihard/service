

FROM traefik:v2.5
## Default module name (put your setting in .env to override)
ARG PLUGIN_MODULE=github.com/ftrihard/service
ADD . /plugins-local/src/${PLUGIN_MODULE}
RUN export GOPATH=/plugins-local/src/${PLUGIN_MODULE}
