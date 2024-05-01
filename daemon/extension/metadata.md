# Overview of HTML Metadata

### General Instructions
We want to send all of this meta to the server in flat JSON where the JSON property name is period-delimited in order to array at as unique value for what will effectively be a set of key-value pairs, e.g.

### Key
- `meta.charset`
- `meta.viewport`
- `base.href`
- `link.stylesheet`
- `script.src`
- `meta.application-name`

Note how for these we drop the attribute name:
- `<meta name>` we drop `.name` 
- `<link rel>` we drop `.rel`

If any of this metadata is duplicated the JSON property name we should use for them, concatenate their values by separating with pipe, unless otherwise specified, such as separate with commas.

### Value
The value should be the obvious other attribute, e.g. for:

- `meta.name` => `content`
- `link.ref` => `href`
- `script` => `src`

And so on. If any are unclear, just ask.

### Relative URLs 
When recording relative URLs, you should adjust them to be absolute using `{base}` or if `<base>` is not provided use the base of the page's URL.

## Categorized Metadata

### Meta Name
- `<meta name="viewport">`
- `<meta name="application-name">`
- `<meta name="description">`
- `<meta name="robots">`
- `<meta name="googlebot">`
- `<meta name="generator">`
- `<meta name="subject">`
- `<meta name="rating">`
- `<meta name="referrer">`
- `<meta name="format-detection">`
- `<meta name="monetization">`

#### For `<name name="google">`
For multiple `name="google"` separate values with commas when concatenating.

### Other Meta
- `<meta charset>`
- `<meta http-equiv>`

### Geo Tags-specific Meta
- `<meta name="ICBM">`
- `<meta name="geo.position">`
- `<meta name="geo.region">`
- `<meta name="geo.placename">`

### Link Rel
- `<link rel="canonical">`
- `<link rel="amphtml">`
- `<link rel="manifest">`
- `<link rel="author">`
- `<link rel="license">`
- `<link rel="me">`
- `<link rel="archives">`
- `<link rel="index">`
- `<link rel="first">`
- `<link rel="last">`
- `<link rel="prev">`
- `<link rel="next">`
- `<link rel="pingback">`
- `<link rel="webmention">`
- `<link rel="micropub">`
- `<link rel="dns-prefetch">`
- `<link rel="preconnect">`
- `<link rel="prefetch">`
- `<link rel="prerender">`

#### For `link.rel`
For `<link rel="{rel}" [type="{type}"] [hreflang="{hreflang}"]>` the JSON property name should be one of these:

- `link.{rel}.{type}`
- `link.{rel}.{type}.{hreflang}`

If there is an `@hreflang` but no `@type` use the type of the current content page.

#### For `link.rel.preload.as`
In the case of `<link rel="preload" as="{as}"> use:

- `link.preload.{as}`

### Miscellaneous
- `<base href>`
- `<script src>`

For `<script src>` only record external URLs.

### Icons

- `<link rel="icon" sizes="192x192" href="/path/to/icon.png">` becomes<br> `{ "link.icon.{size}": "/path/to/icon.png"}`. If `sizes` attribute is not provided, using `{sizes}` as placeholder. 

- `<link rel="apple-touch-icon" href="/path/to/apple-touch-icon.png">` becomes<br> `{ "link.apple-touch-icon": "/path/to/apple-touch-icon.png"}`.

- `<link rel="mask-icon" href="/path/to/icon.svg" color="blue">` becomes<br> `{ "link.mask-icon.blue": "/path/to/icon.svg"}`.


Remember to make the URL values absolute using `<base href>` or by extracting from the page's URL.





