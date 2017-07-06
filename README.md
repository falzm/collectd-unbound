collectd-unbound: Unbound statistics with collectd
==================================================

This utility allows you to access [Unbound][0] statistics from [collectd][1] using the [exec][2] plugin. It hasn't been thoroughly tested, use it at your own risk...

Building
--------

```
go get -u github.com/falzm/collectd-unbound
```

Configuration
-------------

Once the `collectd-unbound` binary is compiled, copy it wherever you want and add the following lines to you collectd configuration:

```
LoadPlugin exec
<Plugin "exec">
    Exec "unbound" "/usr/bin/collectd-unbound"
</Plugin>
```

Note: the utility executes the command `unbound-control stats` to fetch the statistics: make sure the user specified in your `exec` plugin block has the permissions to execute the command. To verify this (replace *unbound* with your user of choice):

```
$ sudo -u unbound unbound-control stats
thread0.num.queries=6080
thread0.num.cachehits=6080
thread0.num.cachemiss=0
thread0.num.prefetch=0
thread0.num.recursivereplies=0
thread0.requestlist.avg=0
thread0.requestlist.max=0
thread0.requestlist.overwritten=0
thread0.requestlist.exceeded=0
thread0.requestlist.current.all=0
thread0.requestlist.current.user=0
thread0.recursion.time.avg=0.000000
thread0.recursion.time.median=0
total.num.queries=6080
total.num.cachehits=6080
total.num.cachemiss=0
total.num.prefetch=0
total.num.recursivereplies=0
total.requestlist.avg=0
total.requestlist.max=0
total.requestlist.overwritten=0
total.requestlist.exceeded=0
total.requestlist.current.all=0
total.requestlist.current.user=0
total.recursion.time.avg=0.000000
total.recursion.time.median=0
time.now=1426870655.404178
time.up=201805.946672
time.elapsed=4.427567
```

See the _STATISTIC COUNTERS_ section from the [unbound-control][4] command manpage for information regarding metrics signification.

TODO
----

 * As of now the plugin only processes the metrics prefixed with *total*. Feel free to hack this according to your needs.

License
-------

Copyright (c) 2015, Marc Falzon.
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

 * Redistributions of source code must retain the above copyright
   notice, this list of conditions and the following disclaimer.

 * Redistributions in binary form must reproduce the above copyright
   notice, this list of conditions and the following disclaimer in the
   documentation and/or other materials provided with the distribution.

 * Neither the name of the authors nor the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.


[0]: https://unbound.net/
[1]: https://collectd.org/
[2]: https://collectd.org/documentation/manpages/collectd-exec.5.shtml
[3]: https://github.com/octo/go-collectd/
[4]: https://www.unbound.net/documentation/unbound-control.html
