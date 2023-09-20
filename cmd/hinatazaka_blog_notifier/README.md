# Running blog_notifier Locally

A simple guide on how to get blog_notifier up and running on your local machine using Docker.

---

## Step 1: Launching with Docker Compose

To start off, initiate the services using the following command:

    docker-compose up hinatazaka-blog-notifier dynamodb-local -d

For first-time users, you'll need to create the table as well:

    docker-compose up create-table -d

Note: This setup will trigger a scraping process every 10 minutes. Ensure you have set up member subscriptions, otherwise, the service won't produce any noticeable activity

---

## Step 2: Stopping blog_notifier

    docker-compose stop hinatazaka-blog-notifier

Remember, if left unchecked, the scraping will continue indefinitely.

---

And that's it! Make sure to manage the service properly to prevent unnecessary scraping. Enjoy your locally running blog_notifier! 🚀
