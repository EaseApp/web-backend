#!/bin/bash

# Runs the server with live code reloading on file change.

# Build the server.
gb build

# Start the server.
bin/main &
SERVERPID=$!

# Wait for changes.
while true; do
  change=$(inotifywait -e close_write,moved_to,create .)
  change=${change#./ * }
  echo "Restarting server..."
  kill $SERVERPID

  # Build the server.
  gb build

  # Start the server.
  bin/main &
  SERVERPID=$!
done

