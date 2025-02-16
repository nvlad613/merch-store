FROM liquibase/liquibase:latest

ENV CHANGELOG_FILE=migrations/master.xml
ENV DB_DRIVER=org.postgresql.Driver
ENV DB_URL=jdbc:postgresql://postgres:5432/appdb
ENV DB_LOGIN=appuser
ENV DB_PASSWORD=apppass

COPY migrations /liquibase/migrations

CMD ["sh", "-c", "liquibase --changeLogFile=${CHANGELOG_FILE} --username=${DB_LOGIN} --password=${DB_PASSWORD} --url=${DB_URL} update"]
