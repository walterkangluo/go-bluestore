# HOW TO USE

## For simple usage

Import this package:

```
import "github.com/DSiSc/craft/log"
```

And log away:

```
log.DebugKV("This is a debug message.", map[string]interface{}{"hello": "world"})
log.Info("This is a info message.")
```

Logs printed as follows:

```
2018-09-03T16:21:57+08:00 |DEBUG| This is a debug message. caller=/home/kang/Workspace/go/src/hello/jtlog/jtlog_test.go:16 hello=world
2018-09-03T16:21:57+08:00 |INFO| This is a info message. caller=/home/kang/Workspace/go/src/hello/jtlog/jtlog_test.go:17
```

`caller` tells us where we log this record.

## Add an appender

By default, the logs output to console, in `TEXT` format.

If you want to add a new appender which output into a specified log file:

```
log.AddFileAppender("/tmp/aaa/aaa.log", log.InfoLevel, log.JsonFmt, true, true)
```

The file params are: file path, log level, log format, whether to show caller, whether to show timestamp.

And do logging:

```
log.DebugKV("This is a debug message on console and file.", map[string]interface{}{"hello": "world"})
log.Info("This is a info message on console and file.")
```

In STDOUT, log output as always. And in log file `/tmp/aaa/aaa.log`:

```
cat /tmp/aaa/aaa.log
{"level":"info","caller":"/home/kang/Workspace/go/src/hello/jtlog/jtlog_test.go:68","time":"2018-09-03T16:26:31+08:00","message":"This is a info message on console and file."}
```

It only log out `Info` record, with `JSON` format.

You can use `log.AddAppender` to add appender with other `io.Writer`.

## Change logging manners

If you want to change global-log-level:

```
log.SetGlobalConfig(config)
```

Or format of timestamp:

```
log.SetTimestampFormat(time.RFC3339Nano)
```

Multiple changes can be done as the following:

```
config := log.GetGlobalConfig()                 // first, get default configurations
config.TimeStampFormat = time.RFC3339Nano       // then, make changes, such as timestamp format
config.Appenders[0].Format = log.JsonFmt        //                          or logging format of the first Appender
log.SetGlobalConfig(config)                     // finally, refresh configurations with the modified config
```

For pros, just compose whole global `Config` is also OK (but not recommended):

```
config := &log.Config{
    ...
}
log.SetGlobalConfig(config)
```
