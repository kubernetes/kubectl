# Contributing

## Running Locally

- Install [GitBook Toolchain](https://toolchain.gitbook.com/setup.html)
- From `docs/book` run `npm ci`  to install node_modules locally (don't run install, it updates the shrinkwrap.json)
- From `docs/book` run `npm audit` to make sure there are no vulnerabilities
- From `docs/book` run `gitbook serve`
- Go to `http://localhost:4000` in a browser

## Adding a Section

- Update `SUMMARY.md` with a new section formatted as `## Section Name`

## Adding a Chapter

- Update `SUMMARY.md` under section with chapter formatted as `* [Name of Chapter](pages/section_chapter.md)`
- Add file `pages/section_chapter.md`

## Adding Examples to a Chapter

```bash
{% method %}
Text Explaining Example
{% sample lang="yaml" %}
Formatted code
{% endmethod %}
```

## Adding Notes to a Chapter

```bash
{% panel style="info", title="Title of Note" %}
Note text
{% endpanel %}
```

Notes may have the following styles:

- success
- info
- warning
- danger

## Building and Publishing a release

- Run `gitbook build`
- Push fies in `_book` to a server (e.g. `firebase deploy`)

## Adding GitBook plugins

- Update `book.json` with the plugin
- Run `npm i <npm-plugin-name>`
- Run `npm strinkwrap`
- Run `npm audit` - fix issues by hand updating `npm-shrinkwrap.json` libraries with the fixed versions
- Run `npm ci` - install the new version
- Run `npm audit` - verify the fixes
- Commit `npm-shrinkwrap.json` and `package.json`

### Cool plugins

See https://github.com/swapagarwal/awesome-gitbook-plugins for more plugins.