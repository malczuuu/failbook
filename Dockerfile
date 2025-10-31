FROM gradle:9.0.0-jdk21-noble AS builder

USER root
COPY . .

ENV GRADLE_OPTS="-Dorg.gradle.daemon=false -Dorg.gradle.parallel=false"

RUN gradle build -i -x test

# Verify that exactly one JAR file was built
RUN test $(ls /home/gradle/build/libs/*.jar | wc -l) -eq 1

FROM eclipse-temurin:21-alpine

WORKDIR /app

EXPOSE 7070

COPY --from=builder /home/gradle/build/libs/*.jar /app/

ENV JAVA_OPTS_DEFAULT="\
-Dfile.encoding=UTF-8 \
-Duser.timezone=UTC \
-XX:+UseContainerSupport \
-XX:MaxRAMPercentage=75.0 \
-XX:+ExitOnOutOfMemoryError \
-XX:+HeapDumpOnOutOfMemoryError \
-XX:HeapDumpPath=/tmp"

ENV EXTRA_JAVA_OPTS=""

ENV APP_ARGS=""

ENTRYPOINT ["sh"]

CMD ["-c", "exec java $JAVA_OPTS_DEFAULT $EXTRA_JAVA_OPTS -jar /app/*.jar $APP_ARGS"]
