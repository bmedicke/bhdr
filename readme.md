# BHDR

BHDR is terminal user interface for Home Assistant.

* VI based keybindings
* customizable by editing JSON
* uses the Home Assistant WebSocket API for the fastest possible response time
* includes a WebSocket log-view for easy troubleshooting

## toc

<!-- vim-markdown-toc GFM -->

* [installation](#installation)
* [configuration](#configuration)
* [key bindings](#key-bindings)

<!-- vim-markdown-toc -->

## installation

```sh
go install github.com/bmedicke/bhdr@latest
```

## configuration

* edit `bhdr.json` in your home folder
* if you don't have one bhdr will create one with `bhdr --create-config`

## key bindings

* all views
  * `q` quit
  * `k` move up
  * `j` move down
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
  * `ctrl-f` move down page
  * `ctrl-b` move up page
  * `d` clear the log
  * `w` write log to `bhdr_log.json`
