# Quick

A template to start new golang projects.

Use MVC project structure.

# TODO

- [ ] add a "error" page to render the message inside a "error layout"
- [ ] go back to "errorMux" and make all middleware as below and treat all the error in the "error mux"?
- [ ] ...
- [ ] create CLI to have generator to add handler, db etc

```go
type HandleFuncWithError func(w http.ResponseWriter, r *http.Request) error
type Middleware func(next HandlerFuncWithError) HandleFuncWithError
```
