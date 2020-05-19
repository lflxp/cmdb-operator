FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/cmdbdemo \
    USER_UID=1001 \
    USER_NAME=cmdbdemo

# install operator binary
COPY build/_output/bin/cmdbdemo ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
