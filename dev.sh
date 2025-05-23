docker run -d --name jaeger -p 6831:6831/udp -p 16686:16686 jaegertracing/all-in-one:1.31.0

# need -p 8600:8600/udp for dns server. Also 8301 9302 tcp and udp for LAN and WAN server communications, 8300 tcp for gRPC.
docker run -d --name consul -p 8500:8500 -p 8600:8600/udp consul:1.11.3