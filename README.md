## Go Axis

 Anonymous cryptocurrency based on zero-knowledge proof technology and refactored ethereum protocol by Golang.


* ### Please refer to the official wiki for the tutorial
   <https://wiki.axis.cash/en/index.html?file=Start/from-the-sourcecode-base-on-centos7>


* ### Self-apply for BetaNet mining license
   <https://axis.cash/license/records.html>

## Building the source

For prerequisites and detailed build instructions please read the
[Installation Instructions](https://github.com/axis-cash/go-axis/wiki/Building-Axis)
on the wiki.

Building axis requires both a Go (version 1.7 or later) and a C++ compiler.
You can install them using your favourite package manager.
Once the dependencies are installed, run

    make 

or, to build the full suite of utilities:

    make all

## Executables

The go-axis project comes with several wrappers/executables found in the `cmd` directory.

| Command    | Description |
|:----------:|-------------|
| **gaxis** | Our main Gaxis CLI client. It is the entry point into the Axis network (alpha or dev net), capable of running as a full node (default). It can be used by other processes as a gateway into the Axis network via JSON RPC endpoints exposed on top of HTTP, WebSocket and/or IPC transports. `gaxis --help` and the [CLI Wiki page](https://github.com/axis-cash/go-Axis/wiki/Command-Line-Options) for command line options. |
| `bootnode` | Stripped down version of our Axis client implementation that only takes part in the network node discovery protocol, but does not run any of the higher level application protocols. It can be used as a lightweight bootstrap node to aid in finding peers in private networks. |


## Running gaxis

Going through all the possible command line flags is out of scope here (please consult our
[CLI Wiki page](https://github.com/axis-cash/go-axis/wiki/Command-Line-Options)), but we've
enumerated a few common parameter combos to get you up to speed quickly on how you can run your
own axis instance.

## Axis networks

axis have 2 networks: **dev**, **main**(mainnet)

For example 

```gaxis --dev ... ``` will connect to a private network for development

```gaxis ...```  will connect to Axis's main network, it is for public testing. License is needed for mining

### Go into console with the Axis network options

By far the most common scenario is people wanting to simply interact with the Axis network:
create accounts; transfer funds; deploy and interact with contracts. 
 **Mining in axis beta network need to be licensed because it is for public testing, More detailed information can be referred to:**
 
<https://wiki.axis.cash/en/index.html?file=Start/from-the-binary-package>

To do so:

```
$ gaxis --${NETWORK_OPTIONS} console
```

**You are connecting axis beta network if startup gaxis __without__ ${NETWORK_OPTIONS}, main network is not online yet, 
 it will be online soon**

This command will:


 * Start up axis's built-in interactive [JavaScript console](https://github.com/axis-cash/console/blob/master/README.md), (via the trailing `console` subcommand) 
   through which you can invoke all official(here just reference ethereum web3 style) [`web3` methods](https://github.com/ethereum/wiki/wiki/JavaScript-API) .
   This too is optional and if you leave it out you can always attach to an already running axis instance with 
   `gaxis --datadir=${DATADIR} attach`.
   
 

### Go into console on the Axis **alpha** network

Transitioning towards developers, if you'd like to play around with creating Axis contracts, you
almost certainly would like to do that **without any real money involved** until you get the hang of the
entire system. In other words, instead of attaching to the main network, you want to join the **alpha**
network with your node, which is fully equivalent to the main network, but with play-Axis only.

```
$ gaxis --dev console
```

The `console` subcommand have the exact same meaning as above. Please see above for their explanations if you've 
skipped to here.

Specifying the `--dev` flag however will reconfigure your axis instance a bit:

 * using the default data directory (`~/.axis` on Linux for example). Note, on OSX
   and Linux this also means that attaching to a running alpha network node requires the use of a custom
   endpoint since `gaxis attach` will try to attach to a production node endpoint by default. E.g.
   `gaxis attach <datadir>/alpha/gaxis.ipc`.
   
*Note: Although there are some internal protective measures to prevent transactions from crossing
over between the main(beta) network and alpha network, you should make sure to always use separate accounts
for play-money and real-money. Unless you manually move accounts, axis will by default correctly
separate the two networks and will not make any accounts available between them.*

### Go into console on the Axis dev network

```
$ gaxis --dev console
```
With dev option, developer should config bootnode in local private network and develop new functions without affect 

outside Axis networks

#### Operating a dev network

Maintaining your own private dev network is more involved as a lot of configurations taken for granted in
the official networks need to be manually set up.

### Configuration

As an alternative to passing the numerous flags to the `gaxis` binary, you can also pass a configuration file via:

```
$ gaxis --config /path/to/your_config.toml
```

To get an idea how the file should look like you can use the `dumpconfig` subcommand to export your existing configuration:

```
$ gaxis --your-favourite-flags dumpconfig
```




### Programatically interfacing axis nodes

As a developer, sooner rather than later you'll want to start interacting with axis and the Axis
network via your own programs and not manually through the console. To aid this, axis has built-in
support for a JSON-RPC based APIs ([standard APIs](https://github.com/ethereum/wiki/wiki/JSON-RPC) and
[axis specific APIs](https://github.com/axis-cash/console/blob/master/README.md)). These can be
exposed via HTTP, WebSockets and IPC (unix sockets on unix based platforms, and named pipes on Windows).

The IPC interface is enabled by default and exposes all the APIs supported by axis, whereas the HTTP
and WS interfaces need to manually be enabled and only expose a subset of APIs due to security reasons.
These can be turned on/off and configured as you'd expect.

HTTP based JSON-RPC API options:

  * `--rpc` Enable the HTTP-RPC server
  * `--rpcaddr` HTTP-RPC server listening interface (default: "localhost")
  * `--rpcport` HTTP-RPC server listening port (default: 8545)
  * `--rpcapi` API's offered over the HTTP-RPC interface (default: "axis,net,web3")
  * `--rpccorsdomain` Comma separated list of domains from which to accept cross origin requests (browser enforced)
  * `--ws` Enable the WS-RPC server
  * `--wsaddr` WS-RPC server listening interface (default: "localhost")
  * `--wsport` WS-RPC server listening port (default: 8546)
  * `--wsapi` API's offered over the WS-RPC interface (default: "axis,net,web3")
  * `--wsorigins` Origins from which to accept websockets requests
  * `--ipcdisable` Disable the IPC-RPC server
  * `--ipcapi` API's offered over the IPC-RPC interface (default: "admin,debug,axis,miner,net,personal,shh,txpool,web3")
  * `--ipcpath` Filename for IPC socket/pipe within the datadir (explicit paths escape it)

You'll need to use your own programming environments' capabilities (libraries, tools, etc) to connect
via HTTP, WS or IPC to a axis node configured with the above flags and you'll need to speak [JSON-RPC](http://www.jsonrpc.org/specification)
on all transports. You can reuse the same connection for multiple requests!

**Note: Please understand the security implications of opening up an HTTP/WS based transport before
doing so! Hackers on the internet are actively trying to subvert Axis nodes with exposed APIs!
Further, all browser tabs can access locally running webservers, so malicious webpages could try to
subvert locally available APIs!**



#### Creating the communction center point with bootnode

With all nodes that you want to run initialized to the desired genesis state, you'll need to start a
bootstrap node that others can use to find each other in your network and/or over the internet. The
clean way is to configure and run a dedicated bootnode:

```
$ bootnode --genkey=boot.key
$ bootnode --nodekey=boot.key
```

With the bootnode online, it will display an [`xnode` URL](https://github.com/axis-cash/go-axis/wiki/xnode-url-format)
that other nodes can use to connect to it and exchange peer information. Make sure to replace the
displayed IP address information (most probably `[::]`) with your externally accessible IP to get the
actual `xnode` URL.

*Note: You could also use a full fledged axis node as a bootnode, but it's the less recommended way.*

*Note: there is bootnodes already available in axis alpha network and axis main network, setup developer's own 
bootnode is supposed to be used for dev network.*

#### Starting up your member nodes

With the bootnode operational and externally reachable (you can try `telnet <ip> <port>` to ensure
it's indeed reachable), start every subsequent axis node pointed to the bootnode for peer discovery
via the `--bootnodes` flag. It will probably also be desirable to keep the data directory of your
private network separated, so do also specify a custom `--datadir` flag.

```
$ gaxis --datadir=path/to/custom/data/folder --bootnodes=<bootnode-xnode-url-from-above>
```

*Note: Since your network will be completely cut off from the axis main and axis alpha networks, you'll also
need to configure a miner to process transactions and create new blocks for you.*

*Note: Mining on the public Axis network need apply license before hand(license@axis.vip). It is will earn axis coins 
in miner's account.*

#### Running a dev network miner


In a dev network setting however, a single CPU miner instance is more than enough for practical
purposes as it can produce a stable stream of blocks at the correct intervals without needing heavy
resources (consider running on a single thread, no need for multiple ones either). To start a axis
instance for mining, run it with all your usual flags, extended by:

```
$ gaxis <usual-flags> --mine --minerthreads=1 --axisbase=2S4kr7ZHFmgue2kLLngtWnAuHMQgV6jyv34SedvHifm1h3oomx59MEqfEmtnw3mCLnSA2FDojgjTA1WWydxHkUUt
```
Which will start mining blocks and transactions on a single CPU thread, crediting all proceedings to
the account specified by `--axisbase`. You can further tune the mining by changing the default gas
limit blocks converge to (`--targetgaslimit`) and the price transactions are accepted at (`--gasprice`).

If beginner want to do mining directlly , There is a script help beginner to create account and start mining. 
[`setup account and start mine` ](https://github.com/axis-cash/go-axis/wiki/start-mine)

## Contribution

Thank you for considering to help out with the source code! We welcome contributions from
anyone on the internet, and are grateful for even the smallest of fixes!

If you'd like to contribute to go-axis, please fork, fix, commit and send a pull request
for the maintainers to review and merge into the main code base. 

Please make sure your contributions adhere to our coding guidelines:

 * Code must adhere to the official Go [formatting](https://golang.org/doc/effective_go.html#formatting) guidelines (i.e. uses [gofmt](https://golang.org/cmd/gofmt/)).
 * Code must be documented adhering to the official Go [commentary](https://golang.org/doc/effective_go.html#commentary) guidelines.
 * Pull requests need to be based on and opened against the `master` branch.
 * Commit messages should be prefixed with the package(s) they modify.
   * E.g. "axis, rpc: make trace configs optional"

Please see the [Developers' Guide](https://github.com/axis-cash/go-axis/wiki/Developers'-Guide)
for more details on configuring your environment, managing project dependencies and testing procedures.

## Community resources

**Wechat:**  AXIS9413

**Discord:**  <https://discord.gg/n5HVxE>

**Twitter:**  <https://twitter.com/AXISdotCASH>

**Telegram:**  <https://t.me/AxisOfficial>

**Gitter:**  <https://gitter.im/axis-cash/Lobby?utm_source=share-link&utm_medium=link&utm_campaign=share-link>


## Other resources

**Official Website:** <https://axis.cash>

**White Paper:** <http://axis-media.s3-website-ap-southeast-1.amazonaws.com/Axis_ENG_V1.06.pdf>

**WIKI:** <https://wiki.axis.cash/zh/index.html?file=home-Home>

**Block Explorer:** <https://explorer.web.axis.cash/blocks.html>

**Introduction Video:** <https://v.qq.com/x/page/s0792e921ok.html>


## License

The go-axis library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html), also
included in our repository in the `COPYING.LESSER` file.

The go-axis binaries (i.e. all code inside of the `cmd` directory) is licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also included
in our repository in the `COPYING` file.

*Note: Go Axis inherit with licenses of ethereum.*
