FROM scratch
USER 10000:10000
EXPOSE 8080
COPY ./app /app
ENTRYPOINT ["/app"]
