Charming is a small server runtime that wraps Prince XML.

Getting started
---------------

First of all, build the project

    go get github.com/gertv/charming
    go install github.com/gertv/charming/charmingd
    
Next, create the main configuration file. This is a JSON file that defines the template directory, the work directory and the port the server will be listening on.  

```json
{
  "templateDir": "home/charming/templates",
  "workDir": "/home/charming/work",
  "listen": ":6060"
}
```

The templates directory is where the stylesheets are living. To add a template, create a subdirectory that contains the CSS stylesheet and a `template.json` file. 

```json
{
  "name": "my-first-template",
  "stylesheet": "style.css"
}
```

Now, you're ready to start the server - `charmingd <main configuration file`. When the server is running, every template will get its own URL to submit work. In the example config above, the URL would be `http://localhost:6060/submit/my-first-template`.

You can submit work by POST'ing the HTML/XML/... document to be transformed to this URL.

Example using `curl`

    curl -L -d @my-input.html http://localhost:6060/submit/my-first-template
  
You will receive a response JSON document with a UUID, status and output URL. You can follow up on the status of the request on `http://localhost:6060/task/<uuid>`. When the status is `done`, just head to the ouput URL and download the generated PDF document.

```json
{"uuid":"aww275wj9wn8","status":"submitted","outputUrl":"http://localhost:6060/task/aww275wj9wn8/output.pdf"}
```
