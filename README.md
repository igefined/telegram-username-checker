# Telegram User Checker

Telegram User Checker is a Go-based application that allows you to check the presence of users in a Telegram.

## Features

- Check if users are exist.
- Batch process multiple users from a file.
- Generate reports.

## Requirements

- Go 1.22.0+
- Telegram API credentials (API ID and Hash)

## Usage

**Set Up Telegram API Credentials**

    - Create a new application on the [Telegram API](https://my.telegram.org/auth) site to obtain your `api_id` and `app_hash`.
    - Create a `.env` file in the root directory of the project and add your credentials:

    ```
    APP_ID = 'your_api_id'
    APP_HASH = 'your_api_hash'
    ```

**Prepare Your User List**

    - Create a file with the user details. The file should have a single column without any headers. The file should be based on ./files directory

    ```
    user1
    user2
    user3
    ```

**Run the Checker**

    - Run the script with the following command:

    ```bash
    go run main.go
    ```

**View the Results**

    - Open the `out/result` file to see the status of each user.