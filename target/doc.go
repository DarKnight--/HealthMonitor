/*

Package target scans all the active targets of the OWTF.

It uses ssdeep fuzzy algorithm to match the web response received from the previously obtained.
It scans the target present in the worklist of the OWTF. If major change is detected
then the process corresponding to the target is sent a pause signal and the user is alerted.

*/
package target
