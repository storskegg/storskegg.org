FROM alpine:3.11
RUN apk add --no-cache nodejs npm go

WORKDIR /storskegg.org
COPY package.json .
RUN npm i
COPY . .
RUN npm run build
RUN go mod vendor
RUN go build
RUN echo $PATH

ENTRYPOINT /storskegg.org/beardsrv

EXPOSE 3001

#CMD ["./photostream"]
