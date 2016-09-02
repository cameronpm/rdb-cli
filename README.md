# Redis RDB manipulation tool

rdb-cli tool takes redis RDB file and converts it to redis protocol (or other) format. This can be then used for ex. to pipe data in protocol format directly to redis using mass insertion feature.

# Usage example
This tool will be used in backup process to restore a running redis instance without the need to have SSH access (ex. AWS ElastiCache)

A primary use case:

 - take RDB snapshot of Redis instance,
 - format RDB file to redis protocol (RESP)
 - pipe it to redis instance using [mass insert](http://redis.io/topics/mass-insert)

After you acquire an RDB snapshot, restoring it to a running instance is as simple as running one command:

```bash
$ ./rdb-cli --format protocol ./redis-backup.rdb | redis-cli --pipe
```

# Similar tools

Similar tools in different languages:

 - python: https://github.com/sripathikrishnan/redis-rdb-tools
 - node.js: https://github.com/codeaholics/node-rdb-tools
 - rust: https://github.com/badboy/rdb-rs
