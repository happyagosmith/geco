# GECO Project

[![Run ci](https://github.com/happyagosmith/geco/actions/workflows/ci.yml/badge.svg)](https://github.com/happyagosmith/geco/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/happyagosmith/geco)](https://goreportcard.com/report/github.com/happyagosmith/geco)

GECO (Generate, Enrich, Combine, Optimize) is a project designed to streamline the process of file generation for a repository based on a template and a model. The unique aspect of GECO is its ability to handle multiple templates associated with a single repository and manage a single model that is the amalgamation of the models of all the templates added to the repository.

When a new template is added to the repository, the model is enriched with the keys required for this new template. If the keys are already present in the model, they are not overridden. Instead, the existing value is used. This ensures that the model is continuously enriched and updated without losing any existing data.

Key Features:

- **File Generation:** GECO can generate files for a repository based on a template and a model, making it easy to maintain consistency across your project.
- **Multiple Templates:** GECO supports multiple templates associated with a single repository, providing flexibility in file generation.
- **Model Enrichment:** Every time a new template is added to the repository, the model is enriched with the required keys from the new template.
- **Non-destructive Updates:** If the keys from a new template are already present in the model, the existing values are preserved. This ensures that adding new templates does not override existing data.

GECO is an ideal solution for projects that require dynamic file generation based on multiple templates and a continuously evolving model.

### Built With

This project is built with:

* [![Cobra][Cobra]][Cobra-url]
* [![go-template][go-template]][go-template-url]
* [![go-sprig][go-sprig]][go-sprig-url]

## Getting Started

### Prerequisites
Ensure that you have a supported version of [Go](https://go.dev/) properly installed and setup. You can find the minimum required version of Go in the [go.mod](go.mod) file.

You can then install the latest release globally by running:
```shell
go install github.com/happyagosmith/geco@latest
```

## Usage

GECO operates with two main commands:

### `geco init .`

This command is used to initialize the repository with a new template. It enriches the `model-geco.yaml` file in the repository and updates the `.geco.yml` file. The `.geco.yml` file includes the list of all the templates used for the repository.

```bash
geco init .
```

### `geco update`

This command updates the generated files if the templates or the models have changed. It ensures that your repository is always up-to-date with the latest changes in your templates and models.

```bash
geco update
```

## GECO templates

GECO utilizes the [go-template text](https://pkg.go.dev/text/template) syntax for defining its templates. However, instead of the standard "{{...}}" separators used in go-templates, GECO adopts "[[...]]" as separators. This approach is designed to prevent conflicts with generated files that are based on go-templates, such as Helm charts.

Besides the standard functions offered by Go's text/template package, our library also encompasses a collection of custom functions and all functions provided by the Sprig library.

### Using Custom GECO Functions in Go Text Templates

GECO extends the standard Go template functionality with a set of custom functions designed to handle YAML files. These functions are:

### `toYaml`

This function converts the input into a YAML formatted string.

```yaml
[[- with .values ]]
[[- toYaml . 2 ]]
[[- end]]
```

### `readYaml`

This function reads a YAML file from a given path relative to the repo folder and returns a format.Yaml object. If the file does not exist, it returns nil.

```yaml
[[- readYaml "file.yaml" ]]
```

### `readTemplateYaml`

This function reads a YAML file from a given path relative to the template folder, executes any templates within the file, and returns a format.Yaml object. If the file does not exist, it returns nil.

```yaml
[[- readTemplateYaml "file.yaml" ]]
```

### `mergeYaml`

This function merges two YAML objects. If the second object is nil, it returns the string representation of the first object. If there is an error during the merge, it returns an empty string and the error.

```yaml
[[- mergeYaml (readTemplateYaml "template.yaml") (readYaml "file.yaml")]]
```

### `extractSections`

This function extracts sections from a format.Yaml object. If the object is nil, it returns an empty slice.

This is useful to document a yaml file

```yaml
[[- $sections := extractSections (readYaml "values.yaml") ]]
[[ range $sections ]]
### [[ .Title ]]

[[ .Description ]]
[[ $subSection := "TBA" -]]
[[ range .Params -]]
[[ if ne .SubSection $subSection ]]
[[- $subSection = .SubSection ]]
[[- if ne "" $subSection]]
####[[ .SubSection ]]
[[- end]]
| Configuration | Description        | Default value      |
|---------------|--------------------|--------------------|
[[ end -]]
[[ if eq .SubSection $subSection -]]
| [[ .Name ]]   | [[ .Description ]] | <pre>[[ .DefaultValue | replace "\n" "<br/>"]] </pre>|
[[ end -]]
[[ end -]]
[[ end -]]
```


## Roadmap

See the [open issues](https://github.com/happyagosmith/geco/issues) for a full list of proposed features (and known issues).
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

Distributed under the Apache Version 2.0 License. See `LICENSE.txt` for more information.

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[license-url]: https://github.com/happyagosmith/geco/blob/main/LICENSE
[go-template]: https://img.shields.io/static/v1?label=text-template&message=latest&color=blue
[go-template-url]: https://pkg.go.dev/text/template
[cobra]: https://img.shields.io/static/v1?label=cobra&message=v1.7.0&color=blue
[cobra-url]: https://github.com/spf13/cobra
[go-sprig]: https://img.shields.io/static/v1?label=sprig&message=latest&color=blue
[go-sprig-url]: https://masterminds.github.io/sprig/