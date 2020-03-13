package main

import (
	"fmt"

	"github.com/m-motawea/redisprotocol/resp"
)

func main() {
	// ri := resp.RedisInt{
	// 	Value: 10,
	// }
	// rs := resp.RedisStr{
	// 	Value: "simple redis str",
	// }
	// rbs := resp.RedisBulkStr{
	// 	Value: "bulk redis string",
	// }
	// ras := resp.RedisArray{
	// 	Value: []interface{}{},
	// }
	// ra := resp.RedisArray{
	// 	Value: []interface{}{ri, rs, rbs, ras},
	// }
	// fmt.Println(resp.Encode(ra))
	// rbs1 := "$-1\r\n"
	// rbs2 := "$6\r\nHello\r\n"
	// fmt.Println(resp.DecodeBulkStr(rbs1))
	// fmt.Println(resp.DecodeBulkStr(rbs2))
	// ri := ":5\r\n"
	// fmt.Println(resp.DecodeInt(ri))
	ra1 := "*-1\r\n"
	ra2 := "*0\r\n\r\n"
	ra3 := "*9\r\n*-1\r\n*0\r\n:5\r\n:6\r\n$13\r\nHello World!\r\n+testsimple\r\n$13\r\nHello World!\r\n*2\r\n:7\r\n*-1\r\n-this is error\r\n"
	fmt.Println(resp.DecodeArray(ra1))
	fmt.Println(resp.DecodeArray(ra2))
	fmt.Println(resp.DecodeArray(ra3))
}
