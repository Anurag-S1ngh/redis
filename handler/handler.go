package handler

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/Anurag-S1ngh/redis/resp"
	"github.com/Anurag-S1ngh/redis/store"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		data, err := resp.Parse(reader)
		if err != nil {
			fmt.Println("error while reading", err)
			return
		}

		if len(data) == 0 {
			conn.Write([]byte("-ERR empty command\r\n"))
			continue
		}

		command := strings.ToLower(data[0])
		switch command {
		case "ping":
			conn.Write([]byte("+PONG\r\n"))

		case "echo":
			if len(data) != 2 {
				conn.Write([]byte("-ERR wrong number of arguments for echo\r\n"))
				continue
			}
			msg := data[1]
			fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(msg), msg)

		case "set":
			if len(data) < 3 {
				fmt.Fprintf(conn, "-ERR wrong number of arguments for SET command '%d'\r\n", len(data))
				continue
			}
			key, value := data[1], data[2]
			var ttlSeconds int
			hasExpiry := false
			if len(data) == 5 {
				if strings.ToLower(data[3]) != "ex" {
					fmt.Fprintf(conn, "-ERR expected EX got=%s\r\n", data[3])
					continue
				}
				ttlSeconds, err = strconv.Atoi(data[4])
				if err != nil {
					fmt.Fprint(conn, "-ERR wrong format for SET command\r\n")
					continue
				}
				hasExpiry = true
			}

			store.SETValue(key, value, ttlSeconds, hasExpiry)
			fmt.Fprint(conn, "+OK\r\n")

		case "get":
			if len(data) != 2 {
				fmt.Fprintf(conn, "-ERR wrong number of arguments for GET command '%d'\r\n", len(data))
				continue
			}
			value, err := store.GETValue(data[1])
			if err != nil {
				fmt.Fprintf(conn, "$-1\r\n")
				continue
			}
			fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(value), value)

		case "del":
			if len(data) < 2 {
				fmt.Fprintf(conn, "-ERR wrong number of arguments for DEL command '%d'\r\n", len(data))
				continue
			}
			count := 0
			for i := range len(data) - 1 {
				if store.Exists(data[i+1]) {
					store.DELValue(data[i+1])
					count++
				}
			}
			fmt.Fprintf(conn, ":%d\r\n", count)

		case "exists":
			if len(data) != 2 {
				fmt.Fprintf(conn, "-ERR wrong number of arguments for DEL command '%d'\r\n", len(data))
				continue
			}
			if ok := store.Exists(data[1]); ok {
				fmt.Fprintf(conn, ":%d\r\n", 1)
			} else {
				fmt.Fprintf(conn, ":%d\r\n", 0)
			}

		default:
			fmt.Fprintf(conn, "-ERR unknown command '%s'\r\n", command)
		}
	}
}
