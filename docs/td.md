# Technical Design - pfolio
## Entities
### Instrument Identifier
An instrument identifier identifies an instrument.
```go
type InstrumentIdentifier struct {
    Symbol string
}
```

### Model
A model contains a list of instruments and their corresponding weight. It describes allocation of value among instruments in a portfolio.

```go
type ModelEntry struct {
    InstrumentIdentifier InstrumentIdentifier
    Weight   float64
}


type Model struct {
    ID         string
    Name       string
    Entries    []ModelEntry
    CreateTime time.Time
    UpdateTime time.Time
}
```

### Portfolio

## Storage
### Database file
The `PFOLIO_DATABASE` environment variable specifies the path to the `.yml` file under which all states are stored. When `PFOLIO_DATABASE` is empty, this defaults to `~/.pfolio_db.yml`.

Every time the program is executed, the `.yml` file is read into a global variable. Before the program exits successfully, it saves the database back to the `.yml` file. 