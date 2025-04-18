FROM golang:1.24 AS build
WORKDIR /app
COPY go.mod go.sum ./
COPY cmd/*.go .
RUN go build -tags lambda.norpc -o main ./...


FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build /app/main ./main
ENTRYPOINT [ "./main" ]