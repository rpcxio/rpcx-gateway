using System;  
using System.IO;  
using System.Net;  
using System.Text;  

namespace example1
{
    class Program
    {
        public static void Main ()  
        {  
            // write request
            WebRequest request = WebRequest.Create ("http://127.0.0.1:9981/Arith"); 
            request.Method = "POST";  
            request.Headers.Add("X-RPCX-MessageID", "12345678");
            request.Headers.Add("X-RPCX-MesssageType", "0");
            request.Headers.Add("X-RPCX-SerializeType", "1");
            request.Headers.Add("X-RPCX-ServicePath", "Arith");
            request.Headers.Add("X-RPCX-ServiceMethod", "Mul");

            string postData = "{\"A\":10, \"B\":20}";  
            byte[] byteArray = Encoding.UTF8.GetBytes (postData);  
            request.ContentType = "application/rpcx";  
            request.ContentLength = byteArray.Length;  
            Stream dataStream = request.GetRequestStream ();  
            dataStream.Write (byteArray, 0, byteArray.Length);    
            dataStream.Close ();  

            
            // Get the response
            WebResponse response = request.GetResponse ();   
            dataStream = response.GetResponseStream ();  
            StreamReader reader = new StreamReader (dataStream);  
            string responseFromServer = reader.ReadToEnd ();  

            Console.WriteLine(responseFromServer);  

            // Clean up the streams
            reader.Close ();  
            dataStream.Close ();  
            response.Close ();  
        } 
    }
}
