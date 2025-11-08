# Migration Dockerfile for mongo-migrate
FROM node:20-alpine

# Install migrate-mongo globally
RUN npm install -g migrate-mongo

# Set working directory
WORKDIR /app

# Copy migration files and config
COPY migrations ./migrations
COPY migrate-mongo-config.js ./

# Create a script to run migrations
RUN echo '#!/bin/sh' > /app/migrate.sh && \
    echo 'echo "Waiting for MongoDB to be ready..."' >> /app/migrate.sh && \
    echo 'sleep 5' >> /app/migrate.sh && \
    echo 'echo "Running migrations..."' >> /app/migrate.sh && \
    echo 'migrate-mongo up' >> /app/migrate.sh && \
    echo 'echo "Migrations completed successfully"' >> /app/migrate.sh && \
    chmod +x /app/migrate.sh

# Run the migration script
CMD ["/app/migrate.sh"]
