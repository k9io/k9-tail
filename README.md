
Join the Key9 Slack channel
---------------------------

[![Slack](./images/slack.png)](https://key9identity.slack.com/)


What is they Key9 Tail?
------------------------

K9-tail is a small program that ”follows” authentication files for the Key9 Identity service.   The only authentication data that is sent to Key9 are SSH logs.   These logs are used to determine when, where, and how a user accessed a system.   For example,  SSH logs tell Key9 want “public key” was used during authentication.  They might also establish a Geo Location, by IP address, of the user logging in.  

k-tail keeps track of where it left off in authentication logs by a “waldo”.  The “waldo” file records the last position of the file.   This prevents k9-tail from “resending” logs it has already sent to Key9.  In the event an authentication log file is truncated,  the “waldo” file is reset to zero. 

Building and installing the Key9 Tail
-------------------------------------

Make sure you have Golang installed! 

<pre>
$ go mod init k9-tail
$ go mod tidy
$ go build
$ sudo mkdir -p /opt/k9/bin
$ sudo cp k9-tail /opt/k9/bin
$ sudo cp k9-tail.service /etc/systemd/system
$ sudo systemctl enable k9-tail
$ sudo systemctl start k9-tail
</pre>

You'll need to have the Key9 master configuration file.   That is located at: 

https://github.com/k9io/k9-ssh/blob/main/etc/k9.yaml

Prebuild Key9 proxy binaries
----------------------------

If you are unable to access a Golang compiler, you can download pre-built/pre-compiled binaries. These binaries are available for various architectures (i386, amd64, arm64, etc) and multiple operating systems (Linux, Solaris, NetBSD, etc).

You can find those binaries at: https://github.com/k9io/k9-binaries/tree/main/k9-tail


