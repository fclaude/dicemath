{% raw %}
# Dicemath - a small Go service for keeping kids busy

I know this may not be the most fun activity for kids, but it turns out that my 5-year-old likes to do simple addition worksheets. Now that many of us are working from home with our kids, which is mostly trying to work from home while our kids prevent it, we have to take anything we can to keep the little ones busy and entertained.

I can't complain. Working from home with two little kids is hard, but the truth is that we are amazingly lucky that everyone close to us so far is healthy, and on top of that, we can work; I'm amazingly grateful for that!

Back to the article, since my son likes to do worksheets where he can count the result of the additions, I built a small website that generates the worksheet. Every morning I print him one sheet, and he goes to either do the exercises or runs away from me because I'll bug him with something he doesn't want to do at the moment, win-win.

You can check the site here. I'll spend the rest of the article explaining how I built and set the site up. It's a small example you can go through in one day, and it's always fun to pause things and do something different. You can find the whole project here (no, I did not write any tests; this was a one-day fun project; that's it).

## The service

I started super simple, only additions of dice. Later I did spend a bit adding other operations. We will provide a webpage where you can give a name (optional) and what operation you want to practice (addition, subtraction, multiplication, division). Then, once you press a "generate" button, the page redirects you to a PDF with a worksheet addressed to the name provided full of exercises for the operation you picked.

## The code

We will visit the code top-down here. The truth is that I wrote it the other way around, but it's shorter to explain it this way. I tend to work bottom-up and play with the code as I build it; in Go, this doesn't come that easy, but for small things having the main function evolving is easy enough.

Let's start with the webpage:

```html
<!doctype html>
<html lang="en">
  <head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=yes">

    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">

	<style>
        body {
            padding-top: 50px;
        }

        .title {
            padding: 40px 15px;
            text-align: center;
        }
    </style>
    <title>Dice-math worksheet generator</title>
  </head>
  <body>
   
  <div class="container">
    <div class="row">
      <div class="title">
        <h1>Dice-math worksheet generator</h1>	
	<p class="lead">Generate a random dice-math worksheet for your child!</p>
      </div>
    </div>
  </div>
  <div class="container">
    <div class="row">
      <form action="/generate/" method="POST" id="generate-form">
      <div class="form-group">
        <label for="name">Child Name (Optional):</label>
        <input type="text" name="name" class="form-control" id="name" maxlength="20">
        <label for="operation">Operation:</label>
        <select id="operation" name="operation">
          <option value="addition">Addition</option>
          <option value="multiplication">Multiplication</option>
          <option value="subtraction">Subtraction</option>
          <option value="division">Division</option>
        </select> 
      </div>
      <button type="submit" class="btn btn-default btn-primary">Generate</button>
      </form> 
    </div>
  </div>

  <div class="container">
    <div class="row"><p></p></div>
    <div class="row">
      <p>Built to keep my children busy, hope it helps you too :)</p>
    </div>
  </div>
  </body>
</html>
```

This sends a post to /generate/ with the name and the operation. The service will respond to that post request with a PDF. We can start by writing the main binary, with comments:

```go
package main

import (
	"github.com/fclaude/dicemath/generator"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

// This function adds 's to a name, unless the name ends in s, then it only adds '
func ap(c string) string {
	if c == "s" {
		return "'"
	}
	return "'s"
}

// We will limit the length of names to 20 (visually, more looks bad in the PDF
// and I don't want someone having fun knocking my site down with super long
// names).
const MaxName = 20

// cleanName takes only alpha characters and spaces (after truncating). The 
// truncating could have happened after, but keeping it cheap :P.
func cleanName(name string) string {
	if len(name) > MaxName {
		log.Printf("Got a name that is too long: %d", len(name))
		name = name[:MaxName]
	}
	res := ""
	for _, c := range name {
		if (c == ' ') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			res += string(c)
		}
	}
	if res == "" {
		return "My"
	}
	return res + ap(string(res[len(res) - 1]))
}

// Given an operation, we return the symbol (in Latex).
func getOperation(op string) string {
	switch op {
	case "addition": return "+"
	case "multiplication": return "$\\times$"
	case "subtraction": return "-"
	case "division": return "/"
	default: return "?"
	}
}

// Our main!
func main() {
	// We start in release by default, I run in dev mode locally.
	mode := gin.ReleaseMode
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}
	gin.SetMode(mode)
	r := gin.Default()

	// Part of the standard example :)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// The root site is our HTML
	r.GET("/", func(c *gin.Context) {
		c.Data(200, "text/html", []byte(page))
	})

	// This generates the PDF
	// 1) cleans the data.
	// 2) generates the PDF (we will dig into this later).
	// 3) Sends the bytes from the PDF back to the client.
	r.POST("/generate/", func(c *gin.Context) {
		name := cleanName(c.PostForm("name"))
		operation := getOperation(c.PostForm("operation"))
		pdfData, err := generator.GeneratePDF(name, operation)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.Header("Content-Disposition", `attachment; filename="worksheet.pdf"`)
		c.Data(200, "application/pdf", pdfData)
	})

	// again, part of the gin example :)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

const page = `HTML for the page` // I want one single binary, and pakr seems an overkill
```

