[![Run on Repl.it](https://repl.it/badge/github/sayakura/Go-IRC-server)](https://repl.it/github/sayakura/Go-IRC-server)

#### From Author

This project has a lot of work to do when you are trying to cover most of the edge cases(the irc server logic). 
I had tried to implement some fancy features and they are like 98% done but I had no time to fill in the last 
piece of the puzzle. They are:
- data persistence, user can choose data persistence mode and what happened is that when the server starts, 
it will read previoust data from the local file system into memory, and I wrote a signal handler to handle 
signal interrupt, it will store all the existing data into file system(in json format) before the server shuts
down. but somehow the golang lib is not able to convert my struct into json representation string. 
- some rfc rules, which are described in the following research section of this readme file, I implemented the 
easy and obvious one, but it added a lot of complexities. 


**Explanation"
└── src
    ├── db.go       ;; database api, provides a way to manipulate the db without changing it directly
    ├── handlers.go 
                    ;; handlers for different irc commands, I have two tables, one for commands when the user is not logged in(mainly the auth ones) , another one contains the normal ones that the user can use
    ├── main.go     ;; entry point, some setup (like db initialization and flag parsing)
    ├── server.go   ;; a big for loop for new connections and go routines to handle each new connection(like a little section)
    └── utils.go    ;; utilities


## Reserch that I had done:

### Client

Each client is distinguished from other clients by **a unique
nickname having a maximum length of nine (9) character**

The nickname supplied with NICK is the name that's used to address you on IRC. The nickname must be unique across the network, so you can't use a nickname that's already in use at the time.

The username supplied with USER is simply the user part in your user@host hostmask that appears to others on IRC, showing where your connection originates from (if you've never seen these, then your client may be hiding them). In the early days of IRC it was typical for people to run their IRC client on multi-user machines, and the username corresponded to their local username on that machine. Some people do still use IRC from multi-user shell servers, but for the most part the username is vestigal.

The real name supplied with USER is used to populate the real name field that appears when someone uses the WHOIS command on your nick. Unlike the other two fields, this field can be fairly long and contain most characters (including spaces). Some people do put their real name here, but many do not.
  
### Channel

A channel is a named group of one or more clients which will all
receive messages addressed to that channel. **The channel is created
implicitly when the first client joins it, and the channel ceases to
exist when the last client leaves it.** While channel exists, any
client can reference the channel using the name of the channel.

**Channels names are strings (beginning with a '&' or '#' character) of
length up to 200 characters. Apart from the the requirement that the
first character being either '&' or '#'; the only restriction on a
channel name is that it may not contain any spaces (' '), a control G
(^G or ASCII 7), or a comma (',' which is used as a list item
separator by the protocol).**


As part of the protocol, a user
may be a part of several channels at once, but a **limit of ten (10)
channels is recommended as being ample for both experienced and
novice users.**

 

## IRC Specification

*  **Character codes**
	No specific character set is specified. The protocol is based on a a
	set of codes which are composed of eight (8) bits, making up an
	octet. Each message may be composed of any number of these octets;
	however, some octet values are used for control codes which act as
	message delimiters.

    Regardless of being an 8-bit protocol, the delimiters and keywords
    are such that protocol is mostly usable from USASCII terminal and a
    telnet connection.

*  **Message**
        Servers and clients send eachother messages which may or may not
        generate a reply.  If the message contains a valid command, as
        described in later sections, the client should expect a reply as
        specified but it is not advised to wait forever for the reply; client
        to server and server to server communication is essentially
        asynchronous in nature.

    Each IRC message may consist of up to three main parts: the prefix
    (optional), the command, and the command parameters (of which there
    may be up to 15).  The prefix, command, and all parameters are
    separated by one (or more) ASCII space character(s) (0x20).

    The presence of a prefix is indicated with a single leading ASCII
    colon character (':', 0x3b), which must be the first character of the
    message itself.  There must be no gap (whitespace) between the colon
    and the prefix.  The prefix is used by servers to indicate the true
    origin of the message.  If the prefix is missing from the message, it
    is assumed to have originated from the connection from which it was
    received.  Clients should not use prefix when sending a message from
    themselves; if they use a prefix, the only valid prefix is the
    registered nickname associated with the client.  If the source
    identified by the prefix cannot be found from the server's internal
    database, or if the source is registered from a different link than
    from which the message arrived, the server must ignore the message
    silently.

    The command must either be a valid IRC command or a three (3) digit
    number represented in ASCII text.

    **IRC messages are always lines of characters terminated with a CR-LF
    (Carriage Return - Line Feed) pair, and these messages shall not
    exceed 512 characters in length, counting all characters including
    the trailing CR-LF.Thus, there are 510 characters maximum allowed
    for the command and its parameters.** There is no provision for
    continuation message lines. 

*   **Message format in 'pseudo' BNF**
    The protocol messages must be extracted from the contiguous stream of
    octets.  The current solution is to designate two characters, CR and
    LF, as message separators.   Empty  messages  are  silently  ignored,
    which permits  use  of  the  sequence  CR-LF  between  messages
    without extra problems.

        The extracted message is parsed into the components <prefix>,
        <command> and list of parameters matched either by <middle> or
        <trailing> components.

    The BNF representation for this is:

        <message>  ::= [':' <prefix> <SPACE> ] <command> <params> <crlf>
        <prefix>   ::= <servername> | <nick> [ '!' <user> ] [ '@' <host> ]
        <command>  ::= <letter> { <letter> } | <number> <number> <number>
        <SPACE>    ::= ' ' { ' ' }
        <params>   ::= <SPACE> [ ':' <trailing> | <middle> <params> ]

        <middle>   ::= <Any *non-empty* sequence of octets not including SPACE
                    or NUL or CR or LF, the first of which may not be ':'>
        <trailing> ::= <Any, possibly *empty*, sequence of octets not including
                        NUL or CR or LF>

        <crlf>     ::= CR LF

### Parsing

#### Messages
An IRC message is a single line, delimited by with a pair of CR ('\r', 0x13) and LF ('\n', 0x10) characters.
When reading messages from a stream, read the incoming data into a buffer. Only parse and process a message once you encounter the \r\n at the end of it. If you encounter an empty message, silently ignore it.
When sending messages, ensure that a pair of \r\n characters follows every single message your software sends out.

Messages have this format:

  [@tags] [:source] <command> <parameters>
The specific parts of an IRC message are:

* tags: Optional metadata on a message, starting with ('@', 0x40).
* source: Optional note of where the message came from, starting with (':', 0x3a). Also called the prefix.
* command: The specific command this message represents.
* parameters: If it exists, data relevant to this specific command.

#### Parameters
Parameters are a series of values separated by one or more ASCII SPACE characters (' ', 0x20). However, to allow a value itself to contain spaces, the final parameter can be prepended by a (':', 0x3a) character. If the final parameter is prefixed with a colon ':', the prefix is stripped and the rest of the message is treated as the final parameter, no matter what characters it contains.

    :irc.example.com CAP * LIST :         ->  ["*", "LIST", ""]

    CAP * LS :multi-prefix sasl           ->  ["*", "LS", "multi-prefix sasl"]

    CAP REQ :sasl message-tags foo        ->  ["REQ", "sasl message-tags foo"]

    :dan!d@localhost PRIVMSG #chan :Hey!  ->  ["#chan", "Hey!"]

    :dan!d@localhost PRIVMSG #chan Hey!   ->  ["#chan", "Hey!"]



    message     =  [ "@" tags SPACE ] [ ":" prefix SPACE ] command
                    [ params ] crlf

    tags        =  tag *[ ";" tag ]
    tag         =  key [ "=" value ]
    key         =  [ vendor "/" ] 1*( ALPHA / DIGIT / "-" )
    value       =  *valuechar
    valuechar   =  <any octet except NUL, BELL, CR, LF, semicolon (`;`) and SPACE>
    vendor      =  hostname

    prefix      =  servername / ( nickname [ [ "!" user ] "@" host ] )

    command     =  1*letter / 3digit

    params      =  *( SPACE middle ) [ SPACE ":" trailing ]
    nospcrlfcl  =  <any octet except NUL, CR, LF, colon (`:`) and SPACE>
    middle      =  nospcrlfcl *( ":" / nospcrlfcl )
    trailing    =  *( ":" / " " / nospcrlfcl )


    SPACE       =  %x20 *( %x20 )   ; space character(s)
    crlf        =  %x0D %x0A        ; "carriage return" "linefeed"