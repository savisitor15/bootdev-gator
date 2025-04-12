# bootdev-gator
Boot.dev project for Gator

## Requirements
`postgres 14+`
`go 1.23+`

## Setup
Create a database table in postgres, example: `gator`
Setup a username with full rights to the new databse
Finally setup the start up config under your home directory
`~/.gatorconfig.json`
Add the following content
    {
        "db_url": [URI of the database created]
    }