So this is our service. We only need to implement the code that generates the PDF. For this I’ll do the same thing as the HTML file, I’ll embed a Latex template into my code:

```latex
\documentclass[12pt,letter]{article}

\usepackage{fullpage}
\usepackage{amsmath}
\usepackage{epsdice}
\usepackage{mathptmx}
\usepackage{anyfontsize}

\title{ [[.Name]] worksheet}
\author{}
\date{}
\begin{document}


\newcommand{\dice}[1]{{\fontsize{50}{50}\selectfont \epsdice{#1}}}
\newcommand{\resultline}{..........}
\newcommand{\childsum}[2]{\begin{equation*}{\LARGE \dice{#1}~~\raisebox{14pt}{[[.Operation]]}~~\dice{#2}~~\raisebox{14pt}{=}~~\resultline}\end{equation*}}

\twocolumn
\maketitle
\thispagestyle{empty}

[[range .Group1]]\childsum{[[.A]]}{[[.B]]}
[[end]]

[[range .Group2]]\childsum{[[.A]]}{[[.B]]}
[[end]]

[[.Separator]]

[[range .Group3]]\childsum{[[.A]]}{[[.B]]}
[[end]]

[[range .Group4]]\childsum{[[.A]]}{[[.B]]}
[[end]]

\end{document}
```

This latex template defines a command childsum that takes two numbers and represents them as two dice operated. We have four groups of exercises, and each is a set of childsums. Now let’s use this from within our code, I put this template as a string named sheetTemplate.

We can start with some structures:

```go
// One exercise is two numbers
type exercise struct {
	A int
	B int
}

// Generates one random exercise, both numbers between 1 and 6 at random.
func genExercise() exercise {
	return exercise{A: rand.Intn(6) + 1, B: rand.Intn(6) + 1}
}

// Our exerciseGroup is a set of for exercise arrays, the name of the child,
// and operation, and a separator (a small hack I needed for divisions, more
// on this later.)
type exerciseGroup struct {
	Group1    []exercise
	Group2    []exercise
	Group3    []exercise
	Group4    []exercise
	Name      string
	Operation string
	Separator string
}
```

We will assume for now that the exerciseGroup is built and will focus on generating the PDF:

```go
const fileName = "worksheet.tex"
const fileNamePdf = "worksheet.pdf"

func GeneratePDF(name, operation string) ([]byte, error) {
	// We generate a temporary working directory.
	dir, err := ioutil.TempDir(os.TempDir(), "worksheet")
	// At any point, if something goes bad, we just return nil and
	// the error.
	if err != nil {
		return nil, err
	}
	// After we are done, we remove it, no reason to keep the PDFs.
	defer os.RemoveAll(dir)

	// We parse our Latex template as a text template in Go
	// we use [[ ]] as delimiters, since the default {{ }} does
	// not play nicely with Latex.
	fileContent, err := template.New("sheet").Delims("[[", "]]").Parse(sheetTemplate)
	if err != nil {
		return nil, err
	}

	// We run the template over a generated sheet, and then write it
	// as worksheet.tex in our temporary folder.
	filePath := path.Join(dir, fileName)
	fileWriter, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	err = fileContent.Execute(fileWriter, generateGroups(name, operation))
	fileWriter.Close()
	if err != nil {
		return nil, err
	}

	// We run latex on our temporary folder.
	err = exec.Command("pdflatex", "-aux-directory="+dir, "-output-directory="+dir, filePath).Run()
	if err != nil {
		return nil, err
	}

	// We read and return the bytes of the generated PDF.
	pdf := path.Join(dir, fileNamePdf)
	return ioutil.ReadFile(pdf)
}
```

And that’s almost it. We are missing the generateGroups function:

