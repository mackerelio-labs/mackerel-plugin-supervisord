# mackerel-plugin-supervisord

mackerel metric plugin for [Supervisor](https://supervisord.org/)

## Synopsis

```shell
mackerel-plugin-supervisord -uri unix:/var/run/supervisor.sock
mackerel-plugin-supervisord -uri http://user:password@127.0.0.1:9001/RPC2
```

### Help

```shell
# mackerel-plugin-supervisord --help
Usage of mackerel-plugin-supervisord:
  -uri uri
        socket uri (default "unix:/var/run/supervisor.sock")
```

## Outputs

```shell
# mackerel-plugin-supervisord
supervisord.init-1.state    20  1766634695
supervisord.init-2.state    20  1766634695
```

- The metric name is `supervisord.#.state`, `#` replace to process_name.
- The value is Process States. about value read more [Subprocesses — Supervisor documentation](https://supervisord.org/subprocess.html#process-states)

## Example of mackerel-agent.conf

```
[plugin.metrics.supervisord]
command = ["/path/to/mackerel-plugin-supervisord"]
```