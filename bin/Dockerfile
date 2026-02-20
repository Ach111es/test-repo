# Lightweight Dockerfile using pre-built binary from bin folder

FROM alpine:latest

# Install ca-certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy pre-built binary from bin directory (context from root)
COPY bin/api_sync_pct_bento .

# Copy configuration file from bin directory
COPY bin/config.json .

# Set execute permission for binary and create logs directory
RUN chmod +x /app/api_sync_pct_bento && \
    mkdir -p /app/logs

# Expose port
EXPOSE 5001

# Run the application
CMD ["./api_sync_pct_bento"]
