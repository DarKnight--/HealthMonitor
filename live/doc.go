/*

Package live scans the internet connectivity.

It implements scanning of the internet connectivity using ping, DNS lookup and HEAD
request. It implements the action to pause OWTF if the internet connectivity is
found to be dead. In case no connectivity is detected HEAD request is only used
to confirm the state.

*/
package live
