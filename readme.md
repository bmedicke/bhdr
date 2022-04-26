# BHDR

BHDR is terminal user interface for Home Assistant.

* VI based keybindings
* customizable by editing JSON
* uses the Home Assistant WebSocket API for the fastest possible response time
* includes a WebSocket log-view for easy troubleshooting

*It's like ncmpcpp for your home!*

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

**example configuration**
```json
{
  "scheme": "ws",
  "server": "127.0.0.1:8123",
  "token": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
}
```

* `"scheme"` the connection protcol, this should be `ws`
* `"server"` point it to your Home Assistance instance
* `"token"` your Home Assistant long-lived access token
  * to get a token go to your Home Assistant profile ([link for locally running server](http://localhost:8123/profile)) and click **create token**

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
