import java.io.IOException;
import java.net.BindException;
import java.net.ServerSocket;
import java.net.Socket;
import java.util.Random;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.logging.*;

/**
 * Registry process that allows peers in the system to register as peers and ask for
 * the location of other peers in the system.
 * <p>
 * Communication with individual peers is managed using RequestProcessor.
 * @author Nathaly Verwaal
 *
 */
public class SubmissionRegistry extends Registry {
	
	/**
	 * Create registry to run at specified port number	
	 * @param aPortNumber port number to attempt running this registry at.
	 */
	public SubmissionRegistry(int aPortNumber) {
		super(aPortNumber);
	}
	
	/**
	 * Starts this registry and accepts connection requests from peers.  For each
	 * connection request, a RequestProcessor object is created and provided to the
	 * thread pool.
	 * @throws IOException if there are problems starting this registry server or if there
	 * are problems communication alonger a connection with a peer.
	 */
	public void start() throws IOException {
		ExecutorService executor = Executors.newFixedThreadPool(THREAD_POOL_SIZE);
		try {
			ServerSocket server = new ServerSocket(portNumber);
			LOGGER.log(Level.INFO, "Server started at " + 
					server.getInetAddress().getLocalHost().getHostAddress() +
					":" + portNumber);
			while (true) {
				Socket sock = server.accept();
				LOGGER.log(Level.INFO, "Connection accepted with " + sock.getRemoteSocketAddress());
				executor.execute(new SubmissionRequestProcessor(sock, this));				
			}
			//server.close();
		} catch (BindException be) {
			LOGGER.log(Level.SEVERE, "Unable to start registry at port " + portNumber);
		}
		executor.shutdown();
	}
	

	/**
	 * Starts the registry server. If a port number is provided as a runtime argument, 
	 * it will be used to start the registry.
	 * Otherwise, the port number provided as an argument will be used. 
	 * <p>
	 * If we can't start the registry server, the stack trace for the exception will
	 * be printed and the program ended.
	 * 
	 * @param args optional port number as a first argument.
	 */
	public static void main(String[] args)  {
		int portNumber = DEFAULT_PORT_NUMBER;
		if (args.length > 0) {
			try {
				portNumber = Integer.parseInt(args[0]);
			} catch (NumberFormatException nfe) {
				System.out.println("Expected first argument to be a port number.  Argument ignored.");
			}
		}
		SubmissionRegistry r = new SubmissionRegistry(portNumber);
		Thread t = new Thread() { 
			public void run() {
				try {
					r.start();
				} catch (IOException e) {
					// Show that an error occurred with exception info
					e.printStackTrace();
					// end program with a error code
					System.exit(1);
				}
			} 
		};
		t.start();
	}
}
