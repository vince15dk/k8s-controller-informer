FROM golang

WORKDIR /code

COPY . /code

RUN go build -o k8s-instance-controller *.go

ENTRYPOINT ["/code/k8s-instance-controller"]