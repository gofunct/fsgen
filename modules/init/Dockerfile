# Build Stage
FROM {{ .docker_build_image }}:{{ .docker_build_image_version }} AS build-stage

LABEL app="build-{{ .app_name }}"
LABEL REPO="https://github.com/{{ .github_username }}/{{ .app_name }}"

ENV PROJPATH=/go/src/github.com/{{ .github_username }}/{{ .app_name }}

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/{{ .github_username }}/{{ .app_name }}
WORKDIR /go/src/github.com/{{ .github_username }}/{{ .app_name }}

RUN make build-alpine

# Final Stage
FROM {{ .docker_image }}

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/{{ .github_username }}/{{ .app_name }}"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/{{ .app_name }}/bin

WORKDIR /opt/{{ .app_name }}/bin

COPY --from=build-stage /go/src/github.com/{{ .github_username }}/{{ .app_name }}/bin/{{ .app_name }} /opt/{{ .app_name }}/bin/
RUN chmod +x /opt/{{ .app_name }}/bin/{{ .app_name }}

# Create appuser
RUN adduser -D -g '' {{ .app_name }}
USER {{ .app_name }}

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/{{ .app_name }}/bin/{{ .app_name }}"]