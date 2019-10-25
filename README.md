# keybasePush

keybasePush is a system for pushing data and messages between devices securely using Keybase.


## Building
To build keybasePush, make sure you have Go installed and set up for your operating system. It may make things easier if you have `$GOPATH/bin` in your path.


Once Go is set up you should be able to download, build, and install with the following command:  

`go get -u github.com/kf5grd/keybasePush`


Note: keybasePush also works in Termux on android, and when built for that platform there will be several event commands you can use to trigger actions on your phone.


## Running
keybasePush _should_ run on several different platforms, though I haven't tested them all.

I will refer to each machine running keybasePush as a 'node'.

Each node needs to also have Keybase running and logged in. Each node will have a name that you can reference when sending a push to that node. By default, the name of a node will be it's Keybase device name. If you want to specify a different name, you can pass the `-name` flag.

Each node will also run a small webserver with a Rest API, which is how you will push messages to other nodes. The default port is 8617, but you can specify a different port by passing the `-port` flag.


### Desktop (Windows/MacOS/Linux)
You should be able to run it by simply typing `keybasePush` in your terminal.


### Android (Termux)
You may want to start keybasePush when your phone boots, or when you first launch Termux.


To start keybasePush when your phone boots, make sure you have the [Termux:Boot](https://wiki.termux.com/wiki/Termux:Boot) plugin installed, and create the following file at `~/.termux/boot/start-keybasepush`:

```bash
#!/data/data/com.termux/files/usr/bin/sh
$HOME/go/bin/keybasePush -name phone 1> /dev/null 2>> $HOME/keybasePush.log
```


Note: This assumer your `$GOPATH` is at `$HOME/go`. If it's not, or if you've placed the keybasePush binary somewhere else, then you will need to make the appropriate adjustments.


## Sending a Push
You can write your own tools to send pushes, but I'll show some examples here using `curl`. These examples will assume the sending node is named `laptop`, and the target node is named `phone`.


### Example 1: A basic push
The most basic of pushes requires that you fill in _at least_ the `target` and `content` fields. This push will just get printed out on the target node's stdout, and you would need to pipe keybasePush into a script that catches the output and reacts in some way. In other words, this doesn't really do anything useful on its own.

```bash
$ curl -H "Content-Type: application/json" -d '{"content":"Hello, world!","target":"phone"}' http://localhost:8617/messages                                                 
{"id":"1dc8741b","type":"message","ack":true,"sender":"laptop","target":"phone","content":"Hello, world!"}
```


You can also set the `title` field. And just like the one above, this push doesn't do much on its own:

```bash
$ curl -H "Content-Type: application/json" -d '{"title":"Hello","content":"Hello, world!","target":"phone"}' http://localhost:8617/messages
{"id":"b01d47b9","type":"message","ack":true,"sender":"laptop","target":"phone","title":"Hello","content":"Hello, world!"}
```

### Example 2: Events
You can specify an event, which can trigger an action on the recieving node. At the time of this writing (and hopefully I won't forget to update the readme in the future!) the only establisheb event commands are available when running on Termux. You can pass any event you want, but if there's no registered command for the event, it will only print to stdout on the target node, just like the other basic commands.


There are currently 3 registered events: `torch`, `openurl`, and `notify`:

- `torch` will turn on the phone's flashlight when the `content` field is set to `on`, and will turn the flashlight off when the `content` field is set to `off`
- `openurl` will open a url passed in the `content` field using the phones default browser
- `notify` will create a notification with the title being the value passed in the push's `title` field, and the body being the value passed in the push's `content` field


Here's an example of sending a notification:

```bash
$ curl -H "Content-Type: application/json" -d '{"event":"notify","title":"Hello","content":"Hello, world!","target":"phone"}' http://localhost:8617/messages
{"id":"bb310c7e","type":"message","ack":true,"sender":"laptop","target":"phone","title":"Hello","content":"Hello, world!","event":"notify"}
```