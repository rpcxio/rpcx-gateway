import java.io.BufferedReader;
import java.io.DataOutputStream;
import java.io.InputStreamReader;
import java.io.Reader;
import java.net.HttpURLConnection;
import java.net.URL;

public class Main {

    public static void main(String[] args) throws Exception {
        String arg = "{\"A\":10, \"B\":20}";
        byte[] requestedPayload = arg.getBytes("utf-8");

        URL url = new URL("http://vps.colobu.com:9981/");
        HttpURLConnection conn = (HttpURLConnection) url.openConnection();
        conn.setDoOutput(true);
        conn.setInstanceFollowRedirects(false);
        conn.setRequestMethod("POST");
        conn.setRequestProperty("Content-Type", "application/rpcx");
        conn.setRequestProperty("charset", "utf-8");
        conn.setRequestProperty("Content-Length", Integer.toString(requestedPayload.length));

        conn.setRequestProperty("X-RPCX-MessageID", "12345678");
        conn.setRequestProperty("X-RPCX-MesssageType", "0");
        conn.setRequestProperty("X-RPCX-SerializeType", "1");
        conn.setRequestProperty("X-RPCX-ServicePath", "Arith");
        conn.setRequestProperty("X-RPCX-ServiceMethod", "Mul");



        conn.setUseCaches(false);
        try (DataOutputStream wr = new DataOutputStream(conn.getOutputStream())) {
            wr.write(requestedPayload);
        }



        // read reply
        Reader in = new BufferedReader(new InputStreamReader(conn.getInputStream(), "UTF-8"));
        for (int c; (c = in.read()) >= 0; )
            System.out.print((char) c);
    }
}
