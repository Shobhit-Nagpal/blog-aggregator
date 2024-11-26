# Blog Aggregator

Blog aggregator is a CLI-tool written in Go to aggregate blogs. It's a project under Boot.dev

## Prerequisites

- PostgreSQL
- Golang
- Goose
- SQLc

## Installation

```bash
go install github.com/Shobhit-Nagpal/blog-aggreator@latest
```

## Configuration

Set up a `.gatorconfig.json` in your home directory with the following structure

```json
{
  "db_url": "<URL>",
  "current_user_name": "USERNAME"
}
```

## Usage

There a commands defined to use the blog aggregator. Some of them are: 

**addfeed**: Takes in a name and url of the feed
**login**: Login as user
**follow**: Follows a feed
**unfollow**: Unfollow a feed
**browse**: Browses posts from feeds you're following
