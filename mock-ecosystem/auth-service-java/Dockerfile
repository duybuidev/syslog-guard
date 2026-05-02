FROM eclipse-temurin:17-alpine
WORKDIR /app
COPY AuthService.java .
RUN javac AuthService.java
EXPOSE 8083
# Tối ưu JVM Memory: Chỉ dùng tối đa 64MB RAM
CMD ["java", "-Xmx64m", "-Xms16m", "AuthService"]
