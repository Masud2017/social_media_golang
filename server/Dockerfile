FROM golang

WORKDIR /app
# COPY go.mod go.sum ./

COPY . .


#RUN go mod download
RUN make download-dep
# COPY *.go ./
#RUN CGO_ENABLED=0 GOOS=linux go build -o /social_media
RUN make build

EXPOSE 4443
CMD ["/social_media"]
