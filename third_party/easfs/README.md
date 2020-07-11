# EASFS (EAsy Static Front-end Server)

A web-server to allow fast development and iterations of static websites, based on YAML and Markdown.

It's similar to Jekyll but pages are rendered at request time.

While I did write all the back-end code (Golang), it currently uses the same front-end as Google DevSite,
because I'm lazy to write CSS. But it's a TODO to move this away from Google DevSite CSS & JS.

## Build
```
bazel build //:easfs
```

## Example
https://www.adamyi.com/ is served using EASFS and https://github.com/adamyi/adamyi.com

## License
It was originally a fork to Google's [Web Fundamentals](https://github.com/google/WebFundamentals) project,
but is now rewritten in Golang. Yet it still uses the CSS & JS files from Web Fundamentals.

Copyright 2018-2019 Adam Yi.

Copyright 2014-2018 Google LLC.

Under [Apache 2.0 License](LICENSE).
