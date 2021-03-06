# templated

Simple go template server that supports HTTP calls inside templates. Configurable via JSON file

Usage:

1. Place templates in `${PWD}/templates`

2. Each template needs to have define block:

   ```html
   {{define "index"}}
   <html>
     <head>
       {{template "head"}}
     </head>
     <body>
       {{template "header"}}
       <h1>Index</h1>
     </body>
   </html>
   {{end}}
   ```

3. Configure server via `${PWD}/serve.json`

   ```json
   {
     "pages": [
       { "url": "/", "template": "index" },
       { "url": "/role/{roleName}", "template": "role" }
     ]
   }
   ```

4. You can access path parameters inside templates. E.g. `<h1>Index — Role {{ .Params.roleName }}</h1>`

5. [Sprig](https://github.com/Masterminds/sprig) functions available

6. Making http requests inside templates: `{{ $json := fetchJSON "http://api-container:80" $.Headers }}` _Because these are templates and it is intended for usage within containers — the timeout is 300ms_

7. Request headers are automatically to the `fetchJSON`

8. `fetchJSON` can fetch local files as well: `{{ $json := fetchJSON "file://mock/data.json" $.Headers }}`

9. Use within you docker container:
   ```Dockerfile
   FROM jedrzejlewandowski/templated:1.7.3
   WORKDIR /app
   ADD serve.json /app/serve.json
   ADD templates /app/templates
   CMD ["/bin/templated"]
   ```
