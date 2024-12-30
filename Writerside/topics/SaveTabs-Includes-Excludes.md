# SaveTabs Include &amp; Exclude Patterns

I envision adding patterns of URLs to include and/or to exclude from logging. 

Patterns will be processed the in order in which users manually sort them. Assumeing the following configuration it will capture anything except FTP URLs and not capture any google.com URLs unless they are at the path `/search`: 

    Include => *://**
    Exclude => ftp://**
    webscheme = http | https
    Exclude => {webscheme}://**.google.com/**
    Include => {webscheme}://**.google.com/search?**

## Regex patterns
I plan to support [RE2 syntax](https://golang.org/s/re2syntax) for includes and excludes that use a Regex pattern, e.g.:

    Include => ^https?://(www\.)?google\.com/search\?

## Glob patterns
I plan to support [Glob syntax](https://stackoverflow.com/a/62985520/102699) for includes and excludes that use a Glob pattern, e.g. _(using a backslash to escape the literal `?`):_

    Include => http?://**google.com/search\?

## Wildcard patterns
I plan to support [Wildcard syntax](https://ss64.com/nt/syntax-wildcards.html) for includes and excludes that use a Wildcard pattern, e.g. _(using a backslash to escape the literal `?`):_

    Include => http?://*google.com/search\?

## URI Template patterns
I plan to support a subset of the [URI Template syntax](https://en.wikipedia.org/wiki/URI_Template) for includes and excludes that use a URI Template:

    webscheme = http | https
    subdomain = www | 
    Include => {webscheme}://{subdomain.}google.com/search/?
