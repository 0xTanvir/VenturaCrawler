# VenturaCrawler

Venturas Ltd Crawling Engineer (Golang) - adidas.jp Store Technical Test

This project is a technical test for the Crawling Engineer role at Venturas Ltd, focusing on the `adidas.jp` store. The project is structured with the `cmd` directory containing all the commands, while the internal directory holds the implementation details.

The crawler is capable of performing `sync`, `async`, and `parallel` scraping with user agent and proxy rotation. Currently, it is configured with a random delay of max `5` seconds and a parallelism level of `1` to comply with source policies.

# Tools Used

- Make: Build tool for running tasks.
- [Colly V2](https://github.com/gocolly/colly): An elegant scraping and crawling framework for Golang.
- [Cobra](https://github.com/spf13/cobra): A library for creating powerful, modern CLI applications in Golang.

# Installation

To set up the crawler, run the following command:

```bash
make install
```

# Testing the Crawler

To test the crawler with a dump limit of `2` products, run:

```bash
make check
```

This is useful for verifying the setup and checking the output.

# Running the Crawler

To run the crawler, use the following command:

```bash
make run
```
