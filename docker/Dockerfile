FROM busybox:latest

MAINTAINER tim@magnetic.io

# Vamp Sava showcase service _SERVICE_
# https://github.com/magneticio/sava

EXPOSE 8080

ADD public /public
ADD _SERVICE_ /

RUN mv /public/index.tmpl / 2> /dev/null || true

ENTRYPOINT ["/_SERVICE_"]
