#
# Step 1: compile the app
#
FROM golang as builder

WORKDIR /app
COPY . .

# compile app
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags "-s -w" \
    -o /app/enumerator \
    cmd/enumerator/main.go



#
# Phase 2: prepare the runtime container, ready for production
#
FROM scratch
EXPOSE 8081

# copy our bot executable
COPY --from=builder /app/enumerator /enumerator

CMD ["/enumerator"]