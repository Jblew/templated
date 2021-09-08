FROM jedrzejlewandowski/schemaver-check:1.2.0 AS schemaVerCheck

FROM golang:1.17-alpine
COPY --from=schemaVerCheck /bin/schemaver-check /bin/schemaver-check

WORKDIR /app
ADD . /app

RUN /app/go-build.sh

CMD ["/bin/sh", "-c", "/app/startpoint.sh"]