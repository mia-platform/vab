FROM --platform=${TARGETPLATFORM} alpine:3.16

ARG TARGETPLATFORM

COPY ${TARGETPLATFORM}/vab /usr/local/bin/

CMD ["/usr/local/bin/vab"]
