# Bandcamper
CLI Bandcamp downloader

Because i listen to lot of unpopular style of music i can't find them in most of popular music streaming platforms.so i deal with Bandcamp.
This is simple CLI tool in development to download from Bandcamp,you can get information of bands (single tracks,albums with them information)then download they works by that provided url. 

## Usage

1. Clone the repository:

```bash
git clone https://github.com/AssassinRobot/Bandcamper.git
cd Bandcamper
```

2. Install dependencies

```bash
go mod tidy
```

3. Get help:

```bash
go run main.go -h
```

Now you can see what is going on

## Authentication

To download resources protected by authentication you can use the following environment variables:

```bash
 export BANDCAMP_USERNAME="your_username"
 export BANDCAMP_COOKIES="client_id=your_client_id; session_id=your_session_id"
```

> Note: to hide the credentials from shell history we added a space before the `export` commands. 

To get your cookies, you need to log to [bandcamp.com](https://bandcamp.com) and open the developer tools in your browser. 
Open the Browser Dev Tools, go to the "Network" tab, and refresh the page. Find the request for the initial HTML page, and look for the "Cookie" header in the request.

