FROM liquibase/liquibase:latest

ENV PROPERTIES_FILE="/liquibase/liquibase.properties"
ENV LIQUIBASE_CMD="update"

COPY migrations /liquibase/migrations
COPY liquibase.properties /liquibase/

CMD "--defaults-file=${PROPERTIES_FILE} ${LIQUIBASE_CMD}"