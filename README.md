## Go Nullable

```go
import "github.com/Uffe-Code/go-nullable/nullable"
```

go-nullable is a library that desires to enable null-able support for struct and primitive types.

It is inspired by the way .NET solves it, and by https://github.com/guregu/null.

This package mainly exposes the struct `Nullable[T]` which will work in the same manner as the nullable
type in C#. It exposes two properties, the boolean `HasValue` and the actual `Value`.

The struct implements `encoding.TextMarshaler`, `encoding.TextUnmarshaler`, `json.Marshaler` and `json.Unmarshaler`.
A null objects MarhsalText will return a blank string.

### Usage

This type can be used as properties for API structs.

```go
type Task struct {
TaskId    int                    `json:"task_id"`
ProjectId nullable.Nullable[int] `json:"project_id"`
Subject   string                 `json:"subject"`
}
```

This struct can be read and written to JSON as expected.

```json
{
  "task_id": 5,
  "project_id": null,
  "subject": "Task"
}
```

will be represented as

```go
task := task{
  TaskId: 5,
  ProjectId: nullable.Null[int]()
  Subject: "Task",
}
```

and if we have a value in ProjectId it will look like this:

```go
task := task{
  TaskId: 5,
  ProjectId: nullable.Value(10),
  Subject: "Task",
}
```

The task can be marshalled to JSON normally:

```go
jsonData, err := json.Marshal(task)
```
