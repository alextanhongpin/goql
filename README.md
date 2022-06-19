# goql


Parse query string with rules.


- all operators should be used once only, e.g. `name=eq:john`, to support SQL IN statement, do `name=in:john,jane`. If there are commas in the the text, then do `name=in:"john,doe",jessie`
- some operators cannot belong together, e.g. gt and gte
- nullable types should be supported (how to differentiate text NULL and actual NULL)?
- array, jsonb operations should be supported
- how to add validation to the values?


Reference:
- https://postgrest.org/en/stable/api.html
