---
Title: "Creating A Bespoke Webserver In Go: A Newcomers Perspective"
Author: George Venn
tags:
  - Golang
  - Web-Dev
---

I have been programming for around a decade, starting with Java 5 in 2015, which I self-taught via the O'Reily Java book.
Overtime I learnt how to perform various tasks in Java, but I wanted to learn more. So I looked at other languages.
The years spent learning about Python, C, C++, C#, JavaScript, TypeScript, Rust & x86/64 Assembly is best left for later.

Currently, as the title suggests, I've been using, (and LOVING) Go; I have so far LOVED learning about the language.
I has been so easy to pickup different concepts without having to learn the specifics of the language. In fact, this
website is hosted on a bespoke Go Webserver, using Templ, Echo & HTMX. 

## The Prologue

I have always liked to 'do it myself' and so I sought to self-host my website using a custom web-server. There were MANY attempts
at this goal:
- NodeJS/Express
- Deno/Oak
- Rust/Shuttlers

After watching ThePrimeagen talk about Golang, I thought that I should give it a go; At first I didn't get along to well with the
syntax (Multiple Return Values, Error-as-value, Version-Control based packages), but I very quickly warmed up to it.
It only took me a couple of hours to write a basic Interpreter working through 'The Interpreter Book', for god's sake!
I have Never gotten up and running with a language this fast. 


## Creating a Blog

One thing I wanted to hand-roll as much as possible is *this blog*. There is currently 3 main components:

- a Folder of markdown files - one per post.
- a Sqlite database containing a list of all posts, with metadata
- a Templating system

The lifecycle of this system is:

- Recursively locate every Markdown file under ./posts
- For each post:
	- hash the current timestamp, take the first 5 bytes, convert to base64
	- read the YAML frontmatter to get the post title & author.
	- check if a post with the same title exists, if so, skip.
	- Store into the database.

There are currently a few rough edges to this system that need some work such as:
- What if the title of a post changes, the current system will see it as an entirely new post and create a new entry, meaning a duplicate is created. 
- Unable to remove posts
- Missing a drafts system, what is the best way to implement?

## Moving from net/http to Echo

Because I wanted to have URL Path values for this blog, I tried to use the stock router's PathValue system. Now despite running the correct version of Go (v1.22) & having access to the correct functions, I couldn't get it to work. Now this may of been because some routes were matching when I didn't want them to, but I decided that I should instead use [Echo](https://echo.labstack.com/) 
as my router. This, honestly, was much easier than I expected thanks to Echo's custom Server.Start Function taking in a stock http.Server Implementation. The mess of a Function to create a webserver needed little-to-no changes. All i needed to do was create the routes on Echo.


## What's Next?

From this point, this project is in desperate need of a refactor. The structure looks more like
*This:* 

<img class="post_img" src="/imgs/dependency.png"></img>

*[XKCD - Dependency](https://xkcd.com/2347)*

Than A stable product.

In the mean time, I hope to keep writing about my adventures in Programming.

