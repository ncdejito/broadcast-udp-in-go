## Problem
I'm trying to receive UDP packets sent to 255.255.255.255 or more specific broadcast addresses (127.255.255.255 for example). I can't for the life of me, even though golang seems to set the SO_BROADCAST sockopt per default (confirmed with debugger). I haven't tried a raw C version yet, but I'm getting a bit stumped.

        address, _ := cmd.PersistentFlags().GetString("address")
        server, _ := cmd.PersistentFlags().GetBool("server")
        port, _ := cmd.PersistentFlags().GetUint("port")

        listenAddr := ":0"
        if server {
            listenAddr = fmt.Sprintf("%s:%d", address, port)
        }

        pc, err := net.ListenPacket("udp4", listenAddr)
        if err != nil {
            panic(err)
        }
        defer pc.Close()

        if server {
            fmt.Printf("Listening on %s\n", pc.LocalAddr().String())
            buf := make([]byte, 1024)
            n, addr, err := pc.ReadFrom(buf)
            if err != nil {
                panic(err)
            }

            fmt.Printf("%s sent this: %s\n", addr, buf[:n])
        } else {
            dstAddr := fmt.Sprintf("%s:%d", address, port)
            addr, err := net.ResolveUDPAddr("udp4", dstAddr)
            if err != nil {
                panic(err)
            }

            fmt.Printf("Sending to %s\n", dstAddr)
            _, err = pc.WriteTo([]byte("data to transmit"), addr)
            if err != nil {
                panic(err)
            }
        }

// I'm gonna try on linux just to make sure, and then uh, try in C?

// This is the code out of go net

func setDefaultSockopts(s, family, sotype int, ipv6only bool) error {
    if family == syscall.AF_INET6 && sotype != syscall.SOCK_RAW {
        // Allow both IP versions even if the OS default
        // is otherwise. Note that some operating systems
        // never admit this option.
        syscall.SetsockoptInt(s, syscall.IPPROTO_IPV6, syscall.IPV6_V6ONLY, boolint(ipv6only))
    }
    if (sotype == syscall.SOCK_DGRAM || sotype == syscall.SOCK_RAW) && family != syscall.AF_UNIX {
        // Allow broadcast.
        return os.NewSyscallError("setsockopt", syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1))
    }
    return nil
}

## Solution

So, I continued down this road this morning in Rust and golang on linux.

TL;DR: you must bind to 0.0.0.0 to receive broadcast messages.

Setting the broadcast flag matters on the sender, as not setting it won't allow to send_to 255.255.255.255 (or any other broadcast for that matter), resulting in a EACCESS.

Setting the broadcast on the receiver doesn't have an effect. What you need to do is bind to "0.0.0.0". You can then further restrict which broadcasts you react to by using a SO_BINDTODEVICE setsockopt call:

            if udpConn, succ := pc.(*net.UDPConn); succ {
                c, err := udpConn.SyscallConn()
                if err != nil {
                    panic(err)
                }
                err = c.Control(func(fd uintptr) {
                    fmt.Printf("Binding socket %d to interface %s\n", fd, ifname)
                    err = syscall.SetsockoptString(int(fd), syscall.SOL_SOCKET, syscall.SO_BINDTODEVICE, ifname)
                    if err != nil {
                        panic(err)
                    }
                })
                if err != nil {
                    panic(err)
                }
            }
You can then either send UDP to 255.255.255.255 (and catch all interfaces, except lo? <- needs a bit more investigation but my curiosity at this point is mostly satiated), or to individual interfaces (which will catch all listeners either bound to that interface, or unbound listeners).

## things i tried that didnt work
https://github.com/aler9/howto-udp-broadcast-golang


## next steps
go through 
[Intro](http://www-net.cs.umass.edu/wireshark-labs/Wireshark_Intro_v8.0.pdf)
[UDP](http://www-net.cs.umass.edu/wireshark-labs/Wireshark_UDP_v8.0.pdf)
[UDP Module](https://github.com/holwech/UDP-module)