```go
// This function returns true is the exercise is valid for the given operation.
// For example, we will not consider valid 3-5; at least my son doesn't handle
// negative numbers yet. For divisions we want integer results. And overall, we
// don't want to repeat exercises (the seen map helps us there.)
func valid(seen map[exercise]int, ex exercise, operation string) bool {
	switch operation {
	case "-":
		return ex.A >= ex.B && seen[ex] == 0
	case "/":
		return ex.A%ex.B == 0 && seen[ex] == 0
	}
	return seen[ex] == 0
}

// Our exercise groups have 5 exercises.
const N = 5

// Here we generate our group
func generateGroups(name, operation string) exerciseGroup {
	// we keep a map of exercises we have generated to prevent repetitions.
	seen := make(map[exercise]int, 36)

	// small helper function that generates a random group, references seen
	// so multiple calls still do not repeat exercises
	randomGroup := func() []exercise {
		result := make([]exercise, N)
		for i := range result {
			ex := genExercise()
			for !valid(seen, ex, operation){
				ex = genExercise()
			}
			seen[ex]++
			result[i] = ex
		}
		return result
	}

	// Separator is an empty string, except for division.
	separator := ""
	if operation == "/" {
		// for divisions we will generate only two groups, since
		// we can't generate 20 divisions that satisfy our definition
		// of valid, and therefore, we put a newpage between group 2
		// and 3.
		separator = "\\newpage"
	}

	// Generate groups 1 and 3.
	group1 := randomGroup()
	group2 := make([]exercise, 0, 5)
	group3 := randomGroup()
	group4 := make([]exercise, 0, 5)

	if operation != "/" {
		// If it's not a division, generate groups 2 and 4.
		group2 = randomGroup()
		group4 = randomGroup()
	}

	// We are done!
	return exerciseGroup{
		Name:      name,
		Operation: operation,
		Separator: separator,
		Group1:    group1,
		Group2:    group2,
		Group3:    group3,
		Group4:    group4,
	}
}
```

And now we are done. Our service is ready … almost at least. We need to deploy it.

## Deployment


I deployed this in digitalocean. I run this as a service, in port 8080, and then setup an nginx tI deployed this using digitalocean. I run this as a service, in port 8080, and then set up an Nginx that proxies to it and handles the SSL for me. I also set a firewall so only Nginx can connect to 8080, but not someone from outside.

I started setting up Nginx:

```bash
sudo add-apt-repository ppa:certbot/certbot
sudo apt-get update
sudo apt-get install certbot python-certbot-nginx
sudo certbot --nginx
```

Then edited /etc/nginx/sites-enabled/default:

```
server {
    # managed by Certbot
    server_name dicemath.recoded.cl; 
    listen [::]:443 ssl ipv6only=on; 
    listen 443 ssl; 
    ssl_certificate /etc/letsencrypt/live/dicemath.recoded.cl/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/dicemath.recoded.cl/privkey.pem; 
    include /etc/letsencrypt/options-ssl-nginx.conf; 
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem;

    # my stuff
    location / {
      proxy_pass http://127.0.0.1:8080;
    }
}

server {
    # redirect to https 
    if ($host = dicemath.recoded.cl) {
        return 301 https://$host$request_uri;
    }

    listen 80 ;
    listen [::]:80 ;
    server_name dicemath.recoded.cl;
    return 404;
}
```

The I added my service as /lib/systemd/system/dicemath.service:

```
[Unit]
Description=Dicemath service
ConditionPathExists=/home/fclaude/dicemath/dicemath
After=network.target
 
[Service]
Type=simple
User=dicemath
Group=dicemath
LimitNOFILE=1000000

Restart=on-failure
RestartSec=5
#StartLimitIntervalSec=60

WorkingDirectory=/home/fclaude/dicemath
ExecStart=/home/fclaude/dicemath/dicemath

# make sure log directory exists and owned by syslog
PermissionsStartOnly=true
ExecStartPre=/bin/mkdir -p /var/log/dicemath
ExecStartPre=/bin/chown syslog:adm /var/log/dicemath
ExecStartPre=/bin/chmod 755 /var/log/dicemath
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=dicemath
 
[Install]
WantedBy=multi-user.target
```

The next step was to create the dicemath user and enable the service by running:

```bash
sudo useradd dicemath -s /sbin/nologin -M
sudo systemctl enable dicemath.service
```

And finally, I setup my machine to upload the binary and restart the service whenever I wanted to deploy. I called this file deploy.sh, and it lives in the root of my project:

```bash
#!/bin/sh

go build .

scp dicemath fclaude@dicemath.recoded.cl:/home/fclaude/dicemath/dicemath.new
ssh fclaude@dicemath.recoded.cl "mv /home/fclaude/dicemath/dicemath /home/fclaude/dicemath/dicemath.old"
ssh fclaude@dicemath.recoded.cl "mv /home/fclaude/dicemath/dicemath.new /home/fclaude/dicemath/dicemath"
ssh fclaude@dicemath.recoded.cl "sudo service dicemath restart"
```

And now, I only need to run deploy.sh whenever I change something in the code, it uploads the new binary as dicemath.new, moves the one that is running, renames dicemath.new to dicemath, and restarts the service. The restart will now spin up the new binary.

Note that I’m developing using Linux, if you were in OSX or Windows, you would need to cross-compile. Cross-compiling is easy in Go, just add `GOOS=linux GOARCH=amd64` in front of the go build command, and your binary will be built for Linux.

{% endraw %}
