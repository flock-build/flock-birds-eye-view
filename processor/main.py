#!/usr/bin/env python
import json
import os
import sys

import pika


def callback(ch, method, properties, body):
    try:
        data = json.loads(body)
        print(f" [x] Received\n{json.dumps(data, indent=4)}")
    except json.JSONDecodeError as e:
        print(f" [x] Failed to decode JSON: {e}")
        print(f"     Received raw body: {body}")


def main():
    connection = pika.BlockingConnection(pika.ConnectionParameters(host="localhost"))
    channel = connection.channel()

    channel.queue_declare(queue="flock-processor")

    channel.basic_consume(
        queue="flock-processor", on_message_callback=callback, auto_ack=True
    )

    print(" [*] Waiting for messages. To exit press CTRL+C")
    channel.start_consuming()


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("Interrupted")
        try:
            sys.exit(0)
        except SystemExit:
            os._exit(0)
