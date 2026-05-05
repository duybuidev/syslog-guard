import com.sun.net.httpserver.HttpServer;
import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;

public class AuthService {
    public static void main(String[] args) throws IOException {
        HttpServer server = HttpServer.create(new InetSocketAddress(8083), 0);

        server.createContext("/health", exchange -> {
            String resp = "Auth Service is healthy";
            exchange.sendResponseHeaders(200, resp.length());
            OutputStream os = exchange.getResponseBody();
            os.write(resp.getBytes());
            os.close();
        });

        server.createContext("/login", exchange -> {
            System.out.println("[INFO] User authenticated successfully.");
            String resp = "Token: jwt_mock_token_123";
            exchange.sendResponseHeaders(200, resp.length());
            OutputStream os = exchange.getResponseBody();
            os.write(resp.getBytes());
            os.close();
        });

        server.createContext("/crash", exchange -> {
            System.out.println("[FATAL] DB Connection lost. Terminating JVM.");
            System.exit(1);
        });

        System.out.println("Auth Service started on port 8083...");
        server.setExecutor(null);
        server.start();
    }
}
