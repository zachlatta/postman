# Postman ![Analytics](https://ga-beacon.appspot.com/UA-34529482-6/postman/readme?pixel)

**UNDER DEVELOPMENT**

<img src="http://i.imgur.com/eF4fOlg.png" alt="Postman Icon" align="right" />
Postman is a command-line utility for batch sending email.

#### Features

* Fast, templated, bulk emails
* Reads template attributes from CSV
* Works with any SMTP server

### Installation

    $ go get github.com/zachlatta/postman

### Usage

    $ postman [flags]

#### Example

```
$ postman -html template.html -text template.txt -csv recipients.csv \
    -sender "Zaphod Beeblebrox <zaphod@beeblebrox.com>" \
    -subject "Hello, World!" -server smtp.beeblebrox.com -port 587 \
    -user zaphod -password Betelgeuse123
```

template.html:

```
<h1>Hello, {{.Name}}! You are a {{.Type}}</h1>
```

template.txt:

```
Hello, {{.Name}}! You are a {{.Type}}.
```

recipients.csv:

```
Email,Name,Type
arthur@dent.com,Arthur Dent,Human
ford@prefect.com,Ford Prefect,Human
martin@gpp.com,Martin,Robot
trillian@mcmillan.com,Trillian,Human
```

## License

The MIT License (MIT)

Copyright (c) 2014 Zach Latta

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
