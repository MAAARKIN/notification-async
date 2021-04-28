## Notification Async

Project created to explain a simple algorithm using worker pool pattern.

### System Requirements

- Golang 1.15

### Getting Start

```bash
git clone https://github.com/MAAARKIN/notification-async.git
```

The project need a csv or text file containing a list of emails like:

```csv
user_email_1@domain.com
user_email_2@domain.com
...
```

To run the sync algorithm:

```bash
go run main.go -filename=<FILE_CSV>  -event=<EVENT_ID> -async=<TRUE/FALSE>

#example
go run main.go -filename=MOCK_DATA100.csv -event=some_event -async=false
```

To run the async algorithm:

```bash
go run main.go -filename=<FILE_CSV> -workers=<NUMBER_OF_WORKERS> -event=<EVENT_ID> -async=<TRUE/FALSE>

#example
go run main.go -filename=MOCK_DATA100.csv -workers=2 -event=some_event -async=true
```