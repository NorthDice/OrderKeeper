FROM ubuntu:latest
LABEL authors="kyolrize"

ENTRYPOINT ["top", "-b"]