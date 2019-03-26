FROM python:3.7

EXPOSE 4000
RUN curl -sL https://deb.nodesource.com/setup_11.x | bash

RUN apt-get update && apt-get install -y nodejs npm && apt-get clean;
RUN npm install gitbook-cli -g

WORKDIR /opt/book/
COPY . /opt/book/
RUN npm install

CMD ["gitbook", "serve"]
