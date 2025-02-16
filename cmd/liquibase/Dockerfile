FROM liquibase/liquibase:latest

ENV PROPERTIES_FILE="/liquibase/liquibase.properties"
ENV DB_URL="jdbc:postgresql://postgres:5432/appdb"
ENV DB_LOGIN="appuser"
ENV DB_PASSWORD="apppass"

ENV LIQUIBASE_CMD="update"

COPY migrations /liquibase/migrations
COPY liquibase.properties /liquibase/

CMD "--defaults-file=${PROPERTIES_FILE} --username=${DB_LOGIN} --password=${DB_PASSWORD} --url=${DB_URL} ${LIQUIBASE_CMD}"