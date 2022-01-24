# Distributed-system-399  Peer to Peer System

## Iteration 1

Your process should go through the following steps for this first iteration.
1. Connect with the registry to request an initial list of peer processes. See the section Registry
Communication Protocol below for more detail on the communication with the registry. The Registry
can ask for information before releasing a peer list.
2. Store this initial list of locations of peer processes.
You must be able to establish the initial connection with the registry and at minimum respond to the request for
your code. If this does not work, you will not be able to submit your code and you don’t qualify for any marks
for this first iteration.

### Registry Communication Protocol
Use the TCP/IP communication protocol for communication with the Registry. Once the connection is
established, the registry will send one of the following five requests. The registry can make these requests in
any order.

1. Team Name 
    ``` text
    <team name request> ::= “get team name”<newline>
    
    <team name response> ::= <team name><newline>
    
    <newline> ::= ‘\n’
    ```
    - Explanation: once the connection is established, the registry sends the string ‘get team name’ followed by a new
line character. The registry then waits to receive the name of the team which must be terminated by a new line
character. (This means a new line character can’t be part of the team name.)

2. Code Request
    ``` text  
    <code request> ::= “get code”<newline>
    
    <code response> ::= <language><newline><code><newline><end_of_code><newline>
    
    <language> ::= [a-zA-Z0-9]+
    
    <code> ::= .*
    
    <end_of_code> ::= “...”
    ```
    - Explanation: once the connection is established, the registry sends the string ‘get code’ followed by a new line character. The registry then waits to receive:
      1. The name of the language that the code is written in followed by a new line character. The name of the language should be expressed as an alpha-numeric string without any punctuation or white space characters. (Use 1 or more a-z character, A-Z characters and 0-9 characters.) 
   
      2. The source that is used to run your peer process. If the code is over multiple files, send the content of each file, one after the other. Any Unicode characters can be used for the code. You can indicate using comments if code is in multiple files. (We won’t run the code however. It is used exclusively to grade the quality of the code.)

      3. Once all the code is sent, indicate all code was send by sending three dots on their own line. (Your code itself should never have three dots on their own line.)

3. Receive Request 
    ``` text 
    <receive request> ::= “receive peers”<newline><numOfPeers><newline><peers>
    
    <numOfPeers> ::= <num>
    
    <peers> ::= <peer> | <peer><peers>
    
    <peer> ::= <ip><colon><port><newline>
    
    <ip> ::= <num><dot><num><dot><num><dot><num>
    
    <port> ::= <num>
    
    <num> ::= [0-9]+
    
    <dot> ::= ‘.’
    
    <colon> ::= ‘:’
    ```
    - Explanation: once the connection is established, the registry sends the string ‘receive peers’ followed by a new line character. The registry does not expect any response for this request. Instead, it sends the number of peers that will be send (on its own line) followed by a list of peers in the form of IP addresses and port numbers. The IP address and port number are separated by a colon. Each peer is on their own line.
        - For example, the request may contain the following:
         ``` text
         receive peers
         2
         136.159.5.27:41
         136.99.21.5:567
         ```  
4. Report Request 
    ``` text
    <report request> ::= “get report”<newline>

    <report response> ::= <numOfPeers><newline><peers><numOfSources><newline><sources>
    
    <sources> := <source> | <source><sources>
    
    <source> ::= <source location><newline><date><newline><numOfPeers><newline><peers>
    
    <source location> ::= <peer>
    
    <date> ::= <year><dash><month><dash><day><space><hour><colon><min><colon><sec>
    
    <day> ::= <two digit num>
    
    <month> ::= <two digit num>
    
    <year> ::= <four digit num>
    
    <hour> ::= <two digit num>
    
    <min> ::= <two digit num>
    
    <sec> ::= <two digit num>
    
    <two digit num> ::= [0-9][0-9]
    
    <four digit num> ::= [0-9][0-9][0-9][0-9]
    
    <dash> ::= ‘-’
    ```
    - Explanation: once the connection is established, the registry sends the string ‘get report’ followed by a new line character. The registry then waits to receive your current list of peers followed by a report that indicates all sources of this peer list. <br><br>
    If no receive request has been received yet, your list of peers will be empty and the list of sources will be empty as well. If you have already got a receive request then the response should contain the list of peers you received from the registry and there is only one source of peers in your list of sources.

      - For example: if the ‘get report’ request is received after the peers in the ‘receive peers’ request example above, the following could be a response to a ‘get report’ request..
         ``` text
         2
         136.159.5.27:41
         136.99.21.5:567
         1
         136.159.5.27:55921
         2021-01-25 15:18:23
         2
         136.159.5.27:41
         136.99.21.5:567
         ```
      - In this example, the first line indicates the length of the peer list. The second and third line contain the information about the peers in the list. The fourth line indicates there was one source that contributed to our list of peers. The fifth line gives the IP address and port of the single source. The sixth line the date that the source provided a list of peers. The remaining lines give the information received from the source. There is no ordering requirement for lists of peers.

5. Close Request 
    ``` text
    <close request> ::= “close”<newline>
    ```
    - Explanation: once the connection is established, the registry sends the string ‘close’ followed by a new line character. The registry does not expect any response for this request.
___
