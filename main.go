package main

// #include <string.h>
// #include <stdbool.h>
// #include <mysql.h>
// #cgo CFLAGS: -O3 -I/usr/include/mysql -fno-omit-frame-pointer
import "C"
import (
	"encoding/json"
	"log"
	"os"
	"unsafe"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func msg(message *C.char, s string) {
	m := C.CString(s)
	defer C.free(unsafe.Pointer(m))

	C.strcpy(message, m)
}

var l = log.New(os.Stderr, "lambda: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)

//export _lambda_sync_init
func _lambda_sync_init(initid *C.UDF_INIT, args *C.UDF_ARGS, message *C.char) C.bool {
	if args.arg_count != 2 {
		msg(message, "`lambda_sync` requires 2 parameters: the ARN and the JSON payload")
		return C.bool(true)
	}

	argsTypes := (*[2]uint32)(unsafe.Pointer(args.arg_type))

	argsTypes[0] = C.STRING_RESULT
	argsTypes[1] = C.STRING_RESULT

	initid.maybe_null = C.bool(true)

	return C.bool(false)
}

//export _lambda_sync
func _lambda_sync(initid *C.UDF_INIT, args *C.UDF_ARGS, result *C.char, length *uint64, isNull *C.char, message *C.char) *C.char {
	c := 2
	argsArgs := (*[1 << 30]*C.char)(unsafe.Pointer(args.args))[:c:c]
	argsLengths := (*[1 << 30]uint64)(unsafe.Pointer(args.lengths))[:c:c]

	var arn *string
	if argsArgs[0] != nil {
		arn = aws.String(C.GoStringN(argsArgs[0], C.int(argsLengths[0])))
	}

	var jsonPayload *string
	if argsArgs[1] != nil {
		jsonPayload = aws.String(C.GoStringN(argsArgs[1], C.int(argsLengths[1])))
	}

	var b []byte
	if jsonPayload != nil {
		b = []byte(*jsonPayload)
	}

	sess, err := session.NewSession()
	if err != nil {
		l.Printf("failed to create aws session: %v\n", err)

		*length = 0
		*isNull = 1
		return nil
	}

	out, err := lambda.New(sess).Invoke(&lambda.InvokeInput{
		FunctionName: arn,
		Payload:      b,
	})
	if err != nil {
		l.Printf("failed to invoke lambda: %v\n", err)

		*length = 0
		*isNull = 1
		return nil
	}

	j, _ := json.Marshal(out.Payload)

	*length = uint64(len(j))
	*isNull = 0
	return C.CString(string(j))
}

//export lambda_async_init
func lambda_async_init(initid *C.UDF_INIT, args *C.UDF_ARGS, message *C.char) C.bool {
	if args.arg_count != 2 {
		msg(message, "`lambda_sync` requires 2 parameters: the ARN and the JSON payload")
		return C.bool(true)
	}

	argsTypes := (*[2]uint32)(unsafe.Pointer(args.arg_type))

	argsTypes[0] = C.STRING_RESULT
	argsTypes[1] = C.STRING_RESULT

	initid.maybe_null = C.bool(true)

	return C.bool(false)
}

//export lambda_async
func lambda_async(initid *C.UDF_INIT, args *C.UDF_ARGS, result *C.char, length *uint64, isNull *C.char, message *C.char) *C.char {
	c := 2
	argsArgs := (*[1 << 30]*C.char)(unsafe.Pointer(args.args))[:c:c]
	argsLengths := (*[1 << 30]uint64)(unsafe.Pointer(args.lengths))[:c:c]

	var arn *string
	if argsArgs[0] != nil {
		arn = aws.String(C.GoStringN(argsArgs[0], C.int(argsLengths[0])))
	}

	var jsonPayload *string
	if argsArgs[1] != nil {
		jsonPayload = aws.String(C.GoStringN(argsArgs[1], C.int(argsLengths[1])))
	}

	var b []byte
	if jsonPayload != nil {
		b = []byte(*jsonPayload)
	}

	sess, err := session.NewSession()
	if err != nil {
		l.Printf("failed to create aws session: %v\n", err)

		*length = 0
		*isNull = 1
		return nil
	}

	go func() {
		_, err := lambda.New(sess).Invoke(&lambda.InvokeInput{
			FunctionName: arn,
			Payload:      b,
		})
		if err != nil {
			l.Printf("failed to invoke lambda: %v\n", err)
		}
	}()

	*length = 0
	*isNull = 1
	return nil
}

func main() {}
