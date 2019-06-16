```
 ______                           _______                                        
|   __ \.--.--.-----.-----.--.--.|    ___|.--.--.-----.----.-----.-----.-----.    (\(\
|   __ <|  |  |     |     |  |  ||    ___||_   _|  _  |   _|  -__|__ --|__ --|   ( =':')
|______/|_____|__|__|__|__|___  ||_______||__.__|   __|__| |_____|_____|_____|   (..(")(")
                          |_____|               |__|                             
```

**BunnyExpress** (c) 2018-19 by [SwordLord - the coding crew](http://www.swordlord.com/)


## Introduction ##

**BunnyExpress** is a CLI based mailbox configuration tool for Postfix and Dovecot. Use **BunnyExpress** to manage your domains, mailboxes, and aliases and have PostFix and Dovecot access that information for authentication as well as for routing and local delivery.

### Why the name BunnyExpress? ###

PostfixAdmin as a name for the tool was taken already. So we decided on a pun about Postfix, Pony Express and a fluffy bunny instead. Yeah, you had to be there.

## Intention ##

We started with **BunnyExpress** because of two reasons:

- We always try to reduce the attack surface on our servers. Running a DB server as well as a web server with PHP enabled and world accessible admin scripts is not necessarily the strategy we prefer for our own systems.

- Running a web server, a DBMS and PHP on a system just to manage a few records is not what we would call lightweight. And besides, it is a pain to set up and administrate.

We try to address both points with **BunnyExpress**. A small CLI based tool to edit and manage domains, mailboxes and aliases, which are stored in a single SQLite3 DB.

If you are looking for a lightweight and hassle free mail account management tool for postfix and dovecot, **BunnyExpress** might be for you.


## Status ##

**BunnyExpress** is somewhat stable. Please have a go and send in your bug reports.


## Installation ##

If you prefer, you can compile your own copy. Download this repository and have a go.

Or you could head over to our Travis based build toolchain and grab a release there.

## Configuration ##

All parameters which can be configured right now are in the file *be.config.js*. If you do not have a config file yet, just run **BunnyExpress** once and the tool will dump a copy for you.  

## Dependencies ##

Please make sure to have SQLite3 binaries installed. There are no further dependencies.

## License ##

**BunnyExpress** is published under the GNU Affero General Public Licence version 3. See the LICENCE file for details.
