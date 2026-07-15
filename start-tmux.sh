#!/bin/bash

# Ensure the 'dcrouter' session exists
if ! tmux has-session -t dcrouter 2>/dev/null; then
  # Create a new session named 'dcrouter', detached
  cd /workspaces/dcrouter
  tmux new-session -d -s dcrouter
 fi
