# BHDR

BHDR is terminal user interface for Home Assistant.

* VI based keybindings
* customizable by editing JSON
* uses the Home Assistant WebSocket API for the fastest possible response time
* includes a WebSocket log-view for easy troubleshooting

*It's like editing your home with Vim!*

## toc

<!-- vim-markdown-toc GFM -->

* [installation](#installation)
* [configuration](#configuration)
* [usage](#usage)
* [key bindings](#key-bindings)

<!-- vim-markdown-toc -->

## installation

```sh
go install github.com/bmedicke/bhdr@latest
```

## configuration

* edit `bhdr.json` in your home folder
* if you don't have one bhdr will create one with `bhdr --create-config`
  * see [bhdr.json](https://github.com/bmedicke/bhdr/blob/main/bhdr.json) for the template
* `"scheme"` the connection protcol, this should be `ws`
* `"server"` point it to your Home Assistance instance
* `"token"` your Home Assistant long-lived access token
  * to get a token go to your Home Assistant profile ([link for locally running server](http://localhost:8123/profile)) and click **create token**
* `"ha-entities"` key-value-pairs of names for entities and their Home Assistant ID
* `"chordmap"` representation of the Vi grammar

## usage

The following flags are available:

* `--create-config` creates a template config in your home folder
* `--show-logs` adds a logs view that outputs websocket messages

## key bindings

* all views
  * `q` quit
  * `k` move up
  * `j` move down
  * `ctrl-f` move down a page
  * `ctrl-b` move up a page
  * `g` move to top
  * `G` move to bottom
  * `]` activate *logs* view
  * `[` activate *switches* view
* *switches* view
  * `h` collapse node, move up tree
  * `H` collapse all nodes
  * `l` expand node
  * `L` expand all nodes
  * `;` toggle entity (light, input_boolean, switch, etc.)
* *logs* view
  * `d` clear the log
  * `w` write log to `bhdr_log.json`
