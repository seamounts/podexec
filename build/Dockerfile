FROM alpine-glibc

ENV KUBE_CONFIG /root/config
COPY /build/config /root/config
ADD _output/amd64/pod-exec /pod-exec
ADD frontend /frontend
ENTRYPOINT [ "/pod-exec" ]
