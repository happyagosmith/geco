<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a name="readme-top"></a>

<!-- PROJECT LOGO -->
<br />
<div align="center">
    <img src="images/logo.png" alt="Logo" width="80" height="80">
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About The Project
TBA

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Built With

This project is built with:

* [![Cobra][Cobra]][Cobra-url]
* [![go-template][go-template]][go-template-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->
## Getting Started

### Prerequisites
TBA

### Installation

#### Build From Source

Ensure that you have a supported version of [Go](https://go.dev/) properly installed and setup. You can find the minimum required version of Go in the [go.mod](go.mod) file.

You can then install the latest release globally by running:
```shell
go install github.com/happyagosmith/geco
```

Or you can install into another directory:
```shell
env GOBIN=/bin go install install github.com/happyagosmith/geco
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- USAGE EXAMPLES -->
## Usage
TBA

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- ROADMAP -->
## Roadmap

- [ ] Add Init sub-command
    - [x] folder
    - [x] single file
    - [ ] merge existing model
    - [ ] generate model.yaml
    - [ ] skip empty file
    - [ ] handle existing file
- [ ] Add Update sub-command
    - [x] add test
    - [x] folder
    - [x] single file
    - [ ] print list generated files 
    - [ ] warning file changed
    - [ ] print list generated files
- [ ] Add Repo sub-command
    - [ ] add
    - [ ] index
    - [ ] list
    - [ ] remove
    - [ ] update

See the [open issues](https://github.com/happyagosmith/geco/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- LICENSE -->
## License

Distributed under the Apache Version 2.0 License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTACT -->
## Contact

Your Name - [@your_twitter](https://twitter.com/your_username) - happyagosmith@gmail.com

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ACKNOWLEDGMENTS -->
## Acknowledgments
TBA

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[license-url]: https://github.com/happyagosmith/geco/blob/main/LICENSE
[go-template]: https://img.shields.io/static/v1?label=text-template&message=latest&color=blue
[go-template-url]: https://pkg.go.dev/text/template
[cobra]: https://img.shields.io/static/v1?label=cobra&message=v1.7.0&color=blue
[cobra-url]: https://github.com/spf13/cobra