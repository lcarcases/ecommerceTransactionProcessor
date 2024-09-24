# Use the official PostgreSQL image from the Docker hub
FROM postgres:latest

# Environment variables to set up the database
ENV POSTGRES_USER=myuser
ENV POSTGRES_PASSWORD=mypassword
ENV POSTGRES_DB=mydatabase

# Copy the initialization script to Docker's init directory
COPY init.sql /docker-entrypoint-initdb.d/
