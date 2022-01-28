import java.io.BufferedReader;
import java.io.BufferedWriter;
import java.io.File;
import java.io.FileWriter;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.OutputStreamWriter;
import java.net.Socket;
import java.util.Date;
import java.util.logging.FileHandler;
import java.util.logging.Level;
import java.util.logging.Logger;
import java.util.logging.SimpleFormatter;

class SubmissionRequestProcessor extends RequestProcessor implements Runnable {
	SubmissionRequestProcessor(Socket aPeer, Registry aSource) {
		super(aPeer, aSource);
	}
	
	@Override
	public void run() {
		try {
			log(Level.INFO, "start");
			Peer addedPeer = addPeer();
			getReport(addedPeer);
			sendPeers();
			addPeer();
			getReport(addedPeer);
			sendPeers();
			getReport(addedPeer);
			closePeer();
		} catch (IOException e) {
			log(Level.WARNING, "Problem processing socket.", e);
		}
	}
	
}