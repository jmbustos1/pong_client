# Usa una imagen base de Go que incluye herramientas esenciales para el desarrollo
FROM golang:1.22-bullseye

# Instala las dependencias del sistema
RUN apt-get update && apt-get install -y \
    libglfw3-dev \
    libgl1-mesa-dev \
    xorg-dev \
    && rm -rf /var/lib/apt/lists/*

RUN apt-get update && apt-get install -y \
    libgl1-mesa-glx \
    libegl1-mesa \
    libxrandr2 \
    libxinerama1 \
    libxcursor1 \
    && rm -rf /var/lib/apt/lists/*

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY *.go ./
COPY assets /app/assets

# Build
# RUN go build -o /docker-app # compila la aplicacion

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080

# Run
# CMD ["/docker-app"]

# ejecutar codigo fuente
CMD ["go", "run", "."] 

