# RSSFeedReader

## Overview

**RSSFeedReader** is a command-line tool written in Go for managing and reading RSS feeds. It allows users to add, remove, and categorize RSS feed links, as well as view and clean feed descriptions.

## Features

- **Add** RSS feeds with optional categories.
- **Remove** RSS feeds by name.
- **Change** the category of existing RSS feeds.
- **Show** all RSS feeds or those in a specific category.
- **Open** an RSS feed and display its content.
- **Clear** the terminal screen.
- **Save** and **Load** feed links from a JSON file.

## Prerequisites

- Go 1.19 or later
- `github.com/PuerkitoBio/goquery` package for HTML parsing

## Installation

1. **Clone the repository:**

   ```sh
   git clone https://github.com/victor247k/RSSFeedReader.git
   ```

2. **Navigate to the project directory:**

   ```sh
   cd RSSFeedReader
   ```

3. **Install the dependencies:**

   ```sh
   go mod tidy
   ```

## Usage

Run the application using:

```sh
go run main.go
```

### Commands

- `.help` - Displays the help message with available commands.
- `.clear` - Clears the terminal screen.
- `.add <url> [category]` - Adds a new RSS feed with an optional category.
- `.remove <name>` - Removes an RSS feed by its name.
- `.category <name> <newCategory>` - Changes the category of an existing RSS feed.
- `.show [category]` - Shows all RSS feeds or those in a specific category.
- `.open <name>` - Opens an RSS feed and displays its content.
- `.exit` - Exits the application and saves the current feed links to a file.

## Example

```sh
RSSFeedReader> .add https://example.com/feed.xml Technology
RSSFeedReader> .show
Technology:
    1. Example Feed;
RSSFeedReader> .open Example Feed
  Title: Example Feed
  Link: https://example.com/feed.xml
  Description: Example description.
  Language: en
  LastBuildDate: Fri, 19 Jul 2024 15:03:23 GMT
  Generator: Example Generator
  Item 1
      Title: Example Item
      Link: https://example.com/item
      Description: Example item description.
```

## Development

To contribute to the project:

1. **Fork the repository.**
2. **Create a new branch:**

   ```sh
   git checkout -b feature/your-feature
   ```

3. **Commit your changes:**

   ```sh
   git commit -am 'Add new feature'
   ```

4. **Push to the branch:**

   ```sh
   git push origin feature/your-feature
   ```

5. **Create a new Pull Request.**

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

