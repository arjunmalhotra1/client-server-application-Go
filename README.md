# serverclient
This is a **Client-Server application** with the following functionalities:
* multiple-threaded server;
* clients;
* External queue between the clients and server;

Clients:
* Are configured from a command line.
* Data is read from a file.
* Can request server to AddItem(), RemoveItem(), GetItem(), GetAllItems()
* Data is in the form of strings;

* Clients can be added / removed while not interfering to the server or other clients ;

Server:
* Has data structure(s) that holds the data in the memory.
  - The data structure keeps the order of items as they added. 
    For example: If client added the following keys in the following order A, B, D, E, C. 
    The GetAllItems returns A, B, D, E, C
	If item D was removed, the GetAllItems return A, B, E, C
* Server is be able to add an item, remove an item, get a single or all item from the data structure;

External queue:
*  Used is Amazon Simple Queue Service (SQS)


Clients send requests to the external queue - while the server reads those and execute them on its data structure. You define the structure of the messages (AddItem, RemoveItem, GetItem, GetAllItems)


# The flow of the project:
1. Multiple clients are sending requests to the queue (and not waiting for the response).
2. Server is reading requests from the queue and processing them, the output of the server is written to a log file
3. Server is able to process items in parallel
4. log messages (debug, error) are written to stdout


# How to run the server and the client
Open 2 terminals.
1. On one terminal run 
 > make client
2. On the other run 
   > make server

# How to run the commands on the client
You could send commands with client app:

`add {key}` - adds {key}.

`remove {key}` - removes {key}.

`get {key}` - gives true or false of the key was present.

`get-all` - gives all the keys

Examples :
```
add 1
add 2
get 1
get-all
remove 1
```
# Further improvements
1. Fix the graceful shutdown. If you notice I do have the signal channel but for some reason I started having dangling go routines. Graceful shutdown would definitely be a big improvement in the project.
2. Use context. context can be another big improvement in this project for things to shutdown gracefully, with point 1.Also context can be passed on to the methods of aws sqs.